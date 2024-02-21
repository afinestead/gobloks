package utilities

import (
	"testing"
)

func TestPieceCreation(t *testing.T) {
	p := NewPiece([]Point{
		{X: 0, Y: 0},
	})
	if p.Size() != 1 {
		t.Errorf("Unexpected piece size, expected 1, got %v", p.Size())
	}
	p2 := NewPiece([]Point{
		{X: 0, Y: 0},
		{X: 0, Y: 1},
	})
	if p2.Size() != 2 {
		t.Errorf("Unexpected piece size, expected 2, got %v", p.Size())
	}
	p3 := NewPiece([]Point{
		{X: 0, Y: 1},
		{X: 0, Y: 1},
	})
	if p3.Size() != 1 {
		t.Errorf("Unexpected piece size, expected 1, got %v", p.Size())
	}
}

func TestPieceNormalization(t *testing.T) {
	p := NewPiece([]Point{
		{X: 0, Y: 0},
	})
	expected := NewPiece([]Point{
		{X: 0, Y: 0},
	})
	result := p.NormalizeToOrigin()

	if result.Size() != expected.Size() {
		t.Errorf("Piece normalization failed! %+v != %+v", result, expected)
	}
	for pt := range expected {
		if !result.Has(pt) {
			t.Errorf("Piece normalization failed! %+v != %+v", result, expected)
		}
	}

	p = NewPiece([]Point{
		{X: 0, Y: 1},
	})
	expected = NewPiece([]Point{
		{X: 0, Y: 0},
	})
	result = p.NormalizeToOrigin()

	if result.Size() != expected.Size() {
		t.Errorf("Piece normalization failed! %+v != %+v", result, expected)
	}
	for pt := range expected {
		if !result.Has(pt) {
			t.Errorf("Piece normalization failed! %+v != %+v", result, expected)
		}
	}

	p = NewPiece([]Point{
		{X: 0, Y: 0},
		{X: 0, Y: 1},
	})
	expected = NewPiece([]Point{
		{X: 0, Y: 0},
		{X: 0, Y: 1},
	})
	result = p.NormalizeToOrigin()
	if result.Size() != expected.Size() {
		t.Errorf("Piece normalization failed! %+v != %+v", result, expected)
	}
	for pt := range expected {
		if !result.Has(pt) {
			t.Errorf("Piece normalization failed! %+v != %+v", result, expected)
		}
	}

	p = NewPiece([]Point{
		{X: 0, Y: 0},
		{X: 0, Y: -1},
	})
	expected = NewPiece([]Point{
		{X: 0, Y: 1},
		{X: 0, Y: 0},
	})
	result = p.NormalizeToOrigin()
	for pt := range expected {
		if !result.Has(pt) {
			t.Errorf("Piece normalization failed! %+v != %+v", result, expected)
		}
	}

	p = NewPiece([]Point{
		{X: -1, Y: 0},
		{X: 0, Y: -1},
	})
	expected = NewPiece([]Point{
		{X: 0, Y: 1},
		{X: 1, Y: 0},
	})
	result = p.NormalizeToOrigin()
	for pt := range expected {
		if !result.Has(pt) {
			t.Errorf("Piece normalization failed! %+v != %+v", result, expected)
		}
	}

	p = NewPiece([]Point{
		{X: 1, Y: 1}, //  # #
		{X: 1, Y: 2}, //    #
		{X: 0, Y: 2}, //
	})
	expected = NewPiece([]Point{
		{X: 1, Y: 0}, //
		{X: 1, Y: 1}, //  # #
		{X: 0, Y: 1}, //    #
	})
	result = p.NormalizeToOrigin()
	for pt := range expected {
		if !result.Has(pt) {
			t.Errorf("Piece normalization failed! %+v != %+v", result, expected)
		}
	}
}

func TestPieceEquality(t *testing.T) {
	p1 := NewPiece([]Point{
		{X: 0, Y: 0},
	})
	p2 := NewPiece([]Point{
		{X: 0, Y: 0},
	})

	expected := true
	result := p1.Is(p2)
	if result != expected {
		t.Errorf("Piece equality failed! %+v != %+v", p1, p2)
	}

	p2 = NewPiece([]Point{
		{X: 0, Y: 1},
	})
	expected = true
	result = p1.Is(p2)
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
	result = p1.Is(p2)
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
	result = p1.Is(p2)
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
	result = p1.Is(p2)
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
	result = p1.Is(p2)
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
	p1.Add(Point{X: 0, Y: 0})
	expected := 4
	if p1.Size() != expected {
		t.Errorf("Piece add failed! expected size %v, got %v", expected, p1.Size())
	}
}

func TestPieceCopy(t *testing.T) {
	p1 := NewPiece([]Point{
		{X: 1, Y: 1}, //  #
		{X: 1, Y: 2}, //  #
		{X: 1, Y: 0}, //  #
	})

	p2 := p1.Copy()

	if p1.Size() != p2.Size() {
		t.Errorf("Piece copy failed! %v != %v", p1.Size(), p2.Size())
	}

	p1.Add(Point{X: 0, Y: 0})

	if p1.Size() == p2.Size() {
		t.Errorf("Piece copy failed!")
	}
}

func TestPieceSetAdd(t *testing.T) {
	ps := PieceSet{}
	if ps.Size() != 0 {
		t.Errorf("PieceSet size failed! %v != %v", ps.Size(), 0)
	}

	p1 := NewPiece([]Point{{X: 0, Y: 0}})
	ps.Add(p1)
	if ps.Size() != 1 {
		t.Errorf("PieceSet size failed! %v != %v", ps.Size(), 1)
	}

	// Try adding same piece again...
	ps.Add(p1)
	if ps.Size() != 1 {
		t.Errorf("PieceSet size failed! %v != %v", ps.Size(), 1)
	}

	// and as a copy...
	ps.Add(p1.Copy())
	if ps.Size() != 1 {
		t.Errorf("PieceSet size failed! %v != %v", ps.Size(), 1)
	}

	// And now a new piece
	p2 := NewPiece([]Point{{X: 0, Y: 0}, {X: 0, Y: 1}})
	ps.Add(p2)
	if ps.Size() != 2 {
		t.Errorf("PieceSet size failed! %v != %v", ps.Size(), 2)
	}
}
