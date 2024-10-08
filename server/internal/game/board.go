package game

import (
	"errors"
	"gobloks/internal/types"
	"gobloks/internal/utilities"
	"math"
	"strings"
	"sync"
)

type Board struct {
	layout     [][]types.Owner
	origins    map[types.PlayerID]types.Point
	maxX, maxY uint
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
		layout:  make([][]types.Owner, diameter),
		origins: make(map[types.PlayerID]types.Point, numPlayers),
		maxX:    diameter,
		maxY:    diameter,
	}
	for i := range board.layout {
		board.layout[i] = make([]types.Owner, diameter)
		for j := range board.layout[i] {
			board.layout[i][j] = types.RESERVED
		}
	}

	circle := utilities.BresenhamCircle(radius, types.Point{X: int(radius), Y: int(radius)})

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
		board.origins[pid] = pt
	}

	return &board, nil
}

func (b *Board) Copy() *Board {
	boardCopy := Board{
		layout:  make([][]types.Owner, b.maxX),
		origins: make(map[types.PlayerID]types.Point, len(b.origins)),
		maxX:    b.maxX,
		maxY:    b.maxY,
	}
	for i := range b.layout {
		boardCopy.layout[i] = make([]types.Owner, b.maxY)
		copy(boardCopy.layout[i], b.layout[i])
	}
	for pid, pt := range b.origins {
		boardCopy.origins[pid] = pt
	}
	return &boardCopy
}

func (b *Board) ToString() string {
	var s string
	var ii, jj uint
	s += "-" + strings.Repeat("---", int(b.maxX)) + "-"
	for ii = 0; ii < b.maxX; ii++ {
		s += "\n|"
		for jj = 0; jj < b.maxY; jj++ {
			s += b.layout[jj][ii].ToString()
		}
		s += "|"
	}
	s += "\n-" + strings.Repeat("---", int(b.maxX)) + "-"
	return s
}

func (b *Board) GetRaw() [][]types.Owner {
	return b.layout
}

func (b *Board) inbounds(square types.Point) bool {
	return (square.X >= 0 && square.X < int(b.maxX) &&
		square.Y >= 0 && square.Y < int(b.maxY)) &&
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

func (b *Board) validPlacement(points utilities.Set[types.Point], owner types.Owner) bool {
	valid := false
	for pt := range points {
		if !b.inbounds(pt) || !b.vacant(pt) || b.hasSelfSide(pt, owner) || b.isOriginForOther(pt, owner) {
			return false
		}
		valid = valid || b.isStartingSquare(pt, owner) || b.HasCorner(pt, owner)
	}
	return valid
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

func (b *Board) HasCorner(pt types.Point, owner types.Owner) bool {
	ul := pt.GetAdjacent(types.UP).GetAdjacent(types.LEFT)
	ur := pt.GetAdjacent(types.UP).GetAdjacent(types.RIGHT)
	dl := pt.GetAdjacent(types.DOWN).GetAdjacent(types.LEFT)
	dr := pt.GetAdjacent(types.DOWN).GetAdjacent(types.RIGHT)

	return ((b.inbounds(ul) && b.occupiedByPlayer(ul, owner)) ||
		(b.inbounds(ur) && b.occupiedByPlayer(ur, owner)) ||
		(b.inbounds(dl) && b.occupiedByPlayer(dl, owner)) ||
		(b.inbounds(dr) && b.occupiedByPlayer(dr, owner)))
}

func (b *Board) Place(points utilities.Set[types.Point], player types.PlayerID) (bool, error) {
	// valid := b.validPlacement(points, types.Owner(player)) // technically already checked validity...
	// if !valid {
	// 	return false, errors.New("invalid placement")
	// }
	for pt := range points {
		err := b.occupy(pt, types.Owner(player))
		if err != nil {
			for clearPt := range points {
				b.vacate(clearPt)
			}
			return false, err
		}
	}
	return true, nil
}

func (b *Board) getOrigin(player types.PlayerID) types.Point {
	return b.origins[player]
}

func (b *Board) getPlacementsAtPoint(
	pt types.Point,
	owner types.Owner,
	pieces PieceSet,
) utilities.LinkedList[utilities.Set[types.Point]] {

	res := &utilities.Node[utilities.Set[types.Point]]{}
	head := res

	var wg sync.WaitGroup
	chFound := make(chan utilities.Set[types.Point])

	findPiecePlacements := func(piece Piece) {
		defer wg.Done()
		attemptedPermutations := utilities.NewSet([]Piece{}, 8)

		for j := 0; j < 2; j++ {
			piece = piece.Reflect(types.X)

			for i := 0; i < 4; i++ {
				piece = piece.Rotate90()
				if attemptedPermutations.Has(piece) {
					continue // Already tried this one
				}
				attemptedPermutations.Add(piece)

				// Check all possible placements
				// TODO: Optimize?
				piecePoints := piece.ToPoints()
				for origin := range piecePoints {
					absPoints := utilities.NewSet([]types.Point{}, len(piecePoints))
					for pp := range piecePoints {
						absPoints.Add(pp.Translate(pt.X-origin.X, pt.Y-origin.Y))
					}
					if b.validPlacement(absPoints, owner) {
						chFound <- absPoints
					}
				}
			}
		}
	}

	wg.Add(pieces.Size())
	for piece := range pieces {
		go findPiecePlacements(piece)
	}

	var resultGroup sync.WaitGroup
	resultGroup.Add(1)
	go func() {
		defer resultGroup.Done()
		for found := range chFound {
			head.Next = &utilities.Node[utilities.Set[types.Point]]{Value: found}
			head = head.Next
		}
	}()

	wg.Wait()
	close(chFound)
	resultGroup.Wait()

	return res
}

func (b *Board) findTerritory(o types.Owner) []types.Point {
	territory := make([]types.Point, 0, b.maxX*b.maxY)
	for ii := 0; ii < int(b.maxX); ii++ {
		for jj := 0; jj < int(b.maxY); jj++ {
			if b.occupiedByPlayer(types.Point{X: ii, Y: jj}, o) {
				territory = append(territory, types.Point{X: ii, Y: jj})
			}
		}
	}
	return territory
}

func (b *Board) findCorners(territory []types.Point, owner types.Owner) []types.Point {
	corners := make([]types.Point, 0, b.maxX*b.maxY)
	for _, pt := range territory {
		corners = append(corners, b.getFreeCorners(pt, owner)...)
	}
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
