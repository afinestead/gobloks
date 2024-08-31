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

/* TODO: Namespace these flags somehow */
// Game status flags (bitset)
const (
	FULL        types.Flags = (1 << 0)
	IN_PROGRESS types.Flags = (1 << 1)
	COMPLETE    types.Flags = (1 << 2)
)

// Player status flags (bitset)
const (
	JOINED    types.Flags = (1 << 0) // has joined the game
	CONNECTED types.Flags = (1 << 1) // has socket connection
	DISABLED  types.Flags = (1 << 2) // not playing
	TIMED_OUT types.Flags = (1 << 3) // has timed out
)

const PID_NONE types.PlayerID = 0

type PlayerProfile struct {
	Name   string
	Color  uint
	Pieces game.PieceSet
}

type PlayerState struct {
	profile     PlayerProfile
	status      types.Flags
	socket      *SocketConnection
	playerTimer *utilities.Timer
	// connectionTimeout *utilities.Timer
}

type GameState struct {
	lock           *sync.Mutex
	players        map[types.PlayerID]*PlayerState
	board          *game.Board
	turn           types.PlayerID
	status         types.Flags
	startingPieces game.PieceSet
	socketManager  *SocketManager
	config         types.GameConfig
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
		&sync.Mutex{},
		playerStates,
		board,
		pids[0],
		0,
		pieces,
		InitSocketManager(len(pids)),
		config,
	}
}

func (gs *GameState) nextTurn() {
	// TODO: Game over when no moves remain

	if !gs.config.TurnBased {
		// TODO: Still need to check winner someplace
		return
	}

	var nextUp types.PlayerID = PID_NONE
	for i := 0; i < len(gs.players); i++ {
		maybeNext := types.PlayerID((int(gs.turn)+i)%len(gs.players)) + 1
		// if player is not disabled and has a piece to play, they are next
		if !gs.players[maybeNext].status.Has(DISABLED) &&
			gs.board.HasPlacement(types.Owner(maybeNext), gs.players[maybeNext].profile.Pieces) {
			nextUp = maybeNext
			break
		} else {
			gs.players[maybeNext].status.Set(DISABLED) // No playable pieces, disable player
		}
	}

	fmt.Println("Next up: ", nextUp)

	if nextUp == PID_NONE {
		winners := gs.determineWinners()
		winString := "Game over\n"
		if len(winners) > 1 {
			for i := 0; i < len(winners)-1; i++ {
				winString += winners[i]
				if len(winners) > 2 {
					winString += ", "
				} else {
					winString += " "
				}
			}
			winString += "and " + winners[len(winners)-1] + " tied!"
		} else {
			winString = winners[0] + " wins!"
		}
		gs.sendGameMessage(winString)
		gs.turn = PID_NONE
		gs.status.Set(COMPLETE)
	} else {
		gs.turn = nextUp
		// Start next players timer
		if gs.config.TimeControl > 0 {
			nextPlayer, _ := gs.getPlayer(gs.turn)
			nextPlayer.playerTimer.Start()
		}
	}
}

func (gs *GameState) receiveMessages(player *PlayerState) {
	for {
		var inMsg types.SocketData
		err := gs.socketManager.Recv(player.socket, &inMsg)
		if err != nil {
			fmt.Println(err)
			break
		}
		if inMsg.Type == CHAT_MESSAGE {
			gs.socketManager.Broadcast(&inMsg)
		}
	}

	// handle socket disconnection type events
	gs.lock.Lock()
	defer gs.lock.Unlock()
	gs.socketManager.Disconnect(player.socket)
	player.socket = nil
	player.status.Clear(CONNECTED)
	fmt.Println("Disconnected", player.profile.Name)
	// player.connectionTimeout = utilities.InitTimer(10, 0, func(...any) {
	// 	player.status = types.JOINED // Remove player from active set
	// 	gs.sendPlayerList()          // broadcast updated player list
	// 	// TODO: Find next valid turn and broadcast it
	// 	gs.sendGameMessage(fmt.Sprintf("%s has left the game", player.profile.Name))
	// })
	// player.connectionTimeout.Start()
}

func (gs *GameState) sendGameMessage(msg string) {
	gs.socketManager.Broadcast(&types.SocketData{
		Type: CHAT_MESSAGE,
		Data: &types.ChatMessage{
			Origin:  types.RESERVED,
			Message: msg,
		},
	})
}

func (gs *GameState) sendPlayerList() {
	players := make([]types.PlayerConfig, 0, len(gs.players))
	for pid, player := range gs.players {
		if player != nil {
			players = append(players, types.PlayerConfig{
				PID:    pid,
				Name:   player.profile.Name,
				Color:  player.profile.Color,
				Status: player.status,
				Time:   player.playerTimer.TimeLeftMs(),
			})
		}
	}

	gs.socketManager.Broadcast(&types.SocketData{
		Type: PLAYER_UPDATE,
		Data: players,
	})
}

