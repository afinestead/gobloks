package server

import (
	"fmt"
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

func handleGameWebsocket(c *gin.Context) {
	g := c.MustGet("manager").(*manager.GameManager)
	gid := c.MustGet("gid").(types.GameID)
	gs, err := g.FindGame(gid)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"message": "access denied"})
		return
	}
	pid := c.MustGet("pid").(types.PlayerID)
	player, err := gs.GetPlayer(pid)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"message": "access denied"})
		return
	}

	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		conn.Close()
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	fmt.Println("Connecting socket PID", pid)
	go gs.ConnectSocket(conn, player)
}

func handleLobbyWebsocket(c *gin.Context) {
	fmt.Println("Connecting lobby socket")
	g := c.MustGet("manager").(*manager.GameManager)

	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		conn.Close()
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	go g.ConnectLobby(conn)
}
