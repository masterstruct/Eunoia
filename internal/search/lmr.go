package search

import (
	"math"

	"github.com/masterstruct/Eunoia/internal/movegen"
)

var lmr [MaxPly + 1][movegen.MaxMoves]int // [depth][move]

func init() {
	initLMRTable()
}

func initLMRTable() {
	for depth := 0; depth <= MaxPly; depth++ {
		for move := range movegen.MaxMoves {
			if depth < lmrMinDepth || move < lmrMinMoves {
				continue
			}
			lmr[depth][move] = computeLMR(depth, move)
		}
	}
}

func computeLMR(depth, move int) int {
	return int(0.99 + math.Log(float64(depth))*math.Log(float64(move))/3.14)
}
