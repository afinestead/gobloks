package utilities

func GeneratePieceSet(degree int) []Piece {
	chResult := make(chan Piece)

	channels := []chan Piece{}
	for ii := 0; ii < degree; ii++ {
		channels = append(channels, make(chan Piece))
	}
	channels = append(channels, nil)

	generator := func(chRecv chan Piece, chEscalate chan Piece) {
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
		baseChannel <- Piece{}
		close(baseChannel)
	} else {
		close(chResult)
	}

	pieceSet := []Piece{}
	for piece := range chResult {
		pieceSet = append(pieceSet, piece)
	}
	return pieceSet
}

func generateNextPieces(piece Piece) PieceSet {
	if len(piece) == 0 {
		return []Piece{NewPiece([]Point{{X: 0, Y: 0}})} // return base piece
	}
	pieceSet := PieceSet{}
	for pt := range piece {
		for _, dir := range []Direction{UP, DOWN, LEFT, RIGHT} {
			newPiece := piece.Copy()
			newPiece.Add(pt.GetAdjacent(dir))
			// Did adding this point change the piece? Add it to the set!
			if newPiece.Size() != piece.Size() {
				pieceSet.Add(newPiece.NormalizeToOrigin())
			}
		}
	}
	return pieceSet
}
