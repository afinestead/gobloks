package utilities

import (
	"cmp"
	"errors"
	"slices"
)

type Piece struct {
	points Set[Point]
	hash   uint64
}

func (p *Piece) ToString() string {
	points := make([]Point, len(p.points))
	ii := 0
	for pt := range p.points {
		points[ii] = pt
		ii++
	}
	slices.SortFunc(points, func(a, b Point) int {
		yCmp := cmp.Compare(b.Y, a.Y)
		if yCmp == 0 {
			return cmp.Compare(a.X, b.X)
		}
		return yCmp
	})
	var s string
	xIdx := 0
	yIdx := -1
	for _, pt := range points {
		if pt.Y != yIdx {
			if yIdx != -1 {
				s += "\n"
			}
			yIdx = pt.Y
			xIdx = 0
		}
		for ; xIdx < pt.X; xIdx++ {
			s += " "
		}
		s += "#"
		xIdx++
	}
	return s
}

func NewPiece(points []Point) Piece {
	piece := Piece{
		NormalizeToOrigin(NewSet[Point](points)),
		0,
	}

	// Compute a hash for each rotation/reflection possible
	// from the "bitset" that this piece makes in an 8x8 matrix
	hashes := make([]uint64, 8)
	idx := 0
	for ii := 0; ii < 2; ii++ { // Reflect twice
		piece.points = Reflect(piece.points, X)
		for _, rot := range []int{0, 90, 180, 270} {
			rotated := Rotate(piece.points, rot)
			rotHash, _ := serialize(rotated)
			hashes[idx] = rotHash
			idx++
		}
	}
	slices.Sort(hashes)
	piece.hash = hashes[0] // Take the lowest hash
	return piece
}

func serialize(points Set[Point]) (uint64, error) {
	var hash uint64 = 0
	for pt := range points {
		if (pt.X < 0) || (pt.Y < 0) {
			return 0, errors.New("piece must be normalized to origin")
		}
		if pt.X >= int(MaxPieceDegree) || pt.Y >= int(MaxPieceDegree) {
			return 0, errors.New("piece exceeds max degree")
		}
		hash |= (1 << (pt.X + (pt.Y * int(MaxPieceDegree))))
	}
	return hash, nil
}

func (p1 *Piece) Is(p2 *Piece) bool {
	return p1.hash == p2.hash
}

func (p *Piece) Size() int {
	return len(p.points)
}

func (p Piece) Add(point Point) Piece {
	newPoints := []Point{point}
	for pt := range p.points {
		newPoints = append(newPoints, pt)
	}
	return NewPiece(newPoints)
}

func (p *Piece) Has(point Point) bool {
	for pt := range p.points {
		if pt == point {
			return true
		}
	}
	return false
}

func (p Piece) Copy() Piece {
	ptsCopty := []Point{}
	for pt := range p.points {
		ptsCopty = append(ptsCopty, pt)
	}
	return NewPiece(ptsCopty)
}

func (p Piece) Rotate(degrees int) Piece {
	return Piece{Rotate(p.points, degrees), p.hash}
}

func (p Piece) Reflect(ax Axis) Piece {
	return Piece{Reflect(p.points, ax), p.hash}
}

type PieceSet []*Piece

func (ps *PieceSet) Size() int {
	return len(*ps)
}

func (ps *PieceSet) Add(piece *Piece) bool {
	if !ps.Has(piece) {
		*ps = append(*ps, piece)
		return true
	}
	return false
}

func (ps *PieceSet) Has(piece *Piece) bool {
	for _, p := range *ps {
		if piece.Is(p) {
			return true
		}
	}
	return false
}
