package game

import (
	"errors"
	"fmt"
	"gobloks/internal/types"
	"gobloks/internal/utilities"
	"math"
	"strings"
	"sync"
)

type Board struct {
	layout     [][]types.Owner
	circle     *utilities.Circle
	MaxX, MaxY uint
}

// TODO: Put this somewhere that makes sense
type PlacementInternal struct {
	piece  Piece
	origin types.Point
}

func NewBoard(players []types.PlayerID, pixelsPerPlayer uint, tighteningFactor float64) (*Board, error) {

	numPlayers := uint(len(players))
	if numPlayers == 0 {
		return nil, errors.New("0 length players slice")
	}

	maxOccupiedPixels := numPlayers * pixelsPerPlayer
	// maxOccupancy ~= tighteningFactor * pi * r * r
	radius := uint(math.Sqrt(float64(maxOccupiedPixels) / tighteningFactor / math.Pi))

	diameter := (radius * 2) + 1

	if diameter < radius {
		// overflow
		return nil, errors.New("radius too large")
	}
	board := Board{
		layout: make([][]types.Owner, diameter),
		MaxX:   diameter,
		MaxY:   diameter,
	}
	for i := range board.layout {
		board.layout[i] = make([]types.Owner, diameter)
		for j := range board.layout[i] {
			board.layout[i][j] = types.RESERVED
		}
	}

	circle := utilities.BresenhamCircle(radius, types.Point{X: int(radius), Y: int(radius)})
	board.circle = circle

	// clear out the playable region
	for pt := range circle.Circumference {
		rowSize := circle.Center.X - pt.X
		var offset int
		if rowSize < 0 {
			rowSize = 1 - rowSize
			offset = circle.Center.X
		} else {
			offset = pt.X

		}

		for i := 0; i < rowSize; i++ {
			board.layout[pt.Y][offset+i] = types.VACANT
		}
	}

	angleDelta := (2 * math.Pi) / float64(numPlayers)
	for ii, pid := range players {
		pt := circle.PointOnCircle(angleDelta * float64(ii))
		board.occupy(pt, types.Owner(pid)|types.ORIGIN|types.VACANT)
	}

	return &board, nil
}

func (b *Board) ToString() string {
	var s string
	var ii, jj uint
	s += "-" + strings.Repeat("---", int(b.MaxX)) + "-"
	for ii = 0; ii < b.MaxX; ii++ {
		s += "\n|"
		for jj = 0; jj < b.MaxY; jj++ {
			s += b.layout[jj][ii].ToString()
		}
		s += "|"
	}
	s += "\n-" + strings.Repeat("---", int(b.MaxX)) + "-"
	return s
}

func (b *Board) GetRaw() [][]types.Owner {
	return b.layout
}

func (b *Board) inbounds(square types.Point) bool {
	return (square.X >= 0 && square.X < int(b.MaxX) &&
		square.Y >= 0 && square.Y < int(b.MaxY)) &&
		(b.layout[square.X][square.Y]&types.RESERVED) == 0
}

func (b *Board) occupy(square types.Point, owner types.Owner) error {
	if !b.inbounds(square) {
		return errors.New("out of bounds")
	} else if !b.vacant(square) {
		return errors.New("square is occupied")
	}
	b.layout[square.X][square.Y] = owner
	return nil
}

func (b *Board) vacate(square types.Point) {
	if b.inbounds(square) {
		b.layout[square.X][square.Y] = types.VACANT
	}
}

func (b *Board) owner(square types.Point) (types.Owner, error) {
	if !b.inbounds(square) {
		return types.RESERVED, errors.New("out of bounds")
	}
	return b.layout[square.X][square.Y], nil
}

func (b *Board) vacant(square types.Point) bool {
	sqOwner, err := b.owner(square)
	if err != nil {
		return false
	}
	return sqOwner.IsVacant()
}

func (b *Board) isStartingSquare(square types.Point, owner types.Owner) bool {
	sqOwner, err := b.owner(square)
	if err != nil {
		return false
	}
	return sqOwner.IsOrigin() && sqOwner.IsSamePlayer(owner)
}

func (b *Board) isOriginForOther(square types.Point, owner types.Owner) bool {
	sqOwner, err := b.owner(square)
	if err != nil {
		return false
	}
	return sqOwner.IsOrigin() && !sqOwner.IsSamePlayer(owner)
}

