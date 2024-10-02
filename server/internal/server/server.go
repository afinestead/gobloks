package server

import (
	"fmt"
	"gobloks/internal/authorization"
	"gobloks/internal/database"
	"gobloks/internal/manager"
	"os"

	"github.com/gin-gonic/gin"
)

func ApiMiddleware(m *manager.GameManager, db *database.DatabaseManager) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set("manager", m)
		c.Set("db", db)
		c.Next()
	}
}

func CORSMiddleware(production bool) gin.HandlerFunc {
	const FRONTEND_ORIGIN string = "http://209.97.144.150"

	var allowedOrigins string
	if production {
		allowedOrigins = FRONTEND_ORIGIN
	} else {
		allowedOrigins = "*"
	}

	return func(c *gin.Context) {
		if production {
			origin := c.Request.Header.Get("Origin")
			if origin != FRONTEND_ORIGIN {
				c.AbortWithStatus(403)
				return
			}
		}

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

func Start(port uint, production bool) {
	authorization.SetupKeys()

	globalGameManager := manager.InitGameManager()
	db, err := database.ConnectDB()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		return
	}

	if production {
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.Default()
	router.SetTrustedProxies(nil)
	router.Use(
		ApiMiddleware(globalGameManager, db),
		CORSMiddleware(production),
		authorization.AuthMiddleware([]gin.HandlerFunc{
			createGame,
			joinGame,
			listGames,
			handleLobbyWebsocket,
		}),
	)

	// These routes are protected by the auth middleware
	router.PUT("/place", placePiece)
	router.GET("/hint", getHint)
	router.GET("/ws/play", handleGameWebsocket)

	// These routes are not protected by the auth middleware
	router.POST("/create", createGame)
	router.POST("/join", joinGame)
	router.GET("/games", listGames)
	router.GET("/ws/lobby", handleLobbyWebsocket)

	router.Run(fmt.Sprintf("0.0.0.0:%d", port))
}
