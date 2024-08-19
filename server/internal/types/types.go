package types

type Direction int
type Axis int
type Owner uint32
type PlayerID uint16
type GameID string

type Coordinate struct{ X, Y uint8 }

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

type PublicGameState struct {
	Board   [][]Owner      `json:"board"`
	Turn    PlayerID       `json:"turn"`
	Players []PlayerConfig `json:"players"`
}

type PrivateGameState struct {
	PID    PlayerID  `json:"pid"`
	Pieces [][]Point `json:"pieces"`
}
