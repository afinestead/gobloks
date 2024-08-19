package game_test

import (
	"testing"
)

func TestGeneratingNextPieces(t *testing.T) {
	p := PieceFromPoints(NewSet[PieceCoord]([]PieceCoord{{0, 0}}))
	generated := generateNextPieces(p)
	expectedSize := 1
	resultSize := generated.Size()

	if expectedSize != resultSize {
		t.Errorf("Next piece generation failed. Expected set of size %v, got size %v\n", expectedSize, resultSize)
	}
	expected := PieceFromPoints(NewSet[PieceCoord]([]PieceCoord{{0, 0}, {0, 1}}))
	for piece := range generated {
		if !piece.Is(expected) {
			t.Errorf("Next piece generation failed. expected:\n%s\ngot:\n%s\n", expected.ToString(), piece.ToString())
		}
	}
}

func TestPieceGenerator(t *testing.T) {

	expectedSizes := []int{0, 1, 2, 4, 9, 21}

	for ii, degree := range []uint8{0, 1, 2, 3, 4, 5} {
		result, err := GeneratePieceSet(degree)
		if err != nil {
			t.Errorf("generator returned error: %s", err)
		}
		resultSize := len(result)
		expectedSize := expectedSizes[ii]
		if resultSize != expectedSize {

			t.Errorf("Unexpected piece set generated for degree %v. expected size %v, got size %v", degree, expectedSize, resultSize)
		}
	}
}
