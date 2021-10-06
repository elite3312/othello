package game

import (
	"fmt"
	"log"
	"othello/board"
	"othello/builtinai"
	"time"

	"github.com/gotk3/gotk3/gtk"
)

const (
	counterTextSize = 24
	timerTextSize   = 13
	nameTextSize    = 13
	maxNameLen      = 20
)

var (
	nullPoint = board.NewPoint(-1, -1)
	// unitSize  = fyne.NewSize(48, 48)
)

type game struct {
	window *gtk.Window
	bd     board.Board
	units  [][]*unit

	counterBlack Text
	counterWhite Text

	passBtn *gtk.ToggleButton
	com1    computer
	com2    computer
	now     board.Color

	blackSpent time.Duration
	whiteSpent time.Duration

	haveHuman bool
	over      bool
}

func New(win *gtk.Window, params Parameter, size int) {
	g := &game{}

	// units := make([][]*unit, size)
	// for i := range units {
	// 	units[i] = make([]*unit, size)
	// }
	// grid := container.New(layout.NewGridLayout(size))
	// for i := 0; i < size; i++ {
	// 	for j := 0; j < size; j++ {
	// 		u := newUnit(g, board.NONE, i, j)
	// 		grid.Add(u)
	// 		units[i][j] = u
	// 	}
	// }

	if params.BlackAgent == AgentBuiltIn {
		if size == 6 {
			g.com1 = builtinai.NewAI6(builtinai.BLACK, params.BlackAILevel)
		} else {
			g.com1 = builtinai.NewAI8(builtinai.BLACK, params.BlackAILevel)
		}
	} else if params.BlackAgent == AgentExternal {
		g.com1 = newCom(board.BLACK, params.BlackPath)
	}
	if params.WhiteAgent == AgentBuiltIn {
		if size == 6 {
			g.com2 = builtinai.NewAI6(builtinai.WHITE, params.WhiteAILevel)
		} else {
			g.com2 = builtinai.NewAI8(builtinai.WHITE, params.WhiteAILevel)
		}
	} else if params.WhiteAgent == AgentExternal {
		g.com2 = newCom(board.WHITE, params.WhitePath)
	}

	g.window = win
	// g.units = units
	g.now = params.GoesFirst
	g.bd = board.NewBoard(size)
	g.over = false
	g.haveHuman = g.com1 == nil || g.com2 == nil
	// g.counterBlack, g.counterWhite = newCounterText()

	// counterTile := container.NewGridWithColumns(2, g.counterBlack.CanvasText(), g.counterWhite.CanvasText())
	// nameText := newNameText(win.Canvas().Size(), params)

	var err error

	vBox, err := gtk.BoxNew(gtk.ORIENTATION_VERTICAL, 0)
	if err != nil {
		panic(err)
	}

	menuB, err := gtk.MenuBarNew()
	if err != nil {
		panic(err)
	}

	newGame, err := gtk.MenuItemNewWithLabel("選項")
	if err != nil {
		panic(err)
	}

	menuB.Append(newGame)
	vBox.PackStart(menuB, false, false, 0)

	res, err := gtk.MenuItemNewWithLabel("新遊戲")
	if err != nil {
		panic(err)
	}
	res.Connect("activate", func() {
		fmt.Println("clicked")
	})

	edit, err := gtk.MenuItemNewWithLabel("編輯棋盤")
	if err != nil {
		panic(err)
	}
	edit.Connect("activate", func() {
		fmt.Println("edit board")
	})

	quitItemBox, err := gtk.BoxNew(gtk.ORIENTATION_HORIZONTAL, 0)
	if err != nil {
		log.Fatal("Unable to create itemBox:", err)
	}

	quitItemPix, err := gtk.ImageNewFromIconName("application-exit", gtk.ICON_SIZE_MENU)
	if err != nil {
		log.Fatal("Unable to create itemPix:", err)
	}
	quitItemBox.PackStart(quitItemPix, false, false, 0)

	quitItemLabel, err := gtk.LabelNew("離開")
	if err != nil {
		log.Fatal("Unable to create itemLabel:", err)
	}
	quitItemBox.PackStart(quitItemLabel, false, false, 0)

	quitMenuItem, err := gtk.MenuItemNew()
	if err != nil {
		log.Fatal("Unable to create newMenuItem:", err)
	}
	quitMenuItem.Connect("activate", func() {
		gtk.MainQuit()
	})
	quitMenuItem.Add(quitItemBox)

	span, err := gtk.MenuNew()
	span.Append(res)
	span.Append(edit)
	span.Append(quitMenuItem)

	newGame.SetSubmenu(span)

	// menuB.Append(subMenu)

	g.passBtn, err = gtk.ToggleButtonNewWithLabel("pass")
	if err != nil {
		panic(err)
	}
	g.passBtn.Connect("clicked", func() {
		g.passBtn.SetSensitive(false)
		g.now = g.now.Opponent()
		// g.update(nullPoint)
	})
	g.passBtn.SetHExpand(true)

	restart, err := gtk.MenuButtonNew()
	if err != nil {
		panic(err)
	}
	restart.Connect("clicked", func() {
		dial, err := gtk.DialogNewWithButtons("confirm", win, gtk.DIALOG_DESTROY_WITH_PARENT, []interface{}{"yes", gtk.RESPONSE_ACCEPT}, []interface{}{"no", gtk.RESPONSE_CANCEL})
		if err != nil {
			panic(err)
		}
		dial.Activate()
	})
	restart.SetHExpand(true)

	editBtn, err := gtk.MenuButtonNew()
	if err != nil {
		panic(err)
	}
	editBtn.SetHExpand(true)

	// if g.com1 != nil || g.com2 != nil {
	// 	go g.round()
	// }
	// g.update(nullPoint)

	btnGrid, err := gtk.GridNew()
	if err != nil {
		panic(err)
	}
	// nb, err := gtk.NotebookNew()
	// if err != nil {
	// 	panic(err)
	// }
	// btnGrid.Attach(nb, 1, 1, 2, 2)
	// btnGrid.SetHAlign(gtk.ALIGN_CENTER)
	btnGrid.SetHExpand(true)
	btnGrid.Add(g.passBtn)

	// restart1, err := gtk.MenuButtonNew()
	// if err != nil {
	// panic(err)
	// }

	body, err := gtk.BoxNew(gtk.ORIENTATION_VERTICAL, 2)
	body.Add(vBox)
	body.Add(btnGrid)

	// menuB.Add(restart1)
	// menuB.Activate()

	// gr, _ := gtk.BoxNew(gtk.ORIENTATION_VERTICAL, 1)
	// gr.Add(menuB)

	// win.Add(gr)
	win.Add(body)

	win.ShowAll()
}

