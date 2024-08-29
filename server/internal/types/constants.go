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

const PID_NONE PlayerID = 0

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

const (
	PLAYER_UPDATE SocketDataType = iota
	PUBLIC_GAME_STATE
	PRIVATE_GAME_STATE
	CHAT_MESSAGE
)

const MANAGED_GAMES_START_SIZE = 256
