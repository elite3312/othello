package game

import (
	"fmt"
	"othello/board"
	"othello/builtinai"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

const (
	COOLDOWN = 0
)

type game struct {
	turn      bool
	over      bool
	bd        board.Board
	lastClick time.Time
	player1   player
	player2   player
	lastMove  board.Point
	winner    board.Color
	available []board.Point
}

func NewGame() *game {
	bd := board.NewBoard(8)

	g := &game{
		turn:      true,
		over:      false,
		bd:        bd,
		lastMove:  board.NewPoint(-9, -9),
		winner:    board.NONE,
		available: bd.AllValidPoint(board.BLACK),
	}

	g.player1 = newHuman(board.BLACK)
	g.player2 = builtinai.NewAI8(builtinai.WHITE)

	return g
}

func (g *game) Update() error {

	if ebiten.IsKeyPressed(ebiten.KeyP) {
		if time.Since(g.lastClick) > COOLDOWN {
			fmt.Println(g.bd)
			g.lastClick = time.Now()
		}
	} else if ebiten.IsKeyPressed(ebiten.KeyV) {
		if time.Since(g.lastClick) > COOLDOWN {
			fmt.Println(g.bd.Visualize())
			g.lastClick = time.Now()
		}
	} else if ebiten.IsKeyPressed(ebiten.KeyR) {
		if time.Since(g.lastClick) > COOLDOWN {
			g.restart()
			g.lastClick = time.Now()
		}
	}

	return nil
}

func (g *game) Draw(screen *ebiten.Image) {
	g.drawBoard(screen)
	g.drawStones(screen)
	if g.over {
		g.drawEnd(screen)
	}

	ebitenutil.DebugPrintAt(screen, fmt.Sprintf("fps: %.02f", ebiten.CurrentFPS()), WIN_WIDTH-65, 0)
}

func (g *game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return outsideWidth, outsideHeight
}

func (g *game) Round() {
	for !g.over {
		bd := g.bd.Copy()
		if g.turn {
			p := g.player1.Move(bd)
			if ok := g.bd.PutPoint(board.BLACK, p); !ok {
				panic("cannot put")
			}
			g.available = g.bd.AllValidPoint(board.WHITE)
			g.lastMove = p
		} else {
			p := g.player2.Move(bd)
			if ok := g.bd.PutPoint(board.WHITE, p); !ok {
				panic("cannot put")
			}
			g.available = g.bd.AllValidPoint(board.BLACK)
			g.lastMove = p
		}
		g.over = g.bd.IsOver()
		if g.over {
			g.winner = g.bd.Winner()
		}
		if len(g.available) != 0 {
			g.turn = !g.turn
		}
	}
}

func (g *game) restart() {
	g.over = true
	*g = *NewGame()
	(*g).turn = !(*g).turn
	go g.Round()
}
