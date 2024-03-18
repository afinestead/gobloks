package utilities

import (
	"fmt"
)

func GeneratePieceSet(degree uint8) (PieceSet, uint, error) {
	if degree > MaxPieceDegree {
		return nil, 0, fmt.Errorf("degree must not exceed %v", MaxPieceDegree)
	}
	chResult := make(chan Piece)

	var ii uint8
	channels := make([]chan Piece, degree+1)
	for ii = 0; ii < degree; ii++ {
		channels[ii] = make(chan Piece)
	}
	channels[degree] = nil

	generator := func(chRecv chan Piece, chEscalate chan Piece) {
		uniquePieces := PieceSet{}
		for piece := range chRecv {
			for nextPiece := range generateNextPieces(piece) {
				if !uniquePieces.Has(nextPiece) { // nextPiece is unique
					uniquePieces.Add(nextPiece) // Add it to the set
					// escalate it to a higher order piece generator
					if chEscalate != nil {
						chEscalate <- nextPiece
					}
					// and send it to the results!
					chResult <- nextPiece
				}
			}
		}
		// All pieces have been dealt with. Close proper channels
		if chEscalate != nil {
			close(chEscalate)
		} else {
			close(chResult)
		}
	}

	for ii = 0; ii < degree; ii++ {
		go generator(channels[ii], channels[ii+1])
	}
	baseChannel := channels[0]
	if baseChannel != nil {
		baseChannel <- Piece{}
		close(baseChannel)
	} else {
		close(chResult)
	}

	var pix uint
	pieces := PieceSet{}
	for piece := range chResult {
		pieces.Add(piece)
		pix += uint(piece.Size())
	}
	return pieces, pix, nil
}

func generateNextPieces(piece Piece) PieceSet {
	pieceSet := PieceSet{}
	if piece.hash == 0 {
		pieceSet.Add(Piece{1, 1})
	} else {
		var shift uint8
		// shift up and over one row+column, if we can
		if getRow(piece.repr, MaxPieceDegree-1) == 0 {
			shift += 1
		}
		if getColumn(piece.repr, MaxPieceDegree-1) == 0 {
			shift += MaxPieceDegree
		}
		piece.repr <<= uint64(shift)

		for pt := range piece.ToPoints() {
			for _, dir := range []Direction{UP, DOWN, LEFT, RIGHT} {
				adj, err := pt.GetAdjacent(dir)
				if err != nil {
					continue
				}
				if !piece.hasPoint(adj) {
					pieceSet.Add(piece.addPoint(adj))
				}
			}
		}
	}
	return pieceSet
}
