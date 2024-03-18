package utilities

import (
	"errors"
	"fmt"
	"math"
	"strings"
)

// TODO: dynamic board size
type Board struct {
	layout     [][]Owner
	circle     *Circle
	MaxX, MaxY uint
}

func NewBoard(players []Owner, pixelsPerPlayer uint, tighteningFactor float64) (*Board, error) {

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
		layout: make([][]Owner, diameter),
		MaxX:   diameter,
		MaxY:   diameter,
	}
	for i := range board.layout {
		board.layout[i] = make([]Owner, diameter)
		for j := range board.layout[i] {
			board.layout[i][j] = RESERVED
		}
	}

	circle := bresenhamCircle(radius, Point{int(radius), int(radius)})
	board.circle = circle

	// clear out the playable region
	for pt := range circle.pixels {
		rowSize := circle.center.X - pt.X
		var offset int
		if rowSize < 0 {
			rowSize = 1 - rowSize
			offset = circle.center.X
		} else {
			offset = pt.X

		}

		for i := 0; i < rowSize; i++ {
			board.layout[pt.Y][offset+i] = VACANT
		}
	}

	angleDelta := (2 * math.Pi) / float64(numPlayers)
	for ii, owner := range players {
		pt := circle.pointOnCircle(angleDelta * float64(ii))
		err := board.occupy(pt, owner)
		if err != nil {
			return nil, fmt.Errorf("unable to occupy origin pt %+v for owner %v", pt, owner)
		}
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
			s += " " + b.layout[jj][ii].ToString() + " "
		}
		s += "|"
	}
	s += "\n-" + strings.Repeat("---", int(b.MaxX)) + "-"
	return s
}

func (b *Board) inbounds(square Point) bool {
	return square.X >= 0 && square.X < int(b.MaxX) && square.Y >= 0 && square.Y < int(b.MaxY)
}

func (b *Board) occupy(square Point, owner Owner) error {
	if !b.vacant(square) {
		return errors.New("square is occupied")
	}
	b.layout[square.X][square.Y] = owner
	return nil
}

func (b *Board) vacate(square Point) {
	b.layout[square.X][square.Y] = VACANT
}

func (b *Board) owner(square Point) Owner {
	return b.layout[square.X][square.Y]
}

func (b *Board) vacant(square Point) bool {
	return b.owner(square) == VACANT
}

func (b *Board) occupiedByPlayer(square Point, owner Owner) bool {
	return b.owner(square) == owner
}

func (b *Board) validPlacement(origin Point, p Piece, owner Owner) bool {
	return true
}

func (b *Board) hasSelfSide(pt Point, owner Owner) bool {
	return true
}

func (b *Board) hasCorner(pt Point, owner Owner) bool {
	if pt.X == 0 && pt.Y == 0 && b.vacant(pt) {
		// 0 0 is always a valid corner, if it is unoccupied
		return true
	}
	ul := pt.GetAdjacent(UP).GetAdjacent(LEFT)
	ur := pt.GetAdjacent(UP).GetAdjacent(RIGHT)
	dl := pt.GetAdjacent(DOWN).GetAdjacent(LEFT)
	dr := pt.GetAdjacent(DOWN).GetAdjacent(RIGHT)

	return ((b.inbounds(ul) && b.occupiedByPlayer(ul, owner)) ||
		(b.inbounds(ur) && b.occupiedByPlayer(ur, owner)) ||
		(b.inbounds(dl) && b.occupiedByPlayer(dl, owner)) ||
		(b.inbounds(dr) && b.occupiedByPlayer(dr, owner)))
}

func (b *Board) Place(origin Point, p Piece, owner Owner) (bool, error) {
	valid := b.validPlacement(origin, p, owner)
	if !valid {
		return false, errors.New("invalid placement")
	}
	pieceCoords := p.ToPoints()
	for pt := range pieceCoords {
		abs_pt := origin.Translate(int(pt.x), int(pt.y))
		err := b.occupy(abs_pt, owner)
		if err != nil {
			for clearPt := range pieceCoords {
				b.vacate(origin.Translate(int(clearPt.x), int(clearPt.y)))
			}
			return false, err
		}
	}
	return true, nil

}
