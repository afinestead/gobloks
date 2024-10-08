package types

type Direction int
type Axis int
type Owner uint32
type PlayerID uint16
type GameID string

type Flags uint32
type SocketDataType uint32

type SocketData struct {
	Type SocketDataType `json:"type"`
	Data interface{}    `json:"data"`
}

type PublicPiece struct {
	Hash uint64  `json:"hash"`
	Body []Point `json:"body"`
}

type GameConfig struct {
	Players     uint    `json:"players" binding:"required,gte=1,lte=65536"`
	BlockDegree uint8   `json:"degree" binding:"required,gte=1,lte=8"`
	Density     float64 `json:"density"`
	TurnBased   bool    `json:"turns"`
	TimeControl uint    `json:"timeSeconds"`
	TimeBonus   uint    `json:"timeBonus"`
	Hints       uint    `json:"hints"`
}

type PlayerConfig struct {
	PID    PlayerID `json:"pid"`
	Name   string   `json:"name" binding:"required,max=32"`
	Color  uint     `json:"color" binding:"required,gt=0,lte=16777215"`
	Status Flags    `json:"status"`
	Time   uint     `json:"timeMs"`
}

type ChatMessage struct {
	Origin  Owner  `json:"origin"`
	Message string `json:"message"`
}

type PrivateGameState struct {
	PID    PlayerID      `json:"pid"`
	Pieces []PublicPiece `json:"pieces"`
	Hints  uint          `json:"hints"`
}

type PublicGameState struct {
	Turn   PlayerID `json:"turn"`
	Status Flags    `json:"status"`
}

type BoardUpdate struct {
	Owner     `json:"owner"`
	Placement `json:"placement"`
}
