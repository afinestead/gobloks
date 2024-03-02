package utilities

import (
	"fmt"
)

func GeneratePieceSet(degree int) (PieceSet, error) {
	if degree > int(MaxPieceDegree) {
		return nil, fmt.Errorf("degree must not exceed %v", MaxPieceDegree)
	}
	chResult := make(chan *Piece)

	// channels := []chan *Piece{}
	channels := make([]chan *Piece, degree+1)
	for ii := 0; ii < degree; ii++ {
		channels[ii] = make(chan *Piece)
	}
	channels[degree] = nil

	generator := func(chRecv chan *Piece, chEscalate chan *Piece) {
		uniquePieces := PieceSet{}
		for piece := range chRecv {
			for _, nextPiece := range generateNextPieces(piece) {
				if uniquePieces.Add(nextPiece) { // nextPiece is unique and was added to the set
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

	for ii := 0; ii < degree; ii++ {
		go generator(channels[ii], channels[ii+1])
	}
	baseChannel := channels[0]
	if baseChannel != nil {
		baseChannel <- &Piece{}
		close(baseChannel)
	} else {
		close(chResult)
	}

	pieces := []*Piece{}
	for piece := range chResult {
		pieces = append(pieces, piece)
	}
	return pieces, nil
}

func generateNextPieces(piece *Piece) PieceSet {
	if piece.Size() == 0 {
		newPiece := NewPiece([]Point{{0, 0}})
		return []*Piece{&newPiece}
	}
	pieceSet := PieceSet{}
	for pt := range piece.points {
		for _, dir := range []Direction{UP, DOWN, LEFT, RIGHT} {
			newPiece := piece.Add(pt.GetAdjacent(dir))
			// Did adding this point change the piece? Add it to the set!
			if newPiece.Size() != piece.Size() {
				pieceSet.Add(&newPiece)
			}
		}
	}
	return pieceSet
}
