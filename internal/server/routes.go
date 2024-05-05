package server

import (
	"fmt"
	"gobloks/internal/authorization"
	"gobloks/internal/manager"
	"gobloks/internal/types"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func CreateGame(c *gin.Context) {
	var config types.GameConfig

	err := c.BindJSON(&config)
	if err != nil {
		fmt.Println(err)
		return
	}

	gm := c.MustGet("manager").(*manager.GameManager)
	gid := gm.CreateGame(config)

	c.IndentedJSON(http.StatusCreated, gid)
}

func JoinGame(c *gin.Context) {
	fmt.Println("joining")

	gid, ok := c.GetQuery("gid")
	if !ok {
		fmt.Println("no gid")
		return
	}

	gm := c.MustGet("manager").(*manager.GameManager)
	gs, err := gm.FindGame(types.GameID(gid))
	if err != nil {
		fmt.Println("no game")
		return
	}

	name, ok := c.GetQuery("name")
	if !ok {
		fmt.Println("no name")
		return
	}
	colorStr, ok := c.GetQuery("color")
	if !ok {
		fmt.Println("no color")
		return
	}
	colorInt, err := strconv.Atoi(colorStr)
	if err != nil || uint(colorInt) > 0xffffff {
		fmt.Println("invalid color")
		return
	}

	pid, err := gs.ConnectPlayer(name, uint(colorInt))
	if err != nil {
		fmt.Println(err)
		return
	}

	authorization.CreateAccessToken(pid, types.GameID(gid), 3600)
}

func PlacePiece(c *gin.Context) {
	fmt.Println("placing")
}
