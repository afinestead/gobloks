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
}

type PlayerState struct {
	profile PlayerProfile
	status  PlayerStatus
	socket  *websocket.Conn
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

	fmt.Println(config)

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

	return &GameState{playerStates, board, pids[0], pieces}
}

func (gs *GameState) ConnectPlayer(name string, color uint) (types.PlayerID, error) {
	/* Assign the new player a PID, if there is one available */

	for pid, player := range gs.players {
		if player.status == NONE { // No player connected at this pid yet
			*player = PlayerState{
				profile: PlayerProfile{
					Name:   name,
					Color:  color,
					Pieces: game.PieceSet{}, // fill in during a "StartGame" function
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
	fmt.Println(player, player.status)
	if player.status != JOINED && player.status != DISCONNECTED {
		return errors.New("invalid player status")
	}
	player.socket = socket
	player.status = CONNECTED

	for {
		_, msg, err := socket.ReadMessage()
		if err != nil {
			break
		}
		fmt.Println(msg)
	}

	socket.Close()

	player.socket = nil
	player.status = DISCONNECTED

	return nil
}
