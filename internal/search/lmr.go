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

var lmrTable [lmrMaxDepth + 1][lmrMaxMoves + 1]int32

func init() {
	initLMRTable()
}

func initLMRTable() {
	for depth := range lmrTable {
		for move := range lmrTable[depth] {
			if depth < lmrMinDepth-1 || move < lmrMinMoves {
				continue
			}
			lmrTable[depth][move] = computeLMR(depth, move)
		}
	}
}

func computeLMR(depth, move int) int32 {
	return int32(lmrBase + math.Log(float64(depth))*math.Log(float64(move))/lmrDivisor)
}

func getLMR(depth, searchedMoves int32) int32 {
	return lmrTable[min(depth, lmrMaxDepth)][min(searchedMoves, lmrMaxMoves)]
}
