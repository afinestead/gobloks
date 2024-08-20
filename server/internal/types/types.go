package types

type Direction int
type Axis int
type Owner uint32
type PlayerID uint16
type GameID string
type SocketDataType uint32

type SocketData struct {
	Type SocketDataType `json:"type"`
	Data interface{}    `json:"data"`
}

type Placement struct {
	Coordinates []Point `json:"coords"`
}

type GameConfig struct {
	Players     uint    `json:"players" binding:"required,gte=1,lte=65536"`
	BlockDegree uint8   `json:"degree" binding:"required,gte=1,lte=8"`
	Density     float64 `json:"density"`
}

type PlayerConfig struct {
	PID   PlayerID `json:"pid"`
	Name  string   `json:"name" binding:"required,max=32"`
	Color uint     `json:"color" binding:"required,gt=0,lte=16777215"`
}

type ChatMessage struct {
	Origin  Owner  `json:"origin"`
	Message string `json:"message"`
}

type ActivePlayers struct {
	Players []PlayerConfig `json:"players"`
}

type PublicGameState struct {
	Board [][]Owner `json:"board"`
	Turn  PlayerID  `json:"turn"`
}

type PrivateGameState struct {
	PID    PlayerID  `json:"pid"`
	Pieces [][]Point `json:"pieces"`
}