func (gs *GameState) sendGameStatus() {
	gs.socketManager.Broadcast(&types.SocketData{
		Type: GAME_STATUS,
		Data: &types.PublicGameState{Turn: gs.turn, Status: gs.status},
	})
}

func (gs *GameState) getPlayer(pid types.PlayerID) (*PlayerState, error) {
	player, ok := gs.players[pid]
	if !ok {
		return nil, errors.New("invalid player id")
	}
	return player, nil
}

func (gs *GameState) AddPlayer(name string, color uint) (types.PlayerID, error) {
	/* Assign the new player a PID, if there is one available */
	gs.lock.Lock()
	defer gs.lock.Unlock()

	if gs.status.Has(FULL) {
		return 0, errors.New("game full")
	}

	var pid int
	for pid = 1; pid <= len(gs.players); pid++ {
		if gs.players[types.PlayerID(pid)] == nil {
			gs.players[types.PlayerID(pid)] = &PlayerState{
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
					types.PlayerID(pid),
				),
				// connectionTimeout: nil,
			}
			break
		}
	}

	fmt.Println("Added player ", pid)

	if pid == len(gs.players) {
		fmt.Println("Game is full")
		gs.status.Set(FULL)
	}

	return types.PlayerID(pid), nil
}

func (gs *GameState) ConnectSocket(socket *websocket.Conn, pid types.PlayerID) error {
	gs.lock.Lock()
	defer gs.lock.Unlock()

	player, err := gs.getPlayer(pid)
	if err != nil {
		return err
	}

	if !player.status.Has(JOINED) {
		return errors.New("invalid player status")
	}

	player.socket = gs.socketManager.Connect(socket)
	player.status.Set(CONNECTED)
	fmt.Println(player.status)

	// begin receiving messages on this socket
	go gs.receiveMessages(player)

	// if player.connectionTimeout != nil {
	// 	player.connectionTimeout.Cancel()
	// }

	fmt.Println("Connected player ", pid)

	var playerPieces []types.PublicPiece

	for piece := range player.profile.Pieces {
		pieceCoords := piece.ToPoints()
		piecePoints := make([]types.Point, 0, pieceCoords.Size())
		for coord := range pieceCoords {
			piecePoints = append(piecePoints, types.Point{X: int(coord.X), Y: int(coord.Y)})
		}
		playerPieces = append(playerPieces, types.PublicPiece{
			Hash: piece.Hash(),
			Body: piecePoints,
		})
	}

	// Send all players the current player list and status to sync up
	gs.sendPlayerList()
	gs.sendGameStatus()
	// Send the player their PID, pieces, and current board on connection
	gs.socketManager.Send(
		player.socket,
		&types.SocketData{
			Type: PRIVATE_GAME_STATE,
			Data: &types.PrivateGameState{PID: pid, Pieces: playerPieces},
		},
	)
	gs.socketManager.Send(player.socket, &types.SocketData{Type: BOARD_STATE, Data: gs.board.GetRaw()})
	gs.sendGameMessage(fmt.Sprintf("%s has joined the game", player.profile.Name))

	return nil
}

func (gs *GameState) PlacePiece(pid types.PlayerID, placement types.Placement) error {
	gs.lock.Lock()
	defer gs.lock.Unlock()

	player, err := gs.getPlayer(pid)
	if err != nil {
		return err
	}

	if gs.config.TurnBased && gs.turn != pid {
		return errors.New("not your turn")
	}

	if player.status.Has(TIMED_OUT) {
		return errors.New("player time expired")
	}

	if !gs.status.Has(FULL) {
		return errors.New("waiting for all players")
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

	fmt.Println("PID", pid)
	fmt.Println("placed piece at", origin, piece.ToString())

	gs.status.Set(IN_PROGRESS)

	gs.socketManager.Broadcast(&types.SocketData{Type: BOARD_STATE, Data: gs.board.GetRaw()})

	player.playerTimer.Pause() // successfully placed piece, pause timer

	player.profile.Pieces.Remove(piece)

	gs.nextTurn()
	// Send updated times and turn to all players
	gs.sendPlayerList()
	gs.sendGameStatus()

	return nil
}

func (gs *GameState) determineWinners() []string {
	winners := make([]string, 0, len(gs.players))
	minScore := 0xffffffff
	for _, player := range gs.players {
		if player != nil {
			score := 0
			for piece := range player.profile.Pieces {
				score += int(piece.Size())
			}
			fmt.Println(player.profile.Name, score)
			if score <= minScore {
				minScore = score
				winners = append(winners, player.profile.Name)
			}
		}
	}
	return winners
}

func (gs *GameState) handleTimeout(args ...any) {
	gs.lock.Lock()
	defer gs.lock.Unlock()

	pid := args[0].(types.PlayerID)
	player, _ := gs.getPlayer(pid)
	player.status.Set(TIMED_OUT | DISABLED)
	gs.nextTurn()
	gs.sendPlayerList()
	gs.sendGameStatus()
}
