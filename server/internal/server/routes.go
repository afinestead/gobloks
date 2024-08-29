package server

import (
	"fmt"
	"gobloks/internal/authorization"
	"gobloks/internal/manager"
	"gobloks/internal/types"
	"net/http"

	"github.com/gin-gonic/gin"
)

func createGame(c *gin.Context) {
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

func joinGame(c *gin.Context) {
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

	pid, err := gs.AddPlayer(config.Name, config.Color)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusConflict, err)
		return
	}

	token, err := authorization.CreateAccessToken(pid, types.GameID(gid), 3600)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, err)
		return
	}

	c.Writer.Header().Set("Access-Token", token)
}

func placePiece(c *gin.Context) {
	fmt.Println("placing")

	g := c.MustGet("manager").(*manager.GameManager)
	gid := c.MustGet("gid").(types.GameID)
	gs, err := g.FindGame(gid)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusNotFound, fmt.Sprintf("no game %s", gid))
		return
	}

	pid := c.MustGet("pid").(types.PlayerID)

	var placement types.Placement
	err = c.BindJSON(&placement)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, err)
		return
	}

	err = gs.PlacePiece(pid, placement)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusConflict, "invalid placement")
		return
	}
}
