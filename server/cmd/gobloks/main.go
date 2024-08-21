package main

import (
	"flag"
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

func CORSMiddleware(production bool) gin.HandlerFunc {
	var allowedOrigins string
	if production {
		allowedOrigins = "http://209.97.144.150:7777/"
	} else {
		allowedOrigins = "*"
	}

	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", allowedOrigins)
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With, Access-Token")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT")
		c.Writer.Header().Set("Access-Control-Expose-Headers", "Access-Token")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}

func main() {
	isProd := flag.Bool("production", false, "true if running in production")
	flag.Parse()

	if *isProd {
		gin.SetMode(gin.ReleaseMode)
	}

	authorization.SetupKeys()
	globalGameManager := manager.InitGameManager()

	router := gin.Default()
	router.Use(
		ApiMiddleware(globalGameManager),
		CORSMiddleware(*isProd),
		authorization.AuthMiddleware([]gin.HandlerFunc{
			server.CreateGame,
			server.JoinGame,
		}),
	)

	router.POST("/create", server.CreateGame)
	router.POST("/join", server.JoinGame)
	router.PUT("/place", server.PlacePiece)
	router.GET("/ws", server.HandleWebsocket)

	router.Run("0.0.0.0:8888")

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
