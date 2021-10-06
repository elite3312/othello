package main

import (
	"othello/game"

	"github.com/gotk3/gotk3/gtk"
)

func main() {
	gtk.Init(nil)

	win, err := gtk.WindowNew(gtk.WINDOW_TOPLEVEL)
	if err != nil {
		panic(err)
	}
	win.SetTitle("othello")
	win.Connect("destroy", func() {
		gtk.MainQuit()
	})

	game.New(win, game.Parameter{}, 8)

	win.SetDefaultSize(500, 500)
	win.SetResizable(false)
	win.SetPosition(gtk.WIN_POS_CENTER_ALWAYS)

	gtk.Main()
}