func (g *game) isBot(cl board.Color) bool {
	if cl == board.BLACK {
		return g.com1 != nil
	} else {
		return g.com2 != nil
	}
}

func (g *game) round() {
	var out string
	var err error
	defer g.cleanAndExit()
	for !g.over {
		if g.isBot(g.now) {
			start := time.Now()
			if g.now == board.BLACK {
				out, err = g.com1.Move(g.bd.String())
			} else {
				out, err = g.com2.Move(g.bd.String())
			}
			spent := time.Since(start)
			fmt.Println(g.now, "side spent:", spent)
			if g.now == board.BLACK {
				g.blackSpent += spent
			} else {
				g.whiteSpent += spent
			}
			if err != nil {
				g.aiError(err)
				break
			}
			g.bd.PutStr(g.now, out)
			g.now = g.now.Opponent()
			g.update(board.StrToPoint(out))
		} else {
			time.Sleep(time.Millisecond * 30)
		}
	}
}

func (g *game) update(current board.Point) {
	g.over = g.bd.IsOver()
	count := g.showValidAndCount(current)
	if count == 0 && !g.over {
		if g.haveHuman {
			// current side is human
			if (g.now == board.BLACK && g.com1 == nil) || (g.now == board.WHITE && g.com2 == nil) {
				// dialog.NewInformation("info", "you have to pass", g.window).Show()
				g.passBtn.SetVisible(true)
			} else { // current is computer
				// dialog.NewInformation("info", "computer have to pass\nit's your turn", g.window).Show()
				g.now = g.now.Opponent()
				g.update(nullPoint)
			}
		} else {
			g.now = g.now.Opponent()
			g.showValidAndCount(current)
		}
	}
	g.refreshCounter()
	if g.over {
		g.gameOver()
	}
}

