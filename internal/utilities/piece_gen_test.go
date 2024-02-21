package utilities

import (
	"testing"
)

func TestGeneratingNextPieces(t *testing.T) {
	p := NewPiece([]Point{
		{X: 0, Y: 0},
	})
	generated := generateNextPieces(p)
	if generated.Size() != 1 {
		t.Errorf("Next piece generation failed.")
	}
	for _, piece := range generated {
		if !piece.Is(NewPiece([]Point{{X: 0, Y: 0}, {X: 0, Y: 1}})) {
			t.Errorf("Next piece generation failed.")
		}
	}
}

func TestPieceGenerator(t *testing.T) {

	expectedSizes := []int{0, 1, 2, 4, 9, 21}

	for ii, degree := range []int{0, 1, 2, 3, 4, 5} {
		result := GeneratePieceSet(degree)
		resultSize := len(result)
		expectedSize := expectedSizes[ii]
		if resultSize != expectedSize {
			t.Errorf("Unexpected piece set generated for degree %v. expected size %v, got size %v", degree, expectedSize, resultSize)
			t.Errorf("Result: %+v", result)
		}
	}
}
