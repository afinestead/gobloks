package manager

import (
	"errors"
	"fmt"
	"gobloks/internal/game"
	"gobloks/internal/types"
)

type PlayerProfile struct {
	Name   string
	Color  uint
	Pieces game.PieceSet
}

type GameState struct {
	players        []PlayerState
	board          *game.Board
	turn           types.PlayerID
	startingPieces game.PieceSet
}

type PlayerState struct {
	profile PlayerProfile
	status  PlayerStatus
	socket  int
}

type PlayerStatus uint

const (
	NONE PlayerStatus = iota
	CONNECTED
	DISCONNECTED
)

func InitGameState(config types.GameConfig) *GameState {

	pieces, setPixels, err := game.GeneratePieceSet(config.BlockDegree) // TODO: cache
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(config)

	playerStates := make([]PlayerState, config.Players)
	pids := make([]types.PlayerID, config.Players)

	for ii := 0; ii < len(pids); ii++ {
		pids[ii] = types.PlayerID(ii)
	}

	board, err := game.NewBoard(pids, setPixels, config.Density)
	if err != nil {
		fmt.Println(err)
	}

	return &GameState{playerStates, board, pids[0], pieces}
}

func (gs *GameState) ConnectPlayer(name string, color uint) (types.PlayerID, error) {
	for pid, player := range gs.players {
		if player.status == NONE { // No player connected at this pid yet
			player = PlayerState{
				profile: PlayerProfile{
					Name:   name,
					Color:  color,
					Pieces: game.PieceSet{}, // fill in during a "StartGame" function
				},
				status: CONNECTED,
				socket: 0,
			}

			return types.PlayerID(pid), nil
		}
	}
	return 0, errors.New("game full")
}
