package builtinai

var (
	VALUE6x6 = [][]int{
		{100, -36, 53, 53, -36, 100},
		{-36, -69, -10, -10, -69, -36},
		{53, -10, -2, -2, -10, 53},
		{53, -10, -2, -2, -10, 53},
		{-36, -69, -10, -10, -69, -36},
		{100, -36, 53, 53, -36, 100},
	}

	VALUE8x8 = [][]int{
		{800, -286, 426, -24, -24, 426, -286, 800},
		{-286, -552, -177, -82, -82, -177, -552, -286},
		{426, -177, 62, 8, 8, 62, -177, 426},
		{-24, -82, 8, -18, -18, 8, -82, -24},
		{-24, -82, 8, -18, -18, 8, -82, -24},
		{426, -177, 62, 8, 8, 62, -177, 426},
		{-286, -552, -177, -82, -82, -177, -552, -286},
		{800, -286, 426, -24, -24, 426, -286, 800},
	}

	TOTAL6x6 int
	TOTAL8x8 int
)

func init() {
	for i := 0; i < len(VALUE6x6); i++ {
		for j := 0; j < len(VALUE6x6); j++ {
			TOTAL6x6 += abs(VALUE6x6[i][j])
		}
	}
	for i := 0; i < len(VALUE8x8); i++ {
		for j := 0; j < len(VALUE8x8); j++ {
			TOTAL8x8 += abs(VALUE8x8[i][j])
		}
	}
}

func (bd aiboard) eval(cl color, opponent color, valueDisk [][]int) int {
	value := 0
	for i := 0; i < bd.size(); i++ {
		for j := 0; j < bd.size(); j++ {
			p := point{x: i, y: j}
			if bd.at(p) == cl {
				value += valueDisk[i][j]
			} else if bd.at(p) == opponent {
				value -= valueDisk[i][j]
			}
		}
	}
	return value
}

func (bd aiboard) countPieces(cl color) int {
	count := 0
	for i := 0; i < bd.size(); i++ {
		for j := 0; j < bd.size(); j++ {
			if bd.at(point{x: i, y: j}) == cl {
				count++
			}
		}
	}
	return count
}

// return the mobility
func (bd aiboard) mobility(cl color) int {
	count := 0
	for i := 0; i < bd.size(); i++ {
		for j := 0; j < bd.size(); j++ {
			if bd.isValidPoint(cl, point{x: i, y: j}) {
				count++
			}
		}
	}
	return count
}

func (ai *AI) changedValue(bd aiboard, cl color, p point, dir [2]int) int {
	delta := 0
	x, y := p.x, p.y
	opponent := cl.reverse()

	x, y = x+dir[0], y+dir[1]
	if bd.at(point{x: x, y: y}) != opponent {
		return 0
	}
	delta += ai.valueDisk[x][y] * 2 // flip opponent to yours, so double

	for {
		x, y = x+dir[0], y+dir[1]
		now := bd.at(point{x: x, y: y})
		if now != opponent {
			if now == cl {
				return delta
			} else {
				return 0
			}
		}
		delta += ai.valueDisk[x][y] * 2 // same as above
	}
}

// don't need to copy
func (ai *AI) evalAfterPut(bd aiboard, currentValue int, p point, cl color) int {
	for i := 0; i < 8; i++ {
		currentValue += ai.changedValue(bd, cl, p, DIRECTION[i])
	}
	currentValue += ai.valueDisk[p.x][p.y]
	return currentValue
}

// don't need to copy board
func (ai *AI) countAfterPut(bd aiboard, currentCount int, p point, cl color) int {
	for i := 0; i < 8; i++ {
		currentCount += bd.countFlipPieces(cl, cl.reverse(), p, DIRECTION[i])
	}
	return currentCount + 1 // include p itself
}

func (ai *AI) heuristicAfterPut(bd aiboard, currentValue int, p point, cl color) int {
	if ai.step == 1 {
		return ai.evalAfterPut(bd, currentValue, p, cl)
	} else {
		return ai.countAfterPut(bd, currentValue, p, ai.color)
	}
}
