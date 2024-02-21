package utilities

const (
	BoardX = 20
	BoardY = 20
)

type Direction int

const (
	UP Direction = iota
	DOWN
	LEFT
	RIGHT
)

type Axis int

const (
	X Axis = iota
	Y
)
