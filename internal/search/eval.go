package search

import (
	"github.com/masterstruct/Eunoia/internal/board"
)

func evaluate(pos *board.Position) int16 {
	var score int

	// PSQT
	score += evaluatePSQT(pos)

	if pos.SideToMove == board.Black {
		return int16(-score)
	}
	return int16(score)
}
