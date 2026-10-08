package search

import (
	"github.com/masterstruct/Eunoia/internal/board"
)

func evaluate(pos *board.Position) int32 {
	var score int32

	// PSQT
	score += int32(evaluatePSQT(pos))

	if pos.SideToMove == board.Black {
		return -score
	}
	return score
}
