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

}
