package main

import (
	"flag"
	"fmt"
	"gobloks/internal/manager"
	"gobloks/internal/server"
	"gobloks/internal/types"
)

func main() {
	isProd := flag.Bool("production", false, "true if running in production")
	flag.Parse()
	server.Start(8888, *isProd)

	globalGameManager := manager.InitGameManager()

	gid := globalGameManager.CreateGame(types.GameConfig{
		Players:     1,
		BlockDegree: 5,
		Density:     1,
		TurnBased:   true,
		TimeControl: 0,
		TimeBonus:   0,
		Hints:       0,
	})
	gs, err := globalGameManager.FindGame(gid)
	if err != nil {
		fmt.Printf("error finding game: %s\n", err)
		return
	}
	pid1, err := gs.AddPlayer("p1", 0xff00ff)
	if err != nil {
		fmt.Printf("error connecting player: %s\n", err)
		return
	}
	fmt.Println(pid1)
	// pid2, err := gs.AddPlayer("p2", 0xffff00)
	// if err != nil {
	// 	fmt.Printf("error connecting player: %s\n", err)
	// 	return
	// }
	// fmt.Println(pid2)

	// gs.TurnIter()

	// err = gs.PlacePiece(pid1, []types.Point{{X: 10, Y: 5}, {X: 9, Y: 5}, {X: 8, Y: 5}, {X: 7, Y: 5}, {X: 6, Y: 5}})
	// if err != nil {
	// 	fmt.Printf("error placing piece: %s\n", err)
	// 	return
	// }

	for {
	}

	// plc := []types.Point{{X: 3, Y: 15}, {X: 3, Y: 16}, {X: 3, Y: 17}, {X: 4, Y: 17}, {X: 4, Y: 18}, {X: 5, Y: 18}}}
	// fmt.Println(plc)
	// ptSet := utilities.NewSet(plc)
	// piece := game.PieceFromPoints(ptSet)
	// fmt.Println(piece.ToString())

	// pieceCoords := piece.ToPoints()
	// for pt := range pieceCoords {
	// 	abs_pt := origin.Translate(int(pt.X), int(pt.Y))
	// 	fmt.Println(abs_pt)
	// }

	// err = gs.PlacePiece(pid1, types.Placement{
	// 	[]types.Point{{X: 9, Y: 4}, {X: 9, Y: 6}},
	// })
	// if err != nil {
	// 	fmt.Printf("error placing piece: %s\n", err)
	// 	return
	// }
	// err = gs.PlacePiece(pid1, types.Placement{
	// 	[]types.Point{{X: 14, Y: 7}},
	// })
	// if err != nil {
	// 	fmt.Printf("error placing piece: %s\n", err)
	// 	return
	// }

	// err = gs.PlacePiece(pid2, types.Placement{
	// 	[]types.Point{{X: 0, Y: 7}},
	// })
	// if err != nil {
	// 	fmt.Printf("error placing piece: %s\n", err)
	// 	return
	// }

	// p := game.PieceFromPoints(utilities.NewSet([]game.PieceCoord{{X: 0, Y: 0}, {X: 0, Y: 1}}))
	// p := game.PieceFromPoints(utilities.NewSet([]game.PieceCoord{{X: 0, Y: 0}, {X: 0, Y: 1}, {X: 1, Y: 1}}))
	// p := game.PieceFromPoints(utilities.NewSet([]game.PieceCoord{{X: 0, Y: 0}, {X: 0, Y: 1}, {X: 1, Y: 1}}))
	// p := game.PieceFromPoints(utilities.NewSet([]game.PieceCoord{{X: 1, Y: 0}, {X: 1, Y: 1}, {X: 0, Y: 1}, {X: 2, Y: 1}, {X: 1, Y: 2}}))
	// fmt.Println(p.ToString())
	// fmt.Println(p.Corners())

	// gs.PlacePiece(pid, types.Placement{
	// 	[]types.Point{{X: 9, Y: 6}, {X: 8, Y: 6}},
	// })

	// gs.PlacePiece(pid, types.Placement{
	// 	[]types.Point{{X: 7, Y: 7}, {X: 7, Y: 8}, {X: 7, Y: 9}},
	// })

	// t1 := time.Now()
	// pieces, setPixels, err := game.GeneratePieceSet(5)
	// t2 := time.Now()
	// fmt.Printf("generated in %v\n", t2.Sub(t1))

	// if err != nil {
	// 	fmt.Printf("error: %s\n", err)
	// } else {
	// 	fmt.Println(len(pieces), setPixels)
	// }

	// board, err := game.NewBoard(players, setPixels, 0.9)

	// if err != nil {
	// 	fmt.Printf("%v\n", err)
	// 	return
	// }

	// board.Place(
	// 	types.Point{X: 10, Y: 5},
	// 	game.PieceFromPoints(utilities.NewSet([]game.PieceCoord{{X: 0, Y: 0}})),
	// 	types.Owner(0),
	// )
	// fmt.Println(board.ToString())

}
