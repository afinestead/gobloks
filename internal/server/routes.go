package server

import (
	"fmt"
	"gobloks/internal/authorization"
	"gobloks/internal/manager"
	"gobloks/internal/types"
	"net/http"

	"github.com/gin-gonic/gin"
)

func CreateGame(c *gin.Context) {
	var config types.GameConfig

	err := c.BindJSON(&config)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, err)
		return
	}

	gm := c.MustGet("manager").(*manager.GameManager)
	gid := gm.CreateGame(config)

	c.IndentedJSON(http.StatusCreated, gid)
}

func JoinGame(c *gin.Context) {
	gid, ok := c.GetQuery("game")
	if !ok {
		c.AbortWithStatusJSON(http.StatusNotFound, "no game provided")
		return
	}

	var config types.PlayerConfig
	err := c.BindJSON(&config)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, err)
		return
	}

	gm := c.MustGet("manager").(*manager.GameManager)
	gs, err := gm.FindGame(types.GameID(gid))
	if err != nil {
		c.AbortWithStatusJSON(http.StatusNotFound, fmt.Sprintf("no game %s", gid))
		return
	}

	pid, err := gs.ConnectPlayer(config.Name, config.Color)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, err)
		return
	}

	token, err := authorization.CreateAccessToken(pid, types.GameID(gid), 3600)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, err)
		return
	}

	c.Writer.Header().Set("access_token", token)
}

func PlacePiece(c *gin.Context) {
	fmt.Println("placing")
}
