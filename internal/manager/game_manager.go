package manager

import (
	"fmt"
	"gobloks/internal/types"
	"math/rand"
)

type GameManager struct {
	mangagedGames map[types.GameID]*GameState
}

func InitGameManager() *GameManager {
	return &GameManager{make(map[types.GameID]*GameState, types.MANAGED_GAMES_START_SIZE)}
}

const letterBytes = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"

func createGameID(n uint8) types.GameID {
	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return types.GameID(b)
}

func (gs *GameManager) CreateGame(config types.GameConfig) types.GameID {
	var gid types.GameID
	for {
		gid = createGameID(4)
		_, ok := gs.mangagedGames[gid]
		if !ok { // gid not in mangagedGames
			break
		}
	}

	gs.mangagedGames[gid] = InitGameState(config)

	fmt.Println(gs.mangagedGames[gid])

	return gid
}

func (gs *GameManager) FindGame(gid types.GameID) (*GameState, error) {
	game, ok := gs.mangagedGames[gid]
	if !ok {
		return nil, fmt.Errorf("invalid game id `%s`", gid)
	}
	return game, nil
}

func (gs *GameManager) DeleteGame(gid types.GameID) error {
	if _, ok := gs.mangagedGames[gid]; !ok {
		return fmt.Errorf("invalid game id `%s`", gid)
	}
	delete(gs.mangagedGames, gid)
	return nil
}
