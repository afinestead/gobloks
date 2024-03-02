package utilities

import (
	"fmt"
	"math"
	"math/bits"
	"slices"
)

type Piece2 struct {
	repr uint64
	hash uint64
}

func NewPiece2(piece uint64) Piece2 {
	p := Piece2{}
	hashes := make([]uint64, 8)
	p.repr = piece
	for ii := 0; ii < 2; ii++ { // Reflect twice
		p.Reflect(X)

		for jj := 0; jj < 4; jj++ { // Rotate 4x
			p.Rotate90()
			hashes[(ii*4)+jj] = p.repr
		}
	}
	slices.Sort(hashes)
	fmt.Println(hashes)
	p.hash = hashes[0] // Take the lowest hash
	return p
}

func (p1 Piece2) Is(p2 Piece2) bool {
	return p1.hash == p2.hash
}

func (p *Piece2) Rotate90() {
	p.repr = rotate64(p.repr)
}

func (p *Piece2) Reflect(ax Axis) {
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

func (p *Piece2) Normalize() {
	p.repr = normalize64(p.repr)
}

func (p *Piece2) ToString() string {
	return stringify64(p.repr)
}

func normalize64(n uint64) uint64 {
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
	var ii uint8
	var res uint64
	for ii = 0; ii < MaxPieceDegree; ii++ {
		res |= uint64(getRow(n, MaxPieceDegree-ii-1)) << (ii * MaxPieceDegree)
	}
	return res
}

func reflectX64(n uint64) uint64 {
	var ii uint8
	var res uint64
	for ii = 0; ii < MaxPieceDegree; ii++ {
		res |= uint64(flipRow(getRow(n, ii))) << (ii * MaxPieceDegree)
	}
	// fmt.Printf("Reflecting %b\nRes %b\n", n, res)
	return res
}

func flipRow(row uint8) uint8 {
	var ii, res uint8
	for ii = 0; ii < MaxPieceDegree; ii++ {
		res |= ((row & (1 << ii)) >> ii) << (MaxPieceDegree - ii - 1)
	}
	return res
}

func stringify64(num uint64) string {
	var s string
	var ii, jj uint8
	for ii = 0; ii < MaxPieceDegree; ii++ {
		row := getRow(num, ii)
		s += "\n"
		for jj = 0; jj < MaxPieceDegree; jj++ {
			col := (row >> jj) & 1
			if col > 0 {
				s += "#"
			} else {
				s += " "
			}
		}
	}
	return s
}

func fromPoints() {

}

func toPoints() {

}
