package utilities

import (
	"testing"
)

func TestPointEquality(t *testing.T) {
	pt1 := Point{X: 0, Y: 0}
	pt2 := Point{X: 0, Y: 0}

	expected := true
	result := pt2 == pt1

	if expected != result {
		t.Error("Point equality failed")
	}
}

func TestPointInequality(t *testing.T) {
	pt1 := Point{X: 0, Y: 0}
	pt2 := Point{X: 0, Y: 1}

	expected := false
	result := pt2 == pt1

	if expected != result {
		t.Error("Point inequality failed")
	}
}

func TestPointRotation(t *testing.T) {

	expected := []Point{
		{X: 1, Y: 0},
		{X: 0, Y: 1},
		{X: -1, Y: 0},
		{X: 0, Y: -1},
		{X: 1, Y: 0},
	}

	pt := Point{X: 1, Y: 0}
	for idx, rotation := range []int{0, 90, 180, 270, 360} {
		result := pt.Rotate(rotation)
		if result != expected[idx] {
			t.Errorf("Point rotation failed, expected %+v, got %+v", expected, result)
		}
	}

	pt2 := Point{X: 1, Y: 0}
	for idx := range []int{1, 2, 3, 4} {
		pt2 = pt2.Rotate(90)
		compare := expected[(idx+1)%len(expected)]
		if pt2 != compare {
			t.Errorf("Point rotation failed, expected %+v, got %+v", compare, pt2)
		}
	}
}

func TestPointReflection(t *testing.T) {
	pt := Point{X: 1, Y: 0}

	result := pt.Reflect(X)
	expected := Point{X: 1, Y: 0}

	if expected != result {
		t.Errorf("Point reflection failed, expected %+v, got %+v", expected, result)
	}

	result = pt.Reflect(Y)
	expected = Point{X: -1, Y: 0}
	if expected != result {
		t.Errorf("Point reflection failed, expected %+v, got %+v", expected, result)
	}

	pt = Point{X: 0, Y: 1}

	result = pt.Reflect(X)
	expected = Point{X: 0, Y: -1}
	if expected != result {
		t.Errorf("Point reflection failed, expected %+v, got %+v", expected, result)
	}

	result = pt.Reflect(Y)
	expected = Point{X: 0, Y: 1}
	if expected != result {
		t.Errorf("Point reflection failed, expected %+v, got %+v", expected, result)
	}
}

func TestPointTranslation(t *testing.T) {
	pt := Point{X: 0, Y: 0}

	result := pt.Translate(1, 1)
	expected := Point{X: 1, Y: 1}
	if expected != result {
		t.Errorf("Point translation failed, expected %+v, got %+v", expected, result)
	}

	result = pt.Translate(-1, -1)
	expected = Point{X: -1, Y: -1}
	if expected != result {
		t.Errorf("Point translation failed, expected %+v, got %+v", expected, result)
	}

	result = pt.Translate(-785, 0)
	expected = Point{X: -785, Y: 0}
	if expected != result {
		t.Errorf("Point translation failed, expected %+v, got %+v", expected, result)
	}
}

func TestPointAdjacency(t *testing.T) {
	directions := []Direction{UP, DOWN, RIGHT, LEFT}
	pt := Point{X: 0, Y: 0}

	expected := []Point{
		{X: 0, Y: 1},  // UP
		{X: 0, Y: -1}, // DOWN
		{X: 1, Y: 0},  // RIGHT
		{X: -1, Y: 0}, // LEFT
	}

	for ii, dir := range directions {
		result := pt.GetAdjacent(dir)
		if expected[ii] != result {
			t.Errorf("Point GetAdjacent (Direction = %v) failed, expected %+v, got %+v", dir, expected[ii], result)
		}
	}

}

func TestPointsNormilzation(t *testing.T) {
	var pts, normalized, expected Set[Point]
	pts = NewSet[Point]([]Point{
		{0, 0},
		{0, -1},
	})
	normalized = NormalizeToOrigin(pts)
	expected = NewSet[Point]([]Point{
		{0, 0},
		{0, 1},
	})
	for res := range expected {
		if !normalized.Has(res) {
			t.Errorf("normalization failed, pt %+v not in %+v", res, normalized)
		}
	}

	pts = NewSet[Point]([]Point{
		{0, 0},
		{1, 0},
		{1, -1},
		{1, -2},
	})
	normalized = NormalizeToOrigin(pts)
	expected = NewSet[Point]([]Point{
		{0, 2},
		{1, 2},
		{1, 1},
		{1, 0},
	})
	for res := range expected {
		if !normalized.Has(res) {
			t.Errorf("normalization failed, pt %+v not in %+v", res, normalized)
		}
	}
}

// func TestPointGetAdjacent(t *testing.T) {
// 	directions := []Direction{UP, DOWN, RIGHT, LEFT}

// 	pt1 := Point{X: BoardX / 2, Y: BoardY / 2}
// 	expected1 := []*Point{
// 		{X: (BoardX / 2), Y: (BoardY / 2) + 1},
// 		{X: (BoardX / 2), Y: (BoardY / 2) - 1},
// 		{X: (BoardX / 2) + 1, Y: (BoardY / 2)},
// 		{X: (BoardX / 2) - 1, Y: (BoardY / 2)},
// 	}

// 	pt2 := Point{X: 0, Y: 0}
// 	expected2 := []*Point{
// 		{X: 0, Y: 1},
// 		nil,
// 		{X: 1, Y: 0},
// 		nil,
// 	}

// 	pt3 := Point{X: (BoardX - 1), Y: (BoardY - 1)}
// 	expected3 := []*Point{
// 		nil,
// 		{X: (BoardX - 1), Y: (BoardY - 1) - 1},
// 		nil,
// 		{X: (BoardX - 1) - 1, Y: (BoardY - 1)},
// 	}

// 	pts := []Point{pt1, pt2, pt3}
// 	expectations := [][]*Point{
// 		expected1, expected2, expected3,
// 	}

// 	for idx := 0; idx < len(pts); idx++ {
// 		pt := pts[idx]
// 		exps := expectations[idx]

// 		for i, dir := range directions {
// 			exp := exps[i]
// 			adj, err := pt.GetAdjacent(dir)
// 			if exp != nil && err != nil {
// 				t.Errorf("Point adjacent (%v) returned error: %s", dir, err)
// 			}
// 			if (adj != nil && exp != nil) && (*adj != *exp) {
// 				t.Errorf("Point adjacent (%v) failed, expected %+v, got %+v", dir, *exp, *adj)
// 			}
// 			if (adj == nil) != (exp == nil) {
// 				t.Errorf("Point adjacent (%v) failed, expected %v, got %v", dir, exp, adj)
// 			}
// 		}
// 	}
// }
