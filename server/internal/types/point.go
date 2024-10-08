package types

import (
	"math"
)

type Point struct {
	X int `json:"x"`
	Y int `json:"y"`
}

type Placement []Point

func (p Point) GetAdjacent(dir Direction) Point {
	if dir == UP {
		return Point{p.X, p.Y + 1}
	} else if dir == DOWN {
		return Point{p.X, p.Y - 1}
	} else if dir == LEFT {
		return Point{p.X - 1, p.Y}
	} else { // dir == types.RIGHT
		return Point{p.X + 1, p.Y}
	}
}

func (pt Point) Translate(x int, y int) Point {
	return Point{pt.X + x, pt.Y + y}
}

func (pt Point) Rotate(degrees int) Point {
	rad := float64(degrees) * (math.Pi / 180)
	cos := float32(math.Cos(rad))
	sin := float32(math.Sin(rad))

	newX := int((float32(pt.X) * cos) - (float32(pt.Y) * sin))
	newY := int((float32(pt.Y) * cos) + (float32(pt.X) * sin))
	return Point{newX, newY}
}

func (pt Point) Reflect(ax Axis) Point {
	if ax == X {
		return Point{pt.X, -pt.Y}
	} else {
		return Point{-pt.X, pt.Y}
	}
}

func (pt Point) Is(other Point) bool {
	return pt.X == other.X && pt.Y == other.Y
}
