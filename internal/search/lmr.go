package search

import (
	"math"

	"github.com/masterstruct/Eunoia/internal/movegen"
)

var lmr [MaxPly][movegen.MaxMoves]int // [depth][move]

func init() {
	initLMRTable()
}

func initLMRTable() {
	for depth := range MaxPly {
		for move := range movegen.MaxMoves {
			lmr[depth][move] = computeLMR(depth, move)
		}
	}
}

func computeLMR(depth, move int) int {
	return int(0.99 + math.Log(float64(depth))*math.Log(float64(move))/3.14)
}
