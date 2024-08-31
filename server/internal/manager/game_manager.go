package manager

import (
	"fmt"
	"gobloks/internal/types"
	"math/rand"
	"sync"
)

type GameManager struct {
	mangagedGames map[types.GameID]*GameState
	lock          *sync.Mutex
}

func InitGameManager() *GameManager {
	return &GameManager{
		make(map[types.GameID]*GameState, types.MANAGED_GAMES_START_SIZE),
		&sync.Mutex{},
	}
}

const letterBytes = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"

func createGameID(n uint8) types.GameID {
	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return types.GameID(b)
}

func (gm *GameManager) CreateGame(config types.GameConfig) types.GameID {
	var gid types.GameID
	gm.lock.Lock()
	defer gm.lock.Unlock()
	for {
		gid = createGameID(4)
		_, ok := gm.mangagedGames[gid]
		if !ok { // unique game ID
			break
		}
	}

	gm.mangagedGames[gid] = InitGameState(config)

	return gid
}

func (gm *GameManager) FindGame(gid types.GameID) (*GameState, error) {
	gm.lock.Lock()
	defer gm.lock.Unlock()

	game, ok := gm.mangagedGames[gid]
	if !ok {
		return nil, fmt.Errorf("invalid game id `%s`", gid)
	}
	return game, nil
}

func (gm *GameManager) DeleteGame(gid types.GameID) error {
	gm.lock.Lock()
	defer gm.lock.Unlock()

	if _, ok := gm.mangagedGames[gid]; !ok {
		return fmt.Errorf("invalid game id `%s`", gid)
	}
	delete(gm.mangagedGames, gid)
	return nil
}
