package manager

import (
	"errors"
	"fmt"
	"gobloks/internal/game"
	"gobloks/internal/types"
	"gobloks/internal/utilities"

	"github.com/gorilla/websocket"
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
}

type PlayerState struct {
	profile           PlayerProfile
	status            types.PlayerStatus
	socket            *SocketConnection
	connectionTimeout *utilities.Timer
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
		playerStates[pid] = &PlayerState{}
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
	}
}

func (gs *GameState) nextTurn() {
	if !gs.config.TurnBased {
		gs.turn = types.PID_NONE
		return
	}

	nextUp := (gs.turn + 1) % types.PlayerID(len(gs.players))
	if gs.board.HasPlacement(types.Owner(nextUp), gs.players[nextUp].profile.Pieces) {
	}
	// 	gs.turn = nextUp
	// } else {
	// 	// gs.nextTurn()
	// }
	gs.turn = nextUp
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

func (gs *GameState) getPlayer(pid types.PlayerID) (*PlayerState, error) {
	player, ok := gs.players[pid]
	if !ok {
		return nil, errors.New("invalid player id")
	}
	return player, nil
}

func (gs *GameState) getActivePlayers() []types.PlayerConfig {
	players := make([]types.PlayerConfig, 0, len(gs.players))
	for pid, player := range gs.players {
		players = append(players, types.PlayerConfig{
			PID:    pid,
			Name:   player.profile.Name,
			Color:  player.profile.Color,
			Status: player.status,
		})
	}
	return players
}

func (gs *GameState) ConnectPlayer(name string, color uint) (types.PlayerID, error) {
	/* Assign the new player a PID, if there is one available */

	for pid, player := range gs.players {
		if player.status == types.NONE { // No player connected at this pid yet
			*player = PlayerState{
				profile: PlayerProfile{
					Name:   name,
					Color:  color,
					Pieces: gs.startingPieces.Copy(),
				},
				status:            types.JOINED,
				socket:            nil,
				connectionTimeout: nil,
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
	if player.status&types.JOINED == 0 {
		return errors.New("invalid player status")
	}

	player.socket = gs.socketManager.Connect(socket)
	player.status |= types.CONNECTED

	if player.connectionTimeout != nil {
		player.connectionTimeout.Cancel()
	}

	fmt.Println("Connected player ", pid)

	// Send all players the current player list to sync up
	gs.socketManager.Broadcast(&types.SocketData{
		Type: types.PLAYER_UPDATE,
		Data: &types.ActivePlayers{Players: gs.getActivePlayers()},
	})

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
	player.status &= ^types.CONNECTED
	player.connectionTimeout = utilities.InitTimer(10, func() {
		player.status = types.NONE // Remove player from active set

		// broadcast updated player list
		gs.socketManager.Broadcast(&types.SocketData{
			Type: types.PLAYER_UPDATE,
			Data: &types.ActivePlayers{Players: gs.getActivePlayers()},
		})

		gs.sendGameMessage(fmt.Sprintf("%s has left the game", player.profile.Name))
	})
	player.connectionTimeout.Start()

	fmt.Println("Disconnected player ", pid)

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

	fmt.Println(gs.board.ToString())

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
	player.profile.Pieces.Remove(piece)

	gs.board.HasPlacement(types.Owner(pid), player.profile.Pieces)

	gs.nextTurn()

	gs.socketManager.Broadcast(&types.SocketData{
		Type: types.PUBLIC_GAME_STATE,
		Data: &types.PublicGameState{
			Board: gs.board.GetRaw(),
			Turn:  gs.turn,
		},
	})

	return nil
}
