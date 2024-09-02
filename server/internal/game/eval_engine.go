package game

import (
	"fmt"
	"gobloks/internal/types"
	"gobloks/internal/utilities"
	"sync"
)

func (b *Board) getPlacements(owner types.Owner, pieces PieceSet, first bool) utilities.Set[PlacementInternal] {
	corners := b.findCorners(owner)

	res := utilities.NewSet([]PlacementInternal{})

	var wg sync.WaitGroup
	var firstClose sync.Once
	chDone := make(chan struct{})
	chFound := make(chan PlacementInternal)

	findPiecePlacement := func(pt types.Point, piece Piece) {
		attemptedPieces := utilities.NewSet([]Piece{}, 8)
		for j := 0; j < 2; j++ {
			piece = piece.Reflect(types.X)
			for i := 0; i < 4; i++ {
				piece = piece.Rotate90()
				if attemptedPieces.Has(piece) {
					continue // Already tried this one
				}
				attemptedPieces.Add(piece)

				// Check all possible placements
				// TODO: Optimize?
				piecePoints := piece.ToPoints()
				for pp := range piecePoints {

					absPieceOrigin := pt.Translate(-int(pp.X), -int(pp.Y))
					if b.validPlacement(absPieceOrigin, piece, owner) {
						select {
						case <-chDone:
							return
						default:
							chFound <- PlacementInternal{piece, absPieceOrigin}
							if first {
								firstClose.Do(func() { close(chDone) })
							}
						}
					}
				}
			}
		}
	}

	playerOrigin := b.origins[types.PlayerID(owner)]

	originPlacementFinder := func(piece Piece) {
		defer wg.Done()
		findPiecePlacement(playerOrigin, piece)
	}

	cornerPlacementFinder := func(piece Piece) {
		defer wg.Done()
		for _, pt := range corners {
			findPiecePlacement(pt, piece)
		}
	}

	var placementFinder func(piece Piece)
	if b.occupiedByPlayer(playerOrigin, owner) {
		placementFinder = cornerPlacementFinder
	} else {
		placementFinder = originPlacementFinder
	}

	for piece := range pieces {
		wg.Add(1)
		go placementFinder(piece)
	}

	var resultGroup sync.WaitGroup
	resultGroup.Add(1)
	go func() {
		defer resultGroup.Done()
		for found := range chFound {
			res.Add(found)
		}
	}()

	wg.Wait()
	close(chFound)
	resultGroup.Wait()

	return res
}

func (b *Board) findTerritory(o types.Owner) []types.Point {
	territory := make([]types.Point, 0, b.MaxX*b.MaxY)
	for ii := 0; ii < int(b.MaxX); ii++ {
		for jj := 0; jj < int(b.MaxY); jj++ {
			if b.occupiedByPlayer(types.Point{X: ii, Y: jj}, o) {
				territory = append(territory, types.Point{X: ii, Y: jj})
			}
		}
	}
	return territory
}

func (b *Board) findCorners(owner types.Owner) []types.Point {
	territory := b.findTerritory(owner)
	fmt.Println("territory ", territory)

	corners := make([]types.Point, 0, b.MaxX*b.MaxY)
	for _, pt := range territory {
		corners = append(corners, b.getFreeCorners(pt, owner)...)
	}

	fmt.Println("corners ", corners)
	return corners
}

func (b *Board) getFreeCorners(pt types.Point, owner types.Owner) []types.Point {
	l := pt.GetAdjacent(types.LEFT)
	r := pt.GetAdjacent(types.RIGHT)
	u := pt.GetAdjacent(types.UP)
	d := pt.GetAdjacent(types.DOWN)

	ul := u.GetAdjacent(types.LEFT)
	ur := u.GetAdjacent(types.RIGHT)
	dl := d.GetAdjacent(types.LEFT)
	dr := d.GetAdjacent(types.RIGHT)

	vacancies := make([]types.Point, 0, 4)

	if b.inbounds(ul) && b.vacant(ul) && !b.occupiedByPlayer(u, owner) && !b.occupiedByPlayer(l, owner) {
		vacancies = append(vacancies, ul)
	}
	if b.inbounds(ur) && b.vacant(ur) && !b.occupiedByPlayer(u, owner) && !b.occupiedByPlayer(r, owner) {
		vacancies = append(vacancies, ur)
	}
	if b.inbounds(dl) && b.vacant(dl) && !b.occupiedByPlayer(d, owner) && !b.occupiedByPlayer(l, owner) {
		vacancies = append(vacancies, dl)
	}
	if b.inbounds(dr) && b.vacant(dr) && !b.occupiedByPlayer(d, owner) && !b.occupiedByPlayer(r, owner) {
		vacancies = append(vacancies, dr)
	}
	return vacancies
}
