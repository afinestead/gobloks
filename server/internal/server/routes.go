package server

import (
	"fmt"
	"gobloks/internal/authorization"
	"gobloks/internal/database"
	"gobloks/internal/manager"
	"gobloks/internal/types"
	"net/http"
	"strconv"

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
	db := c.MustGet("db").(*database.DatabaseManager)
	gid, err := gm.CreateGame(&config, db)

	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, err)
		return
	}

	c.IndentedJSON(http.StatusCreated, gid)
}

func listGames(c *gin.Context) {
	fmt.Println("list games")
	pageStr, ok := c.GetQuery("page")
	if !ok {
		c.AbortWithStatusJSON(http.StatusBadRequest, "no page provided")
		return
	}
	page, err := strconv.ParseUint(pageStr, 10, 64)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, err)
		return
	}

	db := c.MustGet("db").(*database.DatabaseManager)
	games, err := db.GetWaitingGames(50, page)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, err)
		return
	}

	c.IndentedJSON(http.StatusOK, games)
}

func joinGame(c *gin.Context) {
	gidStr, ok := c.GetQuery("game")
	if !ok {
		c.AbortWithStatusJSON(http.StatusNotFound, "no game provided")
		return
	}

	gid, err := strconv.ParseUint(gidStr, 10, 64)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, err)
		return
	}

	var config types.PlayerConfig
	err = c.BindJSON(&config)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, err)
		return
	}

	gm := c.MustGet("manager").(*manager.GameManager)
	gs, err := gm.FindGame(types.GameID(gid))
	if err != nil {
		c.AbortWithStatusJSON(http.StatusNotFound, fmt.Sprintf("no game %d", gid))
		return
	}

	pid, count, err := gs.AddPlayer(config.Name, config.Color)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusConflict, err)
		return
	}

	token, err := authorization.CreateAccessToken(pid, types.GameID(gid), 3600)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, err)
		return
	}

	gm.SendLobbyUpdate(types.GameID(gid), count)

	c.Writer.Header().Set("Access-Token", token)
}

func placePiece(c *gin.Context) {
	g := c.MustGet("manager").(*manager.GameManager)
	gid := c.MustGet("gid").(types.GameID)
	gs, err := g.FindGame(gid)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusNotFound, fmt.Sprintf("no game %d", gid))
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
		fmt.Println(err)
		c.AbortWithStatusJSON(http.StatusConflict, "invalid placement")
		return
	}
}

func getHint(c *gin.Context) {
	g := c.MustGet("manager").(*manager.GameManager)
	gid := c.MustGet("gid").(types.GameID)
	gs, err := g.FindGame(gid)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusNotFound, fmt.Sprintf("no game %d", gid))
		return
	}

	pid := c.MustGet("pid").(types.PlayerID)
	hint, err := gs.GetHint(pid)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusConflict, err)
		return
	}

	c.IndentedJSON(http.StatusOK, hint)
}
