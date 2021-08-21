package builtinai

import "othello/board"

type point struct {
	x, y int
}

func (p point) String() string {
	return string(rune('A'+p.y)) + string(rune('a'+p.x))
}

func (p point) toBoardPoint() board.Point {
	return board.NewPoint(p.x, p.y)
}
