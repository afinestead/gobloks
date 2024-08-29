package manager

import (
	"errors"
	"fmt"
	"gobloks/internal/game"
	"gobloks/internal/types"
	"gobloks/internal/utilities"
	"sync"

	"github.com/gorilla/websocket"
)

const (
	NONE      types.PlayerStatus = 0x0
	JOINED    types.PlayerStatus = 0x1 // has joined the game
	CONNECTED types.PlayerStatus = 0x2 // has socket connection
	ACTIVE    types.PlayerStatus = 0x4 // has playable pieces
	TIMED_OUT types.PlayerStatus = 0x8 // has timed out
)

type PlayerProfile struct {
	Name   string
	Color  uint
	Pieces game.PieceSet
}

type GameState struct {
	players        map[types.PlayerID]*PlayerState
	board          *game.Board
	turn           types.PlayerID
	startingPieces game.PieceSet
	socketManager  *SocketManager
	config         types.GameConfig
	playersLock    *sync.Mutex
}

type PlayerState struct {
	profile     PlayerProfile
	status      types.PlayerStatus
	socket      *SocketConnection
	playerTimer *utilities.Timer
	// connectionTimeout *utilities.Timer
	mtx *sync.Mutex
}

func InitGameState(config types.GameConfig) *GameState {

	pieces, setPixels, err := game.GeneratePieceSet(config.BlockDegree) // TODO: cache
	if err != nil {
		return nil
	}

	playerStates := make(map[types.PlayerID]*PlayerState, config.Players)
	pids := make([]types.PlayerID, config.Players)

	for ii := 1; ii <= len(pids); ii++ {
		pid := types.PlayerID(ii)
		pids[ii-1] = pid
		playerStates[pid] = nil
	}

	board, err := game.NewBoard(pids, setPixels, config.Density)
	if err != nil {
		return nil
	}

	return &GameState{
		playerStates,
		board,
		pids[0],
		pieces,
		InitSocketManager(len(pids)),
		config,
		&sync.Mutex{},
	}
}

func (gs *GameState) nextTurn() (types.PlayerID, error) {
	if !gs.config.TurnBased {
		return types.PID_NONE, nil
	}

	gs.playersLock.Lock()
	defer gs.playersLock.Unlock()
	for i := 0; i < len(gs.players); i++ {
		nextUp := types.PlayerID((int(gs.turn)+i)%len(gs.players)) + 1
		if gs.board.HasPlacement(types.Owner(nextUp), gs.players[nextUp].profile.Pieces) {
			return nextUp, nil
		}
	}

	return types.PID_NONE, errors.New("no valid moves")
}

func (gs *GameState) sendGameMessage(msg string) {
	gs.socketManager.Broadcast(&types.SocketData{
		Type: types.CHAT_MESSAGE,
		Data: &types.ChatMessage{
			Origin:  types.RESERVED,
			Message: msg,
		},
	})
}

func (gs *GameState) sendPlayerList() {
	gs.playersLock.Lock()
	defer gs.playersLock.Unlock()
	players := make([]types.PlayerConfig, 0, len(gs.players))
	for pid, player := range gs.players {
		if player != nil {
			player.mtx.Lock()
			players = append(players, types.PlayerConfig{
				PID:    pid,
				Name:   player.profile.Name,
				Color:  player.profile.Color,
				Status: player.status,
				Time:   player.playerTimer.TimeLeftMs(),
			})
			player.mtx.Unlock()
		}
	}

	gs.socketManager.Broadcast(&types.SocketData{
		Type: types.PLAYER_UPDATE,
		Data: players,
	})
}

func (gs *GameState) hasPlayerStatus(player *PlayerState, status types.PlayerStatus) bool {
	player.mtx.Lock()
	defer player.mtx.Unlock()
	return player.status&status != 0
}

func (gs *GameState) setPlayerStatus(player *PlayerState, status types.PlayerStatus) {
	player.mtx.Lock()
	defer player.mtx.Unlock()
	player.status |= status
}

func (gs *GameState) clearPlayerStatus(player *PlayerState, status types.PlayerStatus) {
	player.mtx.Lock()
	defer player.mtx.Unlock()
	player.status &= ^status
}

func (gs *GameState) getPlayer(pid types.PlayerID) (*PlayerState, error) {
	gs.playersLock.Lock()
	defer gs.playersLock.Unlock()
	player, ok := gs.players[pid]
	if !ok {
		return nil, errors.New("invalid player id")
	}
	return player, nil
}

func (gs *GameState) AddPlayer(name string, color uint) (types.PlayerID, error) {
	/* Assign the new player a PID, if there is one available */
	gs.playersLock.Lock()
	defer gs.playersLock.Unlock()
	for pid, player := range gs.players {
		if player == nil { // No player connected at this pid yet
			gs.players[pid] = &PlayerState{
				profile: PlayerProfile{
					Name:   name,
					Color:  color,
					Pieces: gs.startingPieces.Copy(),
				},
				status: JOINED,
				socket: nil,
				playerTimer: utilities.InitTimer(
					gs.config.TimeControl*1000,
					gs.config.TimeBonus*1000,
					gs.handleTimeout,
					pid,
				),
				// connectionTimeout: nil,
				mtx: &sync.Mutex{},
			}
			return pid, nil
		}
	}
	return 0, errors.New("game full")
}

