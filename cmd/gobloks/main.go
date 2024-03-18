package main

import (
	"fmt"
	"gobloks/internal/utilities"
	"time"
)

func main() {

	var pieceDegree uint8 = 5
	tighteningFactor := 0.9
	players := []utilities.Owner{0, 1, 2, 3, 4, 5, 7, 8}

	t1 := time.Now()
	pieces, setPixels, err := utilities.GeneratePieceSet(pieceDegree)
	t2 := time.Now()
	fmt.Printf("generated in %v\n", t2.Sub(t1))

	if err != nil {
		fmt.Printf("error: %s\n", err)
	} else {
		fmt.Println(len(pieces), setPixels)
	}

	board, err := utilities.NewBoard(players, setPixels, tighteningFactor)

	if err != nil {
		fmt.Printf("%v\n", err)
	} else {
		fmt.Println(board.ToString())
	}

}
