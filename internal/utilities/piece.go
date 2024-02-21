package utilities

import (
	"errors"
	"math"
)

type Piece Set[Point]

func NewPiece(points []Point) Piece {
	return Piece(NewSet[Point](points))
}

func (p1 Piece) Is(p2 Piece) bool {
	if p1.Size() != p2.Size() {
		return false
	}
	compareTo := p1.NormalizeToOrigin()

	for _, rot := range []int{0, 90, 180, 270} {
		rotated := p2.Rotate(rot)
		normalizedRot := rotated.NormalizeToOrigin()
		if Set[Point].Is(Set[Point](normalizedRot), Set[Point](compareTo)) {
			return true
		}
		reflected := rotated.Reflect(X)
		normalizedRef := reflected.NormalizeToOrigin()
		if Set[Point].Is(Set[Point](normalizedRef), Set[Point](compareTo)) {
			return true
		}
	}
	return false
}

func (p Piece) Size() int {
	return len(p)
}

func (p Piece) Add(point Point) {
	Set[Point].Add(Set[Point](p), point)
}

func (p Piece) Has(point Point) bool {
	return Set[Point].Has(Set[Point](p), point)
}

func (p Piece) Copy() Piece {
	pCopy := make(Piece)
	for pt := range p {
		pCopy.Add(pt)
	}
	return pCopy
}

func (p Piece) Rotate(degrees int) Piece {
	shape := []Point{}
	for pt := range p {
		shape = append(shape, pt.Rotate(degrees))
	}
	return NewPiece(shape)
}

func (p Piece) Reflect(ax Axis) Piece {
	shape := []Point{}
	for pt := range p {
		shape = append(shape, pt.Reflect(ax))
	}
	return NewPiece(shape)
}

func (p Piece) translate(x int, y int) Piece {
	shape := []Point{}
	for pt := range p {
		shape = append(shape, pt.Translate(x, y))
	}
	return NewPiece(shape)
}

func (p Piece) minAxisCoordinate(ax Axis) (int, error) {
	if p.Size() == 0 {
		return 0, errors.New("cannot compute min on 0 size piece")
	}
	min := math.MaxInt
	for pt := range p {
		if ax == X && pt.X < min {
			min = pt.X
		} else if ax == Y && pt.Y < min {
			min = pt.Y
		}
	}
	return min, nil
}

func (p Piece) NormalizeToOrigin() Piece {
	minX, _ := p.minAxisCoordinate(X)
	minY, _ := p.minAxisCoordinate(Y)
	return p.translate(-minX, -minY)
}

type PieceSet []Piece

func (ps *PieceSet) Size() int {
	return len(*ps)
}

func (ps *PieceSet) Add(piece Piece) bool {
	if !ps.Has(piece) {
		*ps = append(*ps, piece)
		return true
	}
	return false
}

func (ps *PieceSet) Has(piece Piece) bool {
	for _, p := range *ps {
		if piece.Is(p) {
			return true
		}
	}
	return false
}
