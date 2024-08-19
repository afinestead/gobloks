package manager

import (
	"errors"
	"fmt"
	"gobloks/internal/game"
	"gobloks/internal/types"

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
}

type PlayerState struct {
	profile PlayerProfile
	status  PlayerStatus
	socket  *SocketConnection
}

type PlayerStatus uint

const (
	NONE PlayerStatus = iota
	JOINED
	CONNECTED
	DISCONNECTED
)

func InitGameState(config types.GameConfig) *GameState {

	pieces, setPixels, err := game.GeneratePieceSet(config.BlockDegree) // TODO: cache
	if err != nil {
		fmt.Println(err)
	}

	playerStates := make(map[types.PlayerID]*PlayerState, config.Players)
	pids := make([]types.PlayerID, config.Players)

	for ii := 0; ii < len(pids); ii++ {
		pid := types.PlayerID(ii)
		pids[ii] = pid
		playerStates[pid] = &PlayerState{}
	}

	board, err := game.NewBoard(pids, setPixels, config.Density)
	if err != nil {
		fmt.Println(err)
	}

	return &GameState{
		playerStates,
		board,
		pids[0],
		pieces,
		InitSocketManager(),
	}
}

func (gs *GameState) GetPlayer(pid types.PlayerID) (PlayerProfile, error) {
	player, ok := gs.players[pid]
	if !ok {
		return PlayerProfile{}, errors.New("invalid player id")
	}
	return player.profile, nil
}

func (gs *GameState) ConnectPlayer(name string, color uint) (types.PlayerID, error) {
	/* Assign the new player a PID, if there is one available */

	for pid, player := range gs.players {
		if player.status == NONE { // No player connected at this pid yet
			*player = PlayerState{
				profile: PlayerProfile{
					Name:   name,
					Color:  color,
					Pieces: gs.startingPieces,
				},
				status: JOINED,
				socket: nil,
			}

			return pid, nil
		}
	}
	return 0, errors.New("game full")
}

func (gs *GameState) ConnectSocket(socket *websocket.Conn, pid types.PlayerID) error {
	player := gs.players[pid]
	if player.status != JOINED && player.status != DISCONNECTED {
		return errors.New("invalid player status")
	}

	player.socket = gs.socketManager.Connect(socket)
	player.status = CONNECTED

	fmt.Println("Connected player ", pid)

	// Send all players the current game state to sync up
	gs.socketManager.Broadcast(&types.PublicGameState{
		Board: gs.board.GetRaw(),
		Turn:  gs.turn,
		Players: func() []types.PlayerConfig {
			players := make([]types.PlayerConfig, 0, len(gs.players))
			for pid, player := range gs.players {
				players = append(players, types.PlayerConfig{
					PID:   pid,
					Name:  player.profile.Name,
					Color: player.profile.Color,
				})
			}
			return players
		}(),
	})

	var playerPieces [][]types.Point

	for piece := range player.profile.Pieces {
		pieceCoords := piece.ToPoints()
		piecePoints := make([]types.Point, 0, pieceCoords.Size())
		for coord := range pieceCoords {
			piecePoints = append(piecePoints, types.Point{
				X: int(coord.X),
				Y: int(coord.Y),
			})
		}
		playerPieces = append(playerPieces, piecePoints)
	}

	// Send the player their PID and pieces on connection
	gs.socketManager.Send(player.socket, &types.PrivateGameState{
		PID:    pid,
		Pieces: playerPieces,
	})

	gs.socketManager.Broadcast(&types.ChatMessage{
		Origin:  types.RESERVED,
		Message: fmt.Sprintf("%s has joined the game", player.profile.Name),
	})

	for {
		var inMsg types.ChatMessage
		err := gs.socketManager.Recv(player.socket, &inMsg)
		if err != nil {
			fmt.Println(err)
			break
		}
		if inMsg.Message != "" {
			gs.socketManager.Broadcast(&inMsg)
		}
	}

	gs.socketManager.Disconnect(player.socket)
	player.socket = nil
	player.status = DISCONNECTED

	gs.socketManager.Broadcast(&types.ChatMessage{
		Origin:  types.RESERVED,
		Message: fmt.Sprintf("%s has left the game", player.profile.Name),
	})

	fmt.Println("Disconnected player ", pid)

	return nil
}
