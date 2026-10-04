package search

import (
	"math"
)

const (
	lmrMinDepth = 3
	lmrMinMoves = 2
	lmrMaxDepth = 63
	lmrMaxMoves = 63

	lmrBase    = 0.99
	lmrDivisor = 3.14
)

var lmr [lmrMaxDepth + 1][lmrMaxMoves + 1]int

func init() {
	initLMRTable()
}

func initLMRTable() {
	for depth := range lmr {
		for move := range lmr[depth] {
			if depth < lmrMinDepth-1 || move < lmrMinMoves {
				continue
			}
			lmr[depth][move] = computeLMR(depth, move)
		}
	}
}

func computeLMR(depth, move int) int {
	return int(lmrBase + math.Log(float64(depth))*math.Log(float64(move))/lmrDivisor)
}
