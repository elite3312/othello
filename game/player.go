package game

import "othello/board"

type player interface {
	Move(board.Board) board.Point
}
