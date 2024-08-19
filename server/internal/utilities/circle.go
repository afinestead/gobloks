package utilities

import (
	"cmp"
	"gobloks/internal/types"
	"math"
	"slices"
)

type Circle struct {
	Circumference Set[types.Point]
	Radius        uint
	Center        types.Point
}

func BresenhamCircle(r uint, c types.Point) *Circle {
	x := 0
	y := int(r)
	h := 1 - int(r)

	circlePix := Set[types.Point]{}

	addOctants := func(px, py int) {
		circlePix.Add(types.Point{X: c.X + px, Y: c.Y + py})
		circlePix.Add(types.Point{X: c.X - px, Y: c.Y + py})
		circlePix.Add(types.Point{X: c.X + px, Y: c.Y - py})
		circlePix.Add(types.Point{X: c.X - px, Y: c.Y - py})
		circlePix.Add(types.Point{X: c.X + py, Y: c.Y + px})
		circlePix.Add(types.Point{X: c.X - py, Y: c.Y + px})
		circlePix.Add(types.Point{X: c.X + py, Y: c.Y - px})
		circlePix.Add(types.Point{X: c.X - py, Y: c.Y - px})
	}

	addOctants(x, y)

	for y > x {
		if h >= 0 {
			h += 2*(x-y) + 5
			y--
		} else {
			h += 2*x + 3
		}
		x++
		addOctants(x, y)
	}
	return &Circle{circlePix, r, c}
}

func (circle *Circle) PointOnCircle(theta float64) types.Point {
	floatX := float64(circle.Center.X) + float64(circle.Radius)*math.Cos(theta)
	floatY := float64(circle.Center.Y) + float64(circle.Radius)*math.Sin(theta)

	point := types.Point{X: int(math.Round(floatX)), Y: int(math.Round(floatY))}
	// search the area around the point to snap to the circle
	for {
		for _, dx := range []int{0, -1, 2} {
			for _, dy := range []int{0, -1, 2} {
				point = types.Point{X: point.X + dx, Y: point.Y + dy}
				if circle.Circumference.Has(point) {
					return point
				}
			}
		}
	}
}

func (circle *Circle) circumference() uint {
	return uint(len(circle.Circumference))
}

func (circle *Circle) area() uint {
	// compute at creation? it's just called once...
	var area uint

	var points []types.Point
	for pt := range circle.Circumference {
		points = append(points, pt)
	}
	slices.SortFunc(
		points,
		func(p1, p2 types.Point) int {
			if n := cmp.Compare(p1.Y, p2.Y); n != 0 {
				return n
			}
			return cmp.Compare(p1.X, p2.X)
		},
	)
	minYPt, maxYPt := points[0], points[0]
	for _, pt := range points {
		if pt.Y != minYPt.Y {
			area += uint(1 + (maxYPt.X - minYPt.X))
			minYPt = pt
		}
		maxYPt = pt
	}
	area += uint(1 + (maxYPt.X - minYPt.X))
	return area
}
