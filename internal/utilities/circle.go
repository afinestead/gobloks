package utilities

import "math"

type Circle struct {
	pixels Set[Point]
	radius uint
	center Point
}

func bresenhamCircle(r uint, c Point) *Circle {
	x := 0
	y := int(r)
	h := 1 - int(r)

	circlePix := Set[Point]{}

	addOctants := func(px, py int) {
		circlePix.Add(Point{c.X + px, c.Y + py})
		circlePix.Add(Point{c.X - px, c.Y + py})
		circlePix.Add(Point{c.X + px, c.Y - py})
		circlePix.Add(Point{c.X - px, c.Y - py})
		circlePix.Add(Point{c.X + py, c.Y + px})
		circlePix.Add(Point{c.X - py, c.Y + px})
		circlePix.Add(Point{c.X + py, c.Y - px})
		circlePix.Add(Point{c.X - py, c.Y - px})
	}

	addOctants(x, y)
	for y > x {
		if h >= 0 {
			h += 2*(x-y) + 5 // select SE
			y--
		} else {
			h += 2*x + 3 // Select E
		}
		x++
		addOctants(x, y)
	}
	return &Circle{circlePix, r, c}
}

func (circle *Circle) pointOnCircle(theta float64) Point {
	floatX := float64(circle.center.X) + float64(circle.radius)*math.Cos(theta)
	floatY := float64(circle.center.Y) + float64(circle.radius)*math.Sin(theta)

	point := Point{int(math.Round(floatX)), int(math.Round(floatY))}
	// search the area around the point to snap to the circle
	for {
		for _, dx := range []int{0, -1, 2} {
			for _, dy := range []int{0, -1, 2} {
				point = Point{point.X + dx, point.Y + dy}
				if circle.pixels.Has(point) {
					return point
				}
			}
		}
	}
}

func (circle *Circle) circumference() uint {
	return uint(len(circle.pixels))
}

func (circle *Circle) area() uint {
	return 0
}
