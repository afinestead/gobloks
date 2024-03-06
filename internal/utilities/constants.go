package utilities

import (
	"math"
	"strconv"
)

// Some type definitions
// TODO: Put these in a real spot?
type void struct{} //empty structs occupy 0 memory
type Direction int
type Axis int
type Owner uint32

const (
	MaxPieceDegree uint8  = 8
	ColumnMask     uint64 = 0x8080808080808080
	MatrixMagic    uint64 = 0x02040810204081
)

func (o Owner) ToString() string {
	switch o {
	case VACANT:
		return " "
	case RESERVED:
		return "#"
	}
	return strconv.Itoa(int(o))
}

const (
	VACANT   Owner = math.MaxUint32 - 1
	RESERVED Owner = math.MaxUint32
)

const (
	UP Direction = iota
	DOWN
	LEFT
	RIGHT
)

const (
	X Axis = iota
	Y
)
