package main

import (
	"gobloks/internal/authorization"
	"gobloks/internal/manager"
	"gobloks/internal/server"
	"log"
	"net/http"
	"os"
	"reflect"

	"github.com/gin-gonic/gin"
)

func ApiMiddleware(m *manager.GameManager) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set("manager", m)
		c.Next()
	}
}

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// go can't compare function pointers??
		if reflect.ValueOf(c.Handler()).Pointer() == reflect.ValueOf(server.CreateGame).Pointer() ||
			reflect.ValueOf(c.Handler()).Pointer() == reflect.ValueOf(server.JoinGame).Pointer() {
			// endpoints don't need authentication
			c.Next()
			return
		}

		token := c.GetHeader("access_token")
		_, err := authorization.VerifyAccessToken(token)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": "token invalid"})
			return
		}
		c.Next()
	}
}

func main() {

	if _, err := os.ReadFile("key.priv"); err != nil {
		log.Println("Writing new key into file...")
		file, err := os.OpenFile("key.priv", os.O_WRONLY|os.O_CREATE, 0644)
		if err != nil {
			log.Fatal(err)
		}
		defer file.Close()

		newKey, err := authorization.GenerateKey()
		if err != nil {
			log.Fatal(err)
		}

		_, err = file.WriteString(newKey)
		if err != nil {
			log.Fatal(err)
		}
	}
	globalGameManager := manager.InitGameManager()

	router := gin.Default()
	router.Use(
		ApiMiddleware(globalGameManager),
		AuthMiddleware(),
	)

	router.POST("/create", server.CreateGame)
	router.POST("/join", server.JoinGame)
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
