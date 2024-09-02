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
	cleanupChan   chan types.GameID
}

func InitGameManager() *GameManager {
	manager := &GameManager{
		make(map[types.GameID]*GameState, types.MANAGED_GAMES_START_SIZE),
		&sync.Mutex{},
		make(chan types.GameID),
	}

	go func() {
		for {
			cleanupGame := <-manager.cleanupChan
			fmt.Println("cleaned up game", cleanupGame)
			manager.lock.Lock()
			delete(manager.mangagedGames, cleanupGame)
			manager.lock.Unlock()
		}
	}()

	return manager
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

	gm.mangagedGames[gid] = InitGameState(gid, config, &gm.cleanupChan)

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

func (gm *GameManager) ListGames() []types.GameID {
	gm.lock.Lock()
	defer gm.lock.Unlock()

	// TODO: paginate?
	gids := make([]types.GameID, 0, len(gm.mangagedGames))
	for gid := range gm.mangagedGames {
		gids = append(gids, gid)
	}
	return gids
}
