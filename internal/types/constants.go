package types

const (
	MaxPieceDegree uint8  = 8
	ColumnMask     uint64 = 0x8080808080808080
	MatrixMagic    uint64 = 0x02040810204081
)

const (
	PLAYER_MASK Owner = 0x0000ffff
	VACANT      Owner = (1 << 29)
	ORIGIN      Owner = (1 << 30)
	RESERVED    Owner = (1 << 31)
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

const MANAGED_GAMES_START_SIZE = 256