func (b *Board) occupiedByPlayer(square types.Point, owner types.Owner) bool {
	sqOwner, err := b.owner(square)
	if err != nil {
		return false
	}
	return !sqOwner.IsVacant() && sqOwner.IsSamePlayer(owner)
}

func (b *Board) validPlacement(origin types.Point, p Piece, owner types.Owner) bool {
	validCorner := false
	for pt := range p.ToPoints() {
		absPt := origin.Translate(int(pt.X), int(pt.Y))
		if !b.inbounds(absPt) {
			return false
		}
		if b.hasSelfSide(absPt, owner) {
			return false
		}
		if b.isOriginForOther(absPt, owner) {
			return false
		}

		cornerExists := b.hasCorner(absPt, owner)
		startingSq := b.isStartingSquare(absPt, owner)

		validCorner = validCorner || startingSq || cornerExists
	}
	return validCorner
}

func (b *Board) hasSelfSide(pt types.Point, owner types.Owner) bool {
	l := pt.GetAdjacent(types.LEFT)
	r := pt.GetAdjacent(types.RIGHT)
	u := pt.GetAdjacent(types.UP)
	d := pt.GetAdjacent(types.DOWN)

	return ((b.inbounds(l) && b.occupiedByPlayer(l, owner)) ||
		(b.inbounds(r) && b.occupiedByPlayer(r, owner)) ||
		(b.inbounds(u) && b.occupiedByPlayer(u, owner)) ||
		(b.inbounds(d) && b.occupiedByPlayer(d, owner)))
}

func (b *Board) hasCorner(pt types.Point, owner types.Owner) bool {
	ul := pt.GetAdjacent(types.UP).GetAdjacent(types.LEFT)
	ur := pt.GetAdjacent(types.UP).GetAdjacent(types.RIGHT)
	dl := pt.GetAdjacent(types.DOWN).GetAdjacent(types.LEFT)
	dr := pt.GetAdjacent(types.DOWN).GetAdjacent(types.RIGHT)

	return ((b.inbounds(ul) && b.occupiedByPlayer(ul, owner)) ||
		(b.inbounds(ur) && b.occupiedByPlayer(ur, owner)) ||
		(b.inbounds(dl) && b.occupiedByPlayer(dl, owner)) ||
		(b.inbounds(dr) && b.occupiedByPlayer(dr, owner)))
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

func (b *Board) Place(origin types.Point, p Piece, owner types.Owner) (bool, error) {
	valid := b.validPlacement(origin, p, owner)
	if !valid {
		return false, errors.New("invalid placement")
	}
	pieceCoords := p.ToPoints()
	for pt := range pieceCoords {
		abs_pt := origin.Translate(int(pt.X), int(pt.Y))
		err := b.occupy(abs_pt, owner)
		if err != nil {
			for clearPt := range pieceCoords {
				b.vacate(origin.Translate(int(clearPt.X), int(clearPt.Y)))
			}
			return false, err
		}
	}
	return true, nil
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

func (b *Board) getPlacements(owner types.Owner, pieces PieceSet, first bool) utilities.Set[PlacementInternal] {
	// find players corners

	corners := b.findCorners(owner)

	res := utilities.NewSet([]PlacementInternal{})

	tmpPieces := pieces.Copy()

	var wg sync.WaitGroup
	var firstClose sync.Once
	chDone := make(chan struct{})
	chFound := make(chan PlacementInternal)

	placementFinder := func(piece Piece) {
		defer wg.Done()

		for _, pt := range corners {
			attemptedPieces := utilities.NewSet([]Piece{}, 8)
			for j := 0; j < 2; j++ {
				piece.Reflect(types.X)
				for i := 0; i < 4; i++ {
					piece.Rotate90()
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
	}

	for piece := range tmpPieces {
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

func (b *Board) HasPlacement(owner types.Owner, pieces PieceSet) bool {
	plc := b.getPlacements(owner, pieces, true)

	for p := range plc {
		fmt.Println(p.origin, p.piece.ToString())
		// b.Place(p.origin, p.piece, owner)
		// fmt.Println(b.ToString())
		break
		// return true
	}

	// fmt.Println("placements", plc.Size())
	return plc.Size() > 0
}

func (b *Board) GetPlacements(owner types.Owner, pieces PieceSet) utilities.Set[PlacementInternal] {
	return b.getPlacements(owner, pieces, false)
}
