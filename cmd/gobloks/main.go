package main

import (
	"gobloks/internal/authorization"
	"gobloks/internal/manager"
	"gobloks/internal/server"

	"github.com/gin-gonic/gin"
)

func ApiMiddleware(m *manager.GameManager) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set("manager", m)
		c.Next()
	}
}

func main() {

	authorization.SetupKeys()
	globalGameManager := manager.InitGameManager()

	router := gin.Default()
	router.Use(
		ApiMiddleware(globalGameManager),
		authorization.AuthMiddleware([]gin.HandlerFunc{
			server.CreateGame,
			server.JoinGame,
		}),
	)

	router.POST("/create", server.CreateGame)
	router.POST("/join", server.JoinGame)
	router.GET("/ws", server.HandleWebsocket)

	router.Run("localhost:8888")

	// var pieceDegree uint8 = 5
	// tighteningFactor := 0.9
	// players := []types.Owner{0, 1, 2, 3, 4, 5, 7, 8}

	// t1 := time.Now()
	// pieces, setPixels, err := game.GeneratePieceSet(pieceDegree)
	// t2 := time.Now()
	// fmt.Printf("generated in %v\n", t2.Sub(t1))

	// if err != nil {
	// 	fmt.Printf("error: %s\n", err)
	// } else {
	// 	fmt.Println(len(pieces), setPixels)
	// }

	// board, err := utilities.NewBoard(players, setPixels, tighteningFactor)

	// if err != nil {
	// 	fmt.Printf("%v\n", err)
	// 	return
	// }

	// board.Place(
	// 	utilities.Point{X: 29, Y: 15},
	// 	utilities.PieceFromPoints(utilities.NewSet([]utilities.PieceCoord{{X: 0, Y: 1}, {X: 0, Y: 0}, {X: 1, Y: 0}, {X: 1, Y: 1}})),
	// 	utilities.Owner(0),
	// )
	// fmt.Println(board.ToString())

}
