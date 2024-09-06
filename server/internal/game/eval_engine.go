package game

import (
	"gobloks/internal/types"
	"gobloks/internal/utilities"
)

const (
	WEIGHT_PLACEMENTS float64 = 1.0
	WEIGHT_CORNERS    float64 = 1.0
	WEIGHT_TERRITORY  float64 = 1.0
)

type EvalEngine struct {
	depth    uint
	chRecv   chan *EvalState
	chCancel chan struct{}
	chResult chan map[types.PlayerID]float64
}

type EvalState struct {
	game    *GameState
	players map[types.PlayerID]*PlayerState
}

func InitEvalEngine(depth uint) *EvalEngine {
	return &EvalEngine{
		depth:    depth,
		chRecv:   make(chan *EvalState),
		chCancel: make(chan struct{}),
		chResult: make(chan map[types.PlayerID]float64),
	}
}
func (engine *EvalEngine) Evaluate(state *EvalState) {
	engine.chRecv <- state
}

func (engine *EvalEngine) Start() {
	for state := range engine.chRecv {
		// engine.chCancel <- struct{}{}
		eval := make(map[types.PlayerID]float64, len(state.players))
		err := engine.evaluateGameState(state, 0, eval)
		if err == nil {
			engine.chResult <- eval
		}
	}
}

func (engine *EvalEngine) Stop() {
	close(engine.chRecv)
	close(engine.chCancel)
	close(engine.chResult)
}

func (engine *EvalEngine) evaluateGameState(state *EvalState, curDepth int, curRes map[types.PlayerID]float64) error {
	if curDepth >= int(engine.depth) || state.game.status.Has(COMPLETE) {
		return nil
	}

	// iterate in order, starting with the current turn
	for i := 0; i < len(state.players); i++ {
		next := state.players[types.PlayerID((int(state.game.turn)+i)%len(state.players))+1]
		if !next.status.Has(DISABLED) {
			select {
			case <-engine.chCancel:
				return nil
			default:
				territory := state.game.board.findTerritory(types.Owner(next.pid))
				corners := state.game.board.findCorners(territory, types.Owner(next.pid))
				placements := state.game.board.getPlacements(corners, types.Owner(next.pid), state.players[next.pid].pieces, false)

				eval := float64(len(territory))*WEIGHT_TERRITORY + float64(len(corners))*WEIGHT_CORNERS + float64(len(placements))*WEIGHT_PLACEMENTS
				curRes[next.pid] += eval

				for _, plc := range placements {
					gsCopy := state.game.Copy()
					playersCopy := make(map[types.PlayerID]*PlayerState, len(state.players))
					for pid, player := range state.players {
						playersCopy[pid] = player.Copy()
					}

					gsCopy.board.Place(utilities.NewSet(plc.Coordinates), types.Owner(next.pid))
					engine.evaluateGameState(&EvalState{gsCopy, playersCopy}, curDepth+1, curRes)
				}
			}
		}
	}

	return nil
}
