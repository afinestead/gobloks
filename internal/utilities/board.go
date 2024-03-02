package utilities

// import "errors"

// type Board struct {
// 	layout [BoardX][BoardY]Owner
// }

// func (b *Board) inbounds(square Point) bool {
// 	return square.X >= 0 && square.X < BoardX && square.Y >= 0 && square.Y < BoardY
// }

// func (b *Board) occupy(square Point, owner Owner) (bool, error) {
// 	if !b.vacant(square) {
// 		return false, errors.New("square is occupied")
// 	}
// 	b.layout[square.X][square.Y] = owner
// 	return true, nil
// }

// func (b *Board) vacate(square Point) {
// 	b.layout[square.X][square.Y] = UNOWNED
// }

// func (b *Board) owner(square Point) Owner {
// 	return b.layout[square.X][square.Y]
// }

// func (b *Board) vacant(square Point) bool {
// 	return b.owner(square) == UNOWNED
// }

// func (b *Board) occupiedByPlayer(square Point, owner Owner) bool {
// 	return b.owner(square) == owner
// }

// func (b *Board) validPlacement(origin Point, p Piece, owner Owner) bool {
// 	return true
// }

// // func (b *Board) hasSelfSide(pt Point, owner Owner) bool {

// // }

// func (b *Board) hasCorner(pt Point, owner Owner) bool {
// 	if pt.X == 0 && pt.Y == 0 && b.vacant(pt) {
// 		// 0 0 is always a valid corner, if it is unoccupied
// 		return true
// 	}
// 	ul := pt.GetAdjacent(UP).GetAdjacent(LEFT)
// 	ur := pt.GetAdjacent(UP).GetAdjacent(RIGHT)
// 	dl := pt.GetAdjacent(DOWN).GetAdjacent(LEFT)
// 	dr := pt.GetAdjacent(DOWN).GetAdjacent(RIGHT)

// 	return ((b.inbounds(ul) && b.occupiedByPlayer(ul, owner)) ||
// 		(b.inbounds(ur) && b.occupiedByPlayer(ur, owner)) ||
// 		(b.inbounds(dl) && b.occupiedByPlayer(dl, owner)) ||
// 		(b.inbounds(dr) && b.occupiedByPlayer(dr, owner)))
// }

// func (b *Board) Place(origin Point, p Piece, owner Owner) (bool, error) {
// 	valid := b.validPlacement(origin, p, owner)
// 	if !valid {
// 		return false, errors.New("invalid placement")
// 	}
// 	for pt := range p.points {
// 		abs_pt := pt.Translate(pt.X, pt.Y)
// 		ok, err := b.occupy(abs_pt, owner)
// 		if !ok {
// 			for clearPt := range p.points {
// 				b.vacate(clearPt.Translate(clearPt.X, clearPt.Y))
// 			}
// 			return false, err
// 		}
// 	}
// 	return true, nil

// }
