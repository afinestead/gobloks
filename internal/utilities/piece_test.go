package utilities

import (
	"testing"
)

func TestPieceCreation(t *testing.T) {
	var p Piece
	var expSize int
	var expHash uint64

	checkSize := func(exp, res int) {
		if exp != res {
			t.Errorf("Unexpected piece size, expected %v, got %v", exp, res)
		}
	}
	checkHash := func(exp, res uint64) {
		if exp != res {
			t.Errorf("Unexpected piece hash, expected %v, got %v", exp, res)
		}
	}

	p = NewPiece([]Point{
		{X: 0, Y: 0},
	})
	expSize = 1
	expHash = 0b1
	checkSize(expSize, p.Size())
	checkHash(expHash, p.hash)

	p = NewPiece([]Point{
		{X: 0, Y: 0},
		{X: 0, Y: 1},
	})
	expSize = 2
	expHash = 0b11
	checkSize(expSize, p.Size())
	checkHash(expHash, p.hash)

	p = NewPiece([]Point{
		{X: 0, Y: 1},
		{X: 0, Y: 1},
	})
	expSize = 1
	expHash = 0b1
	checkSize(expSize, p.Size())
	checkHash(expHash, p.hash)

	p = NewPiece([]Point{
		{X: 0, Y: 0},
		{X: 0, Y: 1},
		{X: 0, Y: 2},
	})
	expSize = 3
	expHash = 0b111
	checkSize(expSize, p.Size())
	checkHash(expHash, p.hash)
}

func TestPieceEquality(t *testing.T) {
	p1 := NewPiece([]Point{
		{X: 0, Y: 0},
	})
	p2 := NewPiece([]Point{
		{X: 0, Y: 0},
	})

	expected := true
	result := p1.Is(&p2)
	if result != expected {
		t.Errorf("Piece equality failed! %+v != %+v", p1, p2)
	}

	p2 = NewPiece([]Point{
		{X: 0, Y: 1},
	})
	expected = true
	result = p1.Is(&p2)
	if result != expected {
		t.Errorf("Piece equality failed! %+v != %+v", p1, p2)
	}

	p1 = NewPiece([]Point{
		{X: 0, Y: 0},
		{X: 0, Y: 1},
	})
	p2 = NewPiece([]Point{
		{X: 0, Y: 0},
		{X: 0, Y: 1},
	})
	expected = true
	result = p1.Is(&p2)
	if result != expected {
		t.Errorf("Piece equality failed! %+v != %+v", p1, p2)
	}

	p1 = NewPiece([]Point{
		{X: 0, Y: 0},
		{X: 0, Y: -1},
	})
	p2 = NewPiece([]Point{
		{X: 0, Y: 0},
		{X: 0, Y: 1},
	})
	expected = true
	result = p1.Is(&p2)
	if result != expected {
		t.Errorf("Piece equality failed! %+v != %+v", p1, p2)
	}

	p1 = NewPiece([]Point{
		{X: 0, Y: 0}, //
		{X: 0, Y: 1}, //  # #
		{X: 1, Y: 1}, //  #
	})
	p2 = NewPiece([]Point{
		{X: 1, Y: 1}, //  # #
		{X: 1, Y: 2}, //    #
		{X: 0, Y: 2}, //
	})
	expected = true
	result = p1.Is(&p2)
	if result != expected {
		t.Errorf("Piece equality failed! %+v != %+v", p1, p2)
	}

	p1 = NewPiece([]Point{
		{X: 0, Y: 0}, //
		{X: 0, Y: 1}, //  # #
		{X: 1, Y: 1}, //  #
	})
	p2 = NewPiece([]Point{
		{X: 1, Y: 1}, //  #
		{X: 1, Y: 2}, //  #
		{X: 1, Y: 0}, //  #
	})
	expected = false
	result = p1.Is(&p2)
	if result != expected {
		t.Errorf("Piece equality failed! %+v != %+v", p1, p2)
	}
}

func TestPieceAdd(t *testing.T) {
	p1 := NewPiece([]Point{
		{X: 1, Y: 1}, //  #
		{X: 1, Y: 2}, //  #
		{X: 1, Y: 0}, //  #
	})
	p1 = p1.Add(Point{X: 1, Y: 0})
	expected := 4
	if p1.Size() != expected {
		t.Errorf("Piece add failed! expected size %v, got %v", expected, p1.Size())
	}
}

// func TestPiecePrint(t *testing.T) {
// 	p1 := NewPiece([]Point{
// 		{X: 1, Y: 1}, //  #
// 		{X: 2, Y: 2}, //  #
// 		{X: 1, Y: 2}, //  #
// 		{X: 1, Y: 0}, //  #
// 	})
// 	fmt.Println(p1.ToString())

// 	t.Errorf("error")
// }

func TestPieceCopy(t *testing.T) {
	p1 := NewPiece([]Point{
		{X: 1, Y: 1}, //  #
		{X: 1, Y: 2}, //  #
		{X: 1, Y: 0}, //  #
	})

	p2 := p1

	if p1.Size() != p2.Size() {
		t.Errorf("Piece copy failed! %v != %v", p1.Size(), p2.Size())
	}

	p1 = p1.Add(Point{X: 1, Y: 0})

	if p1.Size() == p2.Size() {
		t.Errorf("Piece copy failed! %v == %v", p1.Size(), p2.Size())
	}
}

func TestPieceSetAdd(t *testing.T) {
	ps := PieceSet{}
	if ps.Size() != 0 {
		t.Errorf("PieceSet size failed! %v != %v", ps.Size(), 0)
	}

	p1 := NewPiece([]Point{{X: 0, Y: 0}})
	ps.Add(&p1)
	if ps.Size() != 1 {
		t.Errorf("PieceSet size failed! %v != %v", ps.Size(), 1)
	}

	// Try adding same piece again...
	ps.Add(&p1)
	if ps.Size() != 1 {
		t.Errorf("PieceSet size failed! %v != %v", ps.Size(), 1)
	}

	// and as a copy...
	p1Cpy := p1.Copy()
	ps.Add(&p1Cpy)
	if ps.Size() != 1 {
		t.Errorf("PieceSet size failed! %v != %v", ps.Size(), 1)
	}

	// And now a new piece
	p2 := NewPiece([]Point{{X: 0, Y: 0}, {X: 0, Y: 1}})
	ps.Add(&p2)
	if ps.Size() != 2 {
		t.Errorf("PieceSet size failed! %v != %v", ps.Size(), 2)
	}
}
