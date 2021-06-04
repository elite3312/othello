package builtinai

import (
	"othello/board"
	"testing"
)

func TestPartialHeuristicChange(t *testing.T) {
	ai := New(board.WHITE, 6, 0)

	bd := board.NewBoard(6)
	bd.AssignBoard("+++++++++++++XXX++++OXX+++O+++++++++")

	p := board.NewPoint(5, 3)

	currentV := ai.heuristic(bd)
	c := bd.Copy()
	if !c.Put(board.WHITE, p) {
		t.Fatal("cannot put")
	}
	newV := ai.heuristic(c)

	aiV := ai.heuristicAfterPut(bd, currentV, p, board.WHITE)

	if newV != aiV {
		t.Error("error, orig:", currentV, "real:", newV, "but:", aiV)
		t.Error(bd.Visualize())
		t.Error(c.Visualize())
	}

	ai = New(board.WHITE, 6, 0)
	bd = board.NewBoard(6)
	bd.AssignBoard("++++++++++++XXOOO++XXOO+O+XXO++XXXO+")

	p = board.NewPoint(1, 4)

	currentV = ai.heuristic(bd)
	c = bd.Copy()
	if !c.Put(board.WHITE, p) {
		t.Error(c.Visualize())
		c.Assign(board.WHITE, p.X, p.Y)
		t.Error(c.Visualize())
		t.Fatal("cannot put")
	}
	newV = ai.heuristic(c)

	aiV = ai.heuristicAfterPut(bd, currentV, p, board.WHITE)

	if newV != aiV {
		t.Error("error, orig:", currentV, "real:", newV, "but:", aiV)
		t.Error(bd.Visualize())
		t.Error(c.Visualize())
	}

	ai = New(board.WHITE, 6, 0)
	bd = board.NewBoard(6)
	bd.AssignBoard("++++++++O+X+XXOOO++XXXXXO+XXO+OOOOO+")

	p = board.NewPoint(1, 4)

	currentV = ai.heuristic(bd)
	c = bd.Copy()
	if !c.Put(board.WHITE, p) {
		t.Error(c.Visualize())
		c.Assign(board.WHITE, p.X, p.Y)
		t.Error(c.Visualize())
		t.Fatal("cannot put")
	}
	newV = ai.heuristic(c)

	aiV = ai.heuristicAfterPut(bd, currentV, p, board.WHITE)

	if newV != aiV {
		t.Error("error, orig:", currentV, "real:", newV, "but:", aiV)
		t.Error(bd.Visualize())
		t.Error(c.Visualize())
	}
}
