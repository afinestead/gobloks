package types

// Some type definitions
// TODO: Put these in a real spot?
type Direction int
type Axis int
type Owner uint32
type PlayerID uint16
type GameID string

type GameConfig struct {
	Players     uint    `json:"players" binding:"required,gte=1,lte=65536"`
	BlockDegree uint8   `json:"degree" binding:"required,gte=1,lte=8"`
	Density     float64 `json:"density"`
}

type PlayerConfig struct {
	Name  string `json:"name" binding:"required,max=32"`
	Color uint   `json:"color" binding:"required,gt=0,lte=16777215"`
}
