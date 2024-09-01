package game

import (
	"errors"
	"fmt"
	"gobloks/internal/types"
	"gobloks/internal/utilities"
	"math"
	"strings"
)

type Board struct {
	layout     [][]types.Owner
	origins    map[types.PlayerID]types.Point
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
		layout:  make([][]types.Owner, diameter),
		origins: make(map[types.PlayerID]types.Point, numPlayers),
		MaxX:    diameter,
		MaxY:    diameter,
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
		if !b.vacant(absPt) {
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

func (b *Board) HasPlacement(owner types.Owner, pieces PieceSet) bool {
	plc := b.getPlacements(owner, pieces, true)

	for p := range plc {
		fmt.Println(p.origin, p.piece.ToString())
		break
	}

	// fmt.Println("placements", plc.Size())
	return plc.Size() > 0
}

func (b *Board) GetPlacements(owner types.Owner, pieces PieceSet) utilities.Set[PlacementInternal] {
	return b.getPlacements(owner, pieces, false)
}
