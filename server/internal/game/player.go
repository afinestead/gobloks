package game

import (
	"gobloks/internal/sockets"
	"gobloks/internal/types"
	"gobloks/internal/utilities"
)

// Player status flags (bitset)
const (
	JOINED    types.Flags = (1 << 0) // has joined the game
	CONNECTED types.Flags = (1 << 1) // has socket connection
	DISABLED  types.Flags = (1 << 2) // not playing
	TIMED_OUT types.Flags = (1 << 3) // has timed out
	WINNER    types.Flags = (1 << 4) // has won
	DRAWN     types.Flags = (1 << 5) // has drawn
)

const PID_NONE types.PlayerID = 0

type Player struct {
	name               string
	color              uint
	state              *PlayerState
	socket             *sockets.Connection
	playerTimer        *utilities.Timer
	connectionTimer    *utilities.Timer
	possiblePlacements utilities.LinkedList[utilities.Set[types.Point]]
	hints              uint
}

type PlayerState struct {
	pid    types.PlayerID
	status types.Flags
	pieces PieceSet
}

func (p *PlayerState) Copy() *PlayerState {
	return &PlayerState{
		pid:    p.pid,
		status: p.status,
		pieces: p.pieces.Copy(),
	}
}