func (gs *GameState) ConnectSocket(socket *websocket.Conn, pid types.PlayerID) error {
	player, err := gs.getPlayer(pid)
	if err != nil {
		return err
	}

	if !gs.hasPlayerStatus(player, JOINED) {
		return errors.New("invalid player status")
	}

	player.socket = gs.socketManager.Connect(socket)
	gs.setPlayerStatus(player, CONNECTED)

	// if player.connectionTimeout != nil {
	// 	player.connectionTimeout.Cancel()
	// }

	fmt.Println("Connected player ", pid)

	// Send all players the current player list to sync up
	gs.sendPlayerList()

	gs.socketManager.Send(
		player.socket,
		&types.SocketData{
			Type: types.PUBLIC_GAME_STATE,
			Data: &types.PublicGameState{
				Board: gs.board.GetRaw(),
				Turn:  gs.turn,
			},
		},
	)

	var playerPieces []types.PublicPiece

	for piece := range player.profile.Pieces {
		pieceCoords := piece.ToPoints()
		piecePoints := make([]types.Point, 0, pieceCoords.Size())
		for coord := range pieceCoords {
			piecePoints = append(piecePoints, types.Point{
				X: int(coord.X),
				Y: int(coord.Y),
			})
		}
		playerPieces = append(playerPieces, types.PublicPiece{
			Hash: piece.Hash(),
			Body: piecePoints,
		})
	}

	// Send the player their PID and pieces on connection
	gs.socketManager.Send(
		player.socket,
		&types.SocketData{
			Type: types.PRIVATE_GAME_STATE,
			Data: &types.PrivateGameState{
				PID:    pid,
				Pieces: playerPieces,
			},
		},
	)

	gs.sendGameMessage(fmt.Sprintf("%s has joined the game", player.profile.Name))

	for {
		var inMsg types.SocketData
		err := gs.socketManager.Recv(player.socket, &inMsg)
		if err != nil {
			fmt.Println(err)
			break
		}
		if inMsg.Type == types.CHAT_MESSAGE {
			gs.socketManager.Broadcast(&inMsg)
		}
	}

	// handle socket disconnection type events
	gs.socketManager.Disconnect(player.socket)
	player.socket = nil
	gs.clearPlayerStatus(player, CONNECTED)
	// player.connectionTimeout = utilities.InitTimer(10, 0, func(...any) {
	// 	player.status = types.JOINED // Remove player from active set
	// 	gs.sendPlayerList()          // broadcast updated player list
	// 	// TODO: Find next valid turn and broadcast it
	// 	gs.sendGameMessage(fmt.Sprintf("%s has left the game", player.profile.Name))
	// 	fmt.Println("Disconnected player ", pid)
	// })
	// player.connectionTimeout.Start()

	return nil
}

func (gs *GameState) PlacePiece(pid types.PlayerID, placement types.Placement) error {
	player, err := gs.getPlayer(pid)
	if err != nil {
		return err
	}

	if gs.config.TurnBased && gs.turn != pid {
		return errors.New("not your turn")
	}

	if gs.hasPlayerStatus(player, TIMED_OUT) {
		return errors.New("player time expired")
	}

	// convert placement to Piece/origin
	relPoints, origin := utilities.NormalizeToOrigin(utilities.NewSet(placement.Coordinates))
	relCoords := utilities.NewSet([]game.PieceCoord{}, relPoints.Size())
	for coord := range relPoints {
		relCoords.Add(game.PieceCoord{X: uint8(coord.X), Y: uint8(coord.Y)})
	}

	piece := game.PieceFromPoints(relCoords)
	if !player.profile.Pieces.Has(piece) {
		return errors.New("player does not have this piece")
	}

	_, err = gs.board.Place(origin, piece, types.Owner(pid))
	if err != nil {
		return err
	}

	player.playerTimer.Pause() // successfully placed piece, pause timer

	player.profile.Pieces.Remove(piece)

	nextUp, err := gs.nextTurn()
	if err != nil {
		winner, _ := gs.determineWinner()
		gs.getPlayer(winner)
		gs.sendGameMessage("Game over!")
		gs.turn = types.PID_NONE
		// gs.sendGameMessage(fmt.Sprintf("%s wins!", p.profile.Name))
	} else {
		gs.turn = nextUp
	}

	gs.socketManager.Broadcast(&types.SocketData{
		Type: types.PUBLIC_GAME_STATE,
		Data: &types.PublicGameState{
			Board: gs.board.GetRaw(),
			Turn:  gs.turn,
		},
	})

	// Start next players timer
	if gs.turn != types.PID_NONE && gs.config.TimeControl > 0 {
		nextPlayer, _ := gs.getPlayer(gs.turn)
		nextPlayer.playerTimer.Start()
	}

	// Send updated times to all players
	gs.sendPlayerList()

	return nil
}

func (gs *GameState) determineWinner() (types.PlayerID, error) {
	return types.PID_NONE, nil
}

func (gs *GameState) handleTimeout(args ...any) {
	pid := args[0].(types.PlayerID)
	player, _ := gs.getPlayer(pid)
	gs.setPlayerStatus(player, TIMED_OUT)
	gs.sendPlayerList()
}
