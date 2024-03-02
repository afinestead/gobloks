package utilities

// Some type definitions
// TODO: Put these in a real spot?
type void struct{} //empty structs occupy 0 memory
type Direction int
type Axis int
type Owner uint

const (
	BoardX         uint32 = 20
	BoardY         uint32 = 20
	MaxPieceDegree uint8  = 8
	ColumnMask     uint64 = 0x8080808080808080
	MatrixMagic    uint64 = 0x02040810204081
)

const UNOWNED Owner = 0

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
