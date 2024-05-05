package utilities

import (
	"gobloks/internal/types"
	"math"
	"testing"
)

func TestPointOnCircle(t *testing.T) {
	for _, radius := range []int{4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20} {
		circle := BresenhamCircle(uint(radius), types.Point{X: radius, Y: radius})
		deg := 0.0
		for deg < (math.Pi * 2) {
			testPt := circle.PointOnCircle(deg)
			if !circle.Circumference.Has(testPt) {
				t.Errorf("error on theta %.2f: circle of radius %v does not contain point %+v", deg, circle.Radius, testPt)
			}
			deg += 0.01
		}
	}
}
