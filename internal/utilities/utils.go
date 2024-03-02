package utilities

import (
	"errors"
	"math"
)

type Point struct {
	X int `json:"x"`
	Y int `json:"y"`
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

func (pt Point) Translate(x int, y int) Point {
	return Point{pt.X + x, pt.Y + y}
}

func Translate(points Set[Point], x int, y int) Set[Point] {
	translated := Set[Point]{}
	for pt := range points {
		translated.Add(pt.Translate(x, y))
	}
	return translated
}

func (pt Point) Rotate(degrees int) Point {
	rad := float64(degrees) * (math.Pi / 180)
	cos := float32(math.Cos(rad))
	sin := float32(math.Sin(rad))

	newX := int((float32(pt.X) * cos) - (float32(pt.Y) * sin))
	newY := int((float32(pt.Y) * cos) + (float32(pt.X) * sin))
	return Point{newX, newY}
}

// Rotate a set of points
func Rotate(points Set[Point], degrees int) Set[Point] {
	rad := float64(degrees) * (math.Pi / 180)
	cos := float32(math.Cos(rad))
	sin := float32(math.Sin(rad))

	rotated := Set[Point]{}
	for pt := range points {
		newX := int((float32(pt.X) * cos) - (float32(pt.Y) * sin))
		newY := int((float32(pt.Y) * cos) + (float32(pt.X) * sin))
		rotated.Add(Point{newX, newY})
	}
	return NormalizeToOrigin(rotated)
}

func (pt Point) Reflect(ax Axis) Point {
	if ax == X {
		return Point{pt.X, -pt.Y}
	} else {
		return Point{-pt.X, pt.Y}
	}
}

// Reflect a set of points across an axis
func Reflect(points Set[Point], ax Axis) Set[Point] {
	reflected := Set[Point]{}
	for pt := range points {
		reflected.Add(pt.Reflect(ax))
	}
	return NormalizeToOrigin(reflected)
}

func NormalizeToOrigin(points Set[Point]) Set[Point] {

	minCoordinate := func() (int, int, error) {
		if points.Size() == 0 {
			return 0, 0, errors.New("cannot compute min on 0 size piece")
		}
		minX, minY := math.MaxInt, math.MaxInt
		for pt := range points {
			if pt.X < minX {
				minX = pt.X
			}
			if pt.Y < minY {
				minY = pt.Y
			}
		}
		return minX, minY, nil
	}

	minX, minY, err := minCoordinate()

	if err != nil {
		return points // 0 size piece
	}

	return Translate(points, -minX, -minY)
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
