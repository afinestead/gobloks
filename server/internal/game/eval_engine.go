package game

import (
	"gobloks/internal/types"
)

const (
	WEIGHT_PLACEMENTS float64 = 1.0
	WEIGHT_CORNERS    float64 = 1.0
	WEIGHT_TERRITORY  float64 = 1.0
	WEIGHT_OPEN_SPACE float64 = 1.0
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
				curRes[next.pid] += engine.evaluatePlayerPosition(state, next.pid, curDepth)
			}
		}
	}

	return nil
}

func (engine *EvalEngine) evaluatePlayerPosition(state *EvalState, pid types.PlayerID, curDepth int) float64 {
	eval := 0.0
	// territory := state.game.board.findTerritory(types.Owner(pid))
	// corners := state.game.board.findCorners(territory, types.Owner(pid))
	// placements := state.game.board.getPlacements(corners, types.Owner(pid), state.players[pid].pieces, false)

	// numPlacements := 0
	// playableArea := 0
	// for p := placements; p != nil; p = p.Next {
	// 	playableArea += len(p.Value)
	// 	numPlacements++
	// }

	// fmt.Println("evaluating PID", pid)
	// fmt.Println(len(territory), len(corners), numPlacements, playableArea)

	// eval := (float64(len(territory))*WEIGHT_TERRITORY +
	// 	float64(len(corners))*WEIGHT_CORNERS +
	// 	float64(numPlacements)*WEIGHT_PLACEMENTS +
	// 	float64(playableArea)*WEIGHT_OPEN_SPACE)

	// for p := placements; p != nil; p = p.Next {
	// 	gsCopy := state.game.Copy()
	// 	playersCopy := make(map[types.PlayerID]*PlayerState, len(state.players))
	// 	for pid, player := range state.players {
	// 		playersCopy[pid] = player.Copy()
	// 	}

	// 	gsCopy.board.Place(utilities.NewSet(p.Value), pid)
	// 	engine.evaluateGameState(&EvalState{gsCopy, playersCopy}, curDepth+1)
	// }

	return eval
}
