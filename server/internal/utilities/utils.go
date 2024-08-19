package utilities

import (
	"errors"
	"gobloks/internal/types"
	"math"
)

// Translate a set of points
func Translate(points Set[types.Point], x int, y int) Set[types.Point] {
	translated := Set[types.Point]{}
	for pt := range points {
		translated.Add(pt.Translate(x, y))
	}
	return translated
}

// Rotate a set of points
func Rotate(points Set[types.Point], degrees int) Set[types.Point] {
	rad := float64(degrees) * (math.Pi / 180)
	cos := float32(math.Cos(rad))
	sin := float32(math.Sin(rad))

	rotated := Set[types.Point]{}
	for pt := range points {
		newX := int((float32(pt.X) * cos) - (float32(pt.Y) * sin))
		newY := int((float32(pt.Y) * cos) + (float32(pt.X) * sin))
		rotated.Add(types.Point{X: newX, Y: newY})
	}
	return NormalizeToOrigin(rotated)
}

// Reflect a set of points across an axis
func Reflect(points Set[types.Point], ax types.Axis) Set[types.Point] {
	reflected := Set[types.Point]{}
	for pt := range points {
		reflected.Add(pt.Reflect(ax))
	}
	return NormalizeToOrigin(reflected)
}

func NormalizeToOrigin(points Set[types.Point]) Set[types.Point] {

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

// {
// 	int x = 0;
// 	int y = radius ;
// 	int h = 1 – radius ;
// 	CirclePoints(x, y, value);
// 	while (y > x) {
// 	if (h < 0) { /* Select E */
// 		h = h + 2 * x + 3;
// 	}
// 	else { /* Select SE */
// 		h = h + 2 * ( x – y ) + 5;
// 		y = y – 1;
// 	}
// 	x = x + 1;
// 	CirclePoints(x, y);
// 	}
// 	}
