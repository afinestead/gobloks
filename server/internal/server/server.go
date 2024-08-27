package server

import (
	"fmt"
	"gobloks/internal/authorization"
	"gobloks/internal/manager"

	"github.com/gin-gonic/gin"
)

func ApiMiddleware(m *manager.GameManager) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set("manager", m)
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

	if production {
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.Default()
	router.Use(
		ApiMiddleware(globalGameManager),
		CORSMiddleware(production),
		authorization.AuthMiddleware([]gin.HandlerFunc{
			createGame,
			joinGame,
		}),
	)

	router.POST("/create", createGame)
	router.POST("/join", joinGame)
	router.PUT("/place", placePiece)
	router.GET("/ws", handleWebsocket)

	router.Run(fmt.Sprintf("0.0.0.0:%d", port))
}
