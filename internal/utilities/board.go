package utilities

// import (
// 	"errors"
// )

// type Board struct {
// 	layout [][]Owner
// }

// func (b *Board) occupy(square Point, owner Owner) (bool, error) {
// 	if b.owner(square) != nil {
// 		return false, errors.New("square is occupied")
// 	}
// 	b.layout[square.X][square.Y] = owner
// 	return true, nil
// }

// func (b *Board) owner(square Point) *Owner {
// 	return &b.layout[square.X][square.Y]
// }

// func (b *Board) occupiedByPlayer(square Point, owner Owner) bool {
// 	return *b.owner(square) == owner
// }

// func (b *Board) validPlacement(origin Point, p Piece, owner Owner) bool {
// 	return true
// }

// func (b *Board) Place(origin Point, p Piece, owner Owner) error {
// 	valid := b.validPlacement(origin, p, owner)
// 	if !valid {
// 		return errors.New("invalid placement")
// 	}
// 	for _, pt := range p {
// 		abs_pt := pt.Translate(pt.X, pt.Y)
// 		b.occupy(abs_pt, owner)
// 	}
// 	return nil

// }
