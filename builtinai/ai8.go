package builtinai

import (
	"fmt"
	"othello/board"
)

const (
	PHASE1DEPTH8 = 10 // 8x8
	PHASE2DEPTH8 = 16 // 8x8

	SIZE8      = 8
	TOTALVALUE = 13752
)

type AI8 struct {
	color    color
	opponent color

	// table map[bboard8]int

	// phase 1 or phase 2
	phase int

	// currently limit depth
	depth int

	// maximum reached depth
	reachedDepth int

	// traversed nodes count
	nodes int

	nodesPool pool
}

func NewAI8(cl color) *AI8 {
	ai := AI8{
		color: cl,
		// table:    make(map[bboard8]int),
		opponent: cl.reverse(),
	}
	ai.nodesPool = newPool(32)

	return &ai
}

func (ai *AI8) Move(bd board.Board) board.Point {
	aibd := newBboard8(bd.String())
	ai.nodes = 0

	ai.setPhase(aibd)
	ai.setDepth()

	best := ai.alphaBetaHelper(aibd, ai.depth)
	ai.printValue(best)

	bestPoint := point{best.loc % SIZE8, best.loc / SIZE8}
	if !aibd.putAndCheck(ai.color, best.loc) {
		panic(fmt.Errorf("cannot put: %v, builtin ai %v", bestPoint, ai.color))
	}
	return bestPoint.toBoardPoint()
}

func (ai *AI8) printValue(best node) {
	if ai.phase == 1 {
		finValue := float64(best.value) / float64(TOTALVALUE) * float64(SIZE8*SIZE8)
		fmt.Printf("built-in AI: {depth: %d, nodes: %d, value: %+.2f}\n", ai.reachedDepth, ai.nodes, finValue)
	} else {
		finValue := best.value
		fmt.Printf("built-in AI: {depth: %d, nodes: %d, value: %+d}\n", ai.reachedDepth, ai.nodes, finValue)
	}
}

func (ai *AI8) setPhase(bd bboard8) {
	emptyCount := bd.emptyCount()
	if emptyCount > PHASE2DEPTH8 {
		ai.phase = 1
	} else {
		ai.phase = 2
	}
}

func (ai *AI8) setDepth() {
	if ai.phase == 1 {
		ai.depth = PHASE1DEPTH8
	} else {
		ai.depth = MAXINT // until end of game
	}
}

func (ai *AI8) heuristic(bd bboard8) int {
	if ai.phase == 1 { // phase 1
		return bd.eval(ai.color)
	} else { // phase 2
		return bd.count(ai.color) - bd.count(ai.opponent)
	}
}

func (ai *AI8) sortedValidNodes(bd bboard8, cl color) (all nodes) {
	all = ai.nodesPool.getClearOne()
	if ai.phase == 1 { // phase 1 sort by eval
		allValid := bd.allValidLoc(cl)
		for loc := 0; loc < SIZE8*SIZE8; loc++ {
			if (u1<<loc)&allValid != 0 {
				tmp := bd.cpy()
				tmp.put(cl, loc)
				all = append(all, node{loc, tmp.eval(cl)})
			}
		}
		all.sortDesc()
	} else { // phase 2 sort by mobility
		op := cl.reverse()
		allValid := bd.allValidLoc(cl)
		for loc := 0; loc < SIZE8*SIZE8; loc++ {
			if (u1<<loc)&allValid != 0 {
				tmp := bd.cpy()
				tmp.put(cl, loc)
				v := tmp.mobility(op)
				all = append(all, node{loc, v})
			}
		}
		// the smaller the opponent's mobility is, the better.
		all.sortAsc()
	}
	return
}

func (ai *AI8) alphaBetaHelper(bd bboard8, depth int) node {
	return ai.alphaBeta(bd, depth, MININT, MAXINT, true)
}

func (ai *AI8) alphaBeta(bd bboard8, depth int, alpha int, beta int, maxLayer bool) node {
	ai.nodes++

	if depth == 0 {
		ai.reachedDepth = ai.depth
		v := ai.heuristic(bd)
		return node{-1, v}
	}
	if bd.isOver() {
		ai.reachedDepth = ai.depth - depth
		v := ai.heuristic(bd)
		return node{-1, v}
	}

	if maxLayer {
		maxValue := MININT
		bestNode := node{-1, maxValue}

		aiValid := ai.sortedValidNodes(bd, ai.color)
		if len(aiValid) == 0 { // 沒地方下，換邊
			ai.nodesPool.freeOne()
			return ai.alphaBeta(bd, depth, alpha, beta, false)
		}

		for i := range aiValid {
			tmp := bd.cpy()
			tmp.put(ai.color, aiValid[i].loc)
			eval := ai.alphaBeta(tmp, depth-1, alpha, beta, false).value

			if eval > maxValue {
				maxValue = eval
				bestNode = aiValid[i]
			}
			alpha = max(alpha, maxValue)
			if beta <= alpha {
				break
			}
		}

		ai.nodesPool.freeOne()
		return node{bestNode.loc, maxValue}
	} else {
		minValue := MAXINT
		bestNode := node{-1, minValue}

		opValid := ai.sortedValidNodes(bd, ai.opponent)
		if len(opValid) == 0 { // 對手沒地方下，換邊
			ai.nodesPool.freeOne()
			return ai.alphaBeta(bd, depth, alpha, beta, true)
		}

		for i := range opValid {
			tmp := bd.cpy()
			tmp.put(ai.opponent, opValid[i].loc)
			eval := ai.alphaBeta(tmp, depth-1, alpha, beta, true).value

			if eval < minValue {
				minValue = eval
				bestNode = opValid[i]
			}

			beta = min(beta, minValue)
			if beta <= alpha {
				break
			}
		}

		ai.nodesPool.freeOne()
		return node{bestNode.loc, minValue}
	}
}
