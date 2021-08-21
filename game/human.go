package game

import (
	"othello/board"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
)

type human struct {
	color board.Color
}

func newHuman(cl board.Color) player {
	return &human{color: cl}
}

func (h *human) Move(bd board.Board) board.Point {
	for ; ; time.Sleep(time.Microsecond * 50) {
		x, y := ebiten.CursorPosition()

		x = int(float64(x-MARGIN_X)/SPACE + FIX)
		y = int(float64(y-MARGIN_Y)/SPACE + FIX)

		if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
			p := board.NewPoint(x, y)
			if bd.PutPoint(h.color, p) {
				return p
			}
		}
	}
}