func (g *game) refreshCounter() {
	blacks := g.bd.CountPieces(board.BLACK)
	whites := g.bd.CountPieces(board.WHITE)
	g.counterBlack.Update(fmt.Sprintf("black: %2d", blacks))
	g.counterWhite.Update(fmt.Sprintf("white: %2d", whites))
}

func (g *game) gameOver() {
	// var text string
	// winner := g.bd.Winner()
	// if winner == board.NONE {
	// 	text = "draw"
	// } else {
	// 	text = winner.String() + " won"
	// }
	// d := dialog.NewInformation("Game Over", text, g.window)
	// d.Resize(fyne.NewSize(250, 0))
	// d.Show()
	// fmt.Println("\ngame over")
	// fmt.Println("black total:", g.blackSpent, ", white total:", g.whiteSpent)
}

func (g *game) showValidAndCount(current board.Point) int {
	// count := 0
	// for i, line := range g.units {
	// 	for j, u := range line {
	// 		cl := g.bd.AtXY(i, j)
	// 		if g.bd.IsValidPoint(g.now, board.NewPoint(i, j)) {
	// 			u.SetResource(possible)
	// 			count++
	// 		} else {
	// 			u.setColor(cl)
	// 		}
	// 		if current.X == i && current.Y == j {
	// 			u.setColorCurrent(cl)
	// 		}
	// 	}
	// }
	// return count
	return 0
}

func (g *game) aiError(err error) {
	// if !g.over {
	// 	d := dialog.NewError(err, g.window)
	// 	d.SetOnClosed(func() { panic(err) })
	// 	d.Show()
	// }
}

func (g *game) cleanAndExit() {
	g.over = true
	if g.com1 != nil {
		g.com1.Close()
	}
	if g.com2 != nil {
		g.com2.Close()
	}
}

type unit struct {
	// g *game
	// widget.Icon
	// x, y  int
	// color board.Color
}

// func newUnit(g *game, cl board.Color, x, y int) *unit {
// 	u := &unit{g: g, color: cl, x: x, y: y}
// 	u.setColor(cl)
// 	u.ExtendBaseWidget(u)
// 	return u
// }

// func (u *unit) Tapped(ev *fyne.PointEvent) {
// 	if u.g.isBot(u.g.now) {
// 		return
// 	}
// 	p := board.NewPoint(u.x, u.y)
// 	if !u.g.bd.PutPoint(u.g.now, p) {
// 		return
// 	}

// 	u.g.now = u.g.now.Opponent()
// 	u.g.update(p)
// }

// func (u *unit) MinSize() fyne.Size {
// 	return unitSize
// }

// func (u *unit) setColor(cl board.Color) {
// 	if cl == board.BLACK {
// 		u.SetResource(blackImg)
// 	} else if cl == board.WHITE {
// 		u.SetResource(whiteImg)
// 	} else {
// 		u.SetResource(noneImg)
// 	}
// }

// func (u *unit) setColorCurrent(cl board.Color) {
// 	if cl == board.BLACK {
// 		u.SetResource(blackCurr)
// 	} else if cl == board.WHITE {
// 		u.SetResource(whiteCurr)
// 	}
// }
