package manager

import (
	"errors"
	"fmt"

	"gobloks/internal/database"
	"gobloks/internal/game"
	"gobloks/internal/sockets"
	"gobloks/internal/types"
	"sync"

	"github.com/gorilla/websocket"
)

type GameManager struct {
	mangagedGames      map[types.GameID]*game.Game
	lock               *sync.Mutex
	lobbySocketManager *sockets.SocketManager
}

func InitGameManager() *GameManager {
	manager := &GameManager{
		make(map[types.GameID]*game.Game, types.MANAGED_GAMES_START_SIZE),
		&sync.Mutex{},
		sockets.InitSocketManager(64), // TODO: Test that this will dynamically resize
	}

	// go func() {
	// 	for range time.Tick(time.Hour * 24) {
	// 		manager.CleanupStale()
	// 	}
	// }()

	return manager
}

// const letterBytes = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"

// func createGameID(n uint8) types.GameID {
// 	b := make([]byte, n)
// 	for i := range b {
// 		b[i] = letterBytes[rand.Intn(len(letterBytes))]
// 	}
// 	return types.GameID(b)
// }

func (gm *GameManager) CreateGame(
	config *types.GameConfig,
	db *database.DatabaseManager,
) (types.GameID, error) {

	g, err := game.InitGame(config, db)
	if err != nil {
		return 0, err
	}
	gm.mangagedGames[g.GID] = g
	return g.GID, nil
}

func (gm *GameManager) FindGame(gid types.GameID) (*game.Game, error) {
	gm.lock.Lock()
	defer gm.lock.Unlock()

	maybeGame := gm.mangagedGames[gid]
	if maybeGame == nil {
		return nil, errors.New("game not found")
	}
	return maybeGame, nil
}

func (gm *GameManager) ListGames(activeOnly bool, page, pageSize int) []types.GameID {
	gm.lock.Lock()
	defer gm.lock.Unlock()

	// TODO: paginate
	gids := make([]types.GameID, 0, len(gm.mangagedGames))
	for gid := range gm.mangagedGames {
		gids = append(gids, gid)
	}
	return gids
}

func (gm *GameManager) ConnectLobby(socket *websocket.Conn) {
	conn := gm.lobbySocketManager.Connect(socket)
	defer gm.lobbySocketManager.Disconnect(conn)

	for {
		var inMsg types.SocketData
		err := gm.lobbySocketManager.Recv(conn, &inMsg)
		if err != nil {
			fmt.Println(err)
			return
		}
	}
}

func (gm *GameManager) SendLobbyUpdate(gid types.GameID, players uint) {
	gm.lobbySocketManager.Broadcast(&types.SocketData{
		Type: sockets.LOBBY_UPDATE,
		Data: &types.LobbyUpdate{
			GID:     gid,
			Players: players,
		},
	})
}

// func (gm *GameManager) CleanupStale() {
// 	gm.lock.Lock()
// 	defer gm.lock.Unlock()
// 	fmt.Println("cleaning up stale games")
// 	for gid := range gm.mangagedGames {
// 		if gm.mangagedGames[gid].IsStale() {
// 			fmt.Println("cleaned up stale game", gid)
// 			delete(gm.mangagedGames, gid)
// 		}
// 	}
// }
