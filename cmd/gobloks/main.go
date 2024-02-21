package main

import (
	"fmt"
	"gobloks/internal/utilities"
	"math"
	"slices"
)

func main() {

	pieces := utilities.GeneratePieceSet(5)

	fmt.Printf("Set len: %v\n", len(pieces))

	type CumSum struct {
		x int
		y int
	}

	cumsum := make(map[CumSum]utilities.Piece)

	for _, piece := range pieces {
		cumX, cumY := 0, 0
		for pt := range piece {
			cumX += pt.X
			cumY += pt.Y
		}
		if cumsum[CumSum{cumX, cumY}] != nil {
			fmt.Println(piece)
			fmt.Println(cumsum[CumSum{cumX, cumY}])
			fmt.Printf("X: %v\tY: %v\n", cumX, cumY)
		}
		cumsum[CumSum{cumX, cumY}] = piece
	}

	fmt.Println()
	fmt.Println("=======================")
	fmt.Println()

	for _, rot := range []int{0, 90, 180, 270} {
		distList := make([]float64, 0, len(pieces))
		for _, piece := range pieces {
			piece = piece.Rotate(rot).NormalizeToOrigin()
			cumDist := 0.0
			for pt := range piece {
				cumDist += math.Sqrt(math.Pow(float64(pt.X), 2) + math.Pow(float64(pt.Y), 2))
			}
			distList = append(distList, cumDist)
		}
		fmt.Println()
		slices.Sort(distList)
		fmt.Println(rot, distList)
	}

}
