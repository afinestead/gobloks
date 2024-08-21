package game

import (
	"errors"
	"gobloks/internal/types"
	"gobloks/internal/utilities"
	"math"
	"math/bits"
)

type Piece struct{ repr, hash uint64 }
type PieceSet utilities.Set[Piece]
type PieceCoord struct{ X, Y uint8 }

func (p PieceCoord) getAdjacent(dir types.Direction) (PieceCoord, error) {
	var pt PieceCoord
	var err error
	if dir == types.UP {
		pt = PieceCoord{p.X, p.Y + 1}
	} else if dir == types.DOWN {
		pt = PieceCoord{p.X, p.Y - 1}
	} else if dir == types.LEFT {
		pt = PieceCoord{p.X - 1, p.Y}
	} else { // dir == RIGHT
		pt = PieceCoord{p.X + 1, p.Y}
	}
	if pt.X >= types.MaxPieceDegree || pt.Y >= types.MaxPieceDegree {
		// under/overflow
		err = errors.New("no adjacent point")
	}
	return pt, err
}

func NewPiece(piece uint64) Piece {
	p := Piece{piece, piece}

	for ii := 0; ii < 2; ii++ { // Reflect twice
		p.Reflect(types.X)

		for jj := 0; jj < 4; jj++ { // Rotate 4x
			p.Rotate90()
			if p.repr < p.hash { // take only the lowest hash
				p.hash = p.repr
			}
		}
	}
	return p
}

func (p *Piece) ToString() string {
	return stringify64(p.repr, '1', '0')
}

func (p1 Piece) Is(p2 Piece) bool {
	return p1.hash == p2.hash
}

func (p Piece) Size() uint8 {
	return uint8(bits.OnesCount64(p.hash))
}

func (p *Piece) Rotate90() {
	p.repr = rotate64(p.repr)
}

func (p *Piece) Reflect(ax types.Axis) {
	var ii uint8
	var ref uint64

	for ii = 0; ii < types.MaxPieceDegree; ii++ {
		if ax == types.Y {
			ref = reflectY64(p.repr)
		} else {
			ref = reflectX64(p.repr)
		}
	}
	p.repr = ref
}

func (p *Piece) Normalize() {
	p.repr = normalize64(p.repr)
}

func (p Piece) addPoint(pt PieceCoord) Piece {
	return NewPiece(p.repr | pointTo64(pt))
}

func (p Piece) hasPoint(pt PieceCoord) bool {
	return (p.repr & pointTo64(pt)) > 0
}

func (p Piece) ToPoints() utilities.Set[PieceCoord] {
	pointSet := utilities.Set[PieceCoord]{}
	var bitMask uint64
	var countedBits, ii uint8
	for countedBits < p.Size() {
		bitMask = (1 << ii)
		if (p.repr & bitMask) > 0 {
			countedBits++
			pointSet.Add(pointFrom64(bitMask))
		}
		ii++
	}
	return pointSet
}

func PieceFromPoints(points utilities.Set[PieceCoord]) Piece {
	var v uint64
	for pt := range points {
		v |= pointTo64(pt)
	}
	return NewPiece(v)
}

func ValidPieceCoords(points utilities.Set[PieceCoord]) bool {
	if points.Size() == 1 {
		return true
	}
	for pt := range points {
		// technically don't need to check the last point, but it's quick
		hasAdj := false
		for _, dir := range []types.Direction{types.UP, types.DOWN, types.LEFT, types.RIGHT} {
			adj, err := pt.getAdjacent(dir)
			if err == nil && points.Has(adj) {
				hasAdj = true
				break
			}
		}
		if !hasAdj {
			return false
		}
	}
	return true
}

func lsXY(n uint64) (uint8, uint8) {
	var ii, row, tz uint8
	lsx := types.MaxPieceDegree
	lsy := types.MaxPieceDegree * types.MaxPieceDegree
	for ii = types.MaxPieceDegree; ii > 0; ii-- {
		row = getRow(n, ii-1)
		tz = uint8(bits.TrailingZeros8(row))
		if tz < lsx {
			lsx = tz
		}
		if row > 0 {
			lsy = (ii - 1) * types.MaxPieceDegree
		}
	}
	return lsx, lsy
}

func normalize64(n uint64) uint64 {
	lsx, lsy := lsXY(n)
	return n >> (lsx + lsy)
}

func rotate64(n uint64) uint64 {
	var res uint64
	var lsx, lsy, ii uint8
	lsx = types.MaxPieceDegree
	lsy = types.MaxPieceDegree * types.MaxPieceDegree
	for ii = 0; ii < types.MaxPieceDegree; ii++ {
		newRow := getColumn(n, ii)
		rowShift := types.MaxPieceDegree * (types.MaxPieceDegree - ii - 1)
		tz := uint8(bits.TrailingZeros8(newRow))
		if newRow > 0 {
			if rowShift < lsy {
				lsy = rowShift
			}
			if tz < lsx {
				lsx = tz
			}
		}
		res |= uint64(newRow) << rowShift
	}
	return res >> (lsx + lsy)
}

func getColumn(n uint64, col uint8) uint8 {
	return uint8((((n << (types.MaxPieceDegree - 1 - col)) & types.ColumnMask) * types.MatrixMagic) >> ((types.MaxPieceDegree * types.MaxPieceDegree) - types.MaxPieceDegree) & uint64(math.MaxUint8))
}

func getRow(n uint64, row uint8) uint8 {
	return uint8(n >> (row * types.MaxPieceDegree))
}

func reflectY64(n uint64) uint64 {
	return bits.ReverseBytes64(n)
}

func reflectX64(n uint64) uint64 {
	return bits.ReverseBytes64(bits.Reverse64(n))
}

func stringify64(num uint64, filled, unfilled rune) string {
	var s string
	var ii, jj uint8
	for ii = 0; ii < types.MaxPieceDegree; ii++ {
		row := getRow(num, ii)
		s += "\n"
		for jj = 0; jj < types.MaxPieceDegree; jj++ {
			col := (row >> jj) & 1
			if col > 0 {
				s += string(filled)
			} else {
				s += string(unfilled)
			}
		}
	}
	return s
}

func pointTo64(pt PieceCoord) uint64 {
	return (1 << pt.X) << (pt.Y * types.MaxPieceDegree)
}

func pointFrom64(n uint64) PieceCoord {
	tz := uint8(bits.TrailingZeros64(n))
	return PieceCoord{tz % types.MaxPieceDegree, tz / types.MaxPieceDegree}
}

func (ps *PieceSet) Size() int {
	return len(*ps)
}

func (ps *PieceSet) Add(piece Piece) {
	piece.repr = piece.hash
	utilities.Set[Piece](*ps).Add(piece)
}

func (ps *PieceSet) Has(piece Piece) bool {
	piece.repr = piece.hash
	return utilities.Set[Piece](*ps).Has(piece)
}

func (ps *PieceSet) Remove(piece Piece) {
	piece.repr = piece.hash
	utilities.Set[Piece](*ps).Remove(piece)
}

func (ps *PieceSet) Copy() PieceSet {
	cpy := PieceSet{}
	for piece := range *ps {
		cpy.Add(piece)
	}
	return cpy
}
