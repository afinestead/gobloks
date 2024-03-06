package utilities

import (
	"errors"
	"math"
	"math/bits"
)

type Piece struct{ repr, hash uint64 }
type PieceSet Set[Piece]
type PieceCoord struct{ x, y uint8 }

func (p PieceCoord) GetAdjacent(dir Direction) (PieceCoord, error) {
	var pt PieceCoord
	var err error
	if dir == UP {
		pt = PieceCoord{p.x, p.y + 1}
	} else if dir == DOWN {
		pt = PieceCoord{p.x, p.y - 1}
	} else if dir == LEFT {
		pt = PieceCoord{p.x - 1, p.y}
	} else { // dir == RIGHT
		pt = PieceCoord{p.x + 1, p.y}
	}
	if pt.x >= MaxPieceDegree || pt.y >= MaxPieceDegree {
		// under/overflow
		err = errors.New("no adjacent point")
	}
	return pt, err
}

func NewPiece(piece uint64) Piece {
	p := Piece{piece, piece}

	for ii := 0; ii < 2; ii++ { // Reflect twice
		p.Reflect(X)

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

func (p *Piece) Reflect(ax Axis) {
	var ii uint8
	var ref uint64

	for ii = 0; ii < MaxPieceDegree; ii++ {
		if ax == Y {
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

func (p Piece) ToPoints() Set[PieceCoord] {
	pointSet := Set[PieceCoord]{}
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

func FromPoints(points Set[PieceCoord]) Piece {
	var v uint64
	for pt := range points {
		v |= pointTo64(pt)
	}
	return NewPiece(v)
}

func lsXY(n uint64) (uint8, uint8) {
	var ii, row, tz uint8
	lsx := MaxPieceDegree
	lsy := MaxPieceDegree * MaxPieceDegree
	for ii = MaxPieceDegree; ii > 0; ii-- {
		row = getRow(n, ii-1)
		tz = uint8(bits.TrailingZeros8(row))
		if tz < lsx {
			lsx = tz
		}
		if row > 0 {
			lsy = (ii - 1) * MaxPieceDegree
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
	lsx = MaxPieceDegree
	lsy = MaxPieceDegree * MaxPieceDegree
	for ii = 0; ii < MaxPieceDegree; ii++ {
		newRow := getColumn(n, ii)
		rowShift := MaxPieceDegree * (MaxPieceDegree - ii - 1)
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
	return uint8((((n << (MaxPieceDegree - 1 - col)) & ColumnMask) * MatrixMagic) >> ((MaxPieceDegree * MaxPieceDegree) - MaxPieceDegree) & uint64(math.MaxUint8))
}

func getRow(n uint64, row uint8) uint8 {
	return uint8(n >> (row * MaxPieceDegree))
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
	for ii = 0; ii < MaxPieceDegree; ii++ {
		row := getRow(num, ii)
		s += "\n"
		for jj = 0; jj < MaxPieceDegree; jj++ {
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
	return (1 << pt.x) << (pt.y * MaxPieceDegree)
}

func pointFrom64(n uint64) PieceCoord {
	tz := uint8(bits.TrailingZeros64(n))
	return PieceCoord{tz % MaxPieceDegree, tz / MaxPieceDegree}
}

func (ps *PieceSet) Size() int {
	return len(*ps)
}

func (ps *PieceSet) Add(piece Piece) {
	piece.repr = piece.hash
	Set[Piece](*ps).Add(piece)
}

func (ps *PieceSet) Has(piece Piece) bool {
	piece.repr = piece.hash
	return Set[Piece](*ps).Has(piece)
}
