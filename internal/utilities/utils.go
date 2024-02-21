package utilities

import (
	"math"
)

type void struct{} //empty structs occupy 0 memory

type Point struct {
	X int `json:"x"`
	Y int `json:"y"`
}

type Owner struct {
	PID int `json:"pid"`
}

func (p Point) GetAdjacent(dir Direction) Point {
	if dir == UP {
		return Point{p.X, p.Y + 1}
	} else if dir == DOWN {
		return Point{p.X, p.Y - 1}
	} else if dir == LEFT {
		return Point{p.X - 1, p.Y}
	} else { // dir == RIGHT
		return Point{p.X + 1, p.Y}
	}
}

func (p Point) Rotate(degrees int) Point {
	rad := float64(degrees) * (math.Pi / 180)
	cos := float32(math.Cos(rad))
	sin := float32(math.Sin(rad))

	newX := int((float32(p.X) * cos) - (float32(p.Y) * sin))
	newY := int((float32(p.Y) * cos) + (float32(p.X) * sin))
	return Point{newX, newY}
}

func (p Point) Reflect(ax Axis) Point {
	if ax == X {
		return Point{p.X, -p.Y}
	} else {
		return Point{-p.X, p.Y}
	}
}

func (p Point) Translate(x int, y int) Point {
	return Point{p.X + x, p.Y + y}
}

type Set[T comparable] map[T]void

func NewSet[T comparable](items []T) Set[T] {
	s := make(Set[T])
	for _, item := range items {
		s.Add(item)
	}
	return s
}

func (s Set[T]) Has(v T) bool {
	_, ok := s[v]
	return ok
}

func (s Set[T]) Add(v T) {
	s[v] = void{}
}

func (s Set[T]) Remove(v T) {
	delete(s, v)
}

func (s Set[T]) Clear() {
	s = make(map[T]void)
}

func (s Set[T]) Size() int {
	return len(s)
}

func (s1 Set[T]) Is(s2 Set[T]) bool {
	if s1.Size() != s2.Size() {
		return false
	}
	for elem := range s1 {
		if !s2.Has(elem) {
			return false
		}
	}
	return true
}
