package main

import (
	"fmt"
	"gobloks/internal/utilities"
	"time"
)

func main() {

	t1 := time.Now()
	pieces, err := utilities.GeneratePieceSet(8)
	t2 := time.Now()
	fmt.Printf("generated in %v\n", t2.Sub(t1))

	if err != nil {
		fmt.Printf("error: %s\n", err)
	} else {
		fmt.Println(len(pieces))
	}

	players := []utilities.Owner{0, 1, 2, 3}
	board, err := utilities.NewBoard(2, players)

	if err != nil {
		fmt.Printf("%v\n", err)
	} else {
		fmt.Println(board.ToString())
	}

}
