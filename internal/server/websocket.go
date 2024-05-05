package server

import (
	"gobloks/internal/manager"
	"gobloks/internal/types"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		// Allow all connections
		return true
	},
}

func HandleWebsocket(c *gin.Context) {
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	g := c.MustGet("manager").(*manager.GameManager)
	gid := c.MustGet("gid").(types.GameID)
	gs, err := g.FindGame(gid)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "token invalid"})
		return
	}

	pid := c.MustGet("pid").(types.PlayerID)

	go gs.ConnectSocket(conn, pid)
}
