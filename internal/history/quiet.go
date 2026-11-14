package history

import "github.com/masterstruct/Eunoia/internal/board"

const (
	quietBonusBase  = -250
	quietBonusScale = 300

	quietMalusBase  = -250
	quietMalusScale = 300
)

type quietHistory [2][64][64]int16

func (h *quietHistory) Get(stm board.Color, from, to board.Square) int32 {
	return int32(h[stm][from][to])
}

func (h *quietHistory) applyBonus(stm board.Color, from, to board.Square, depth int32) {
	gravity(
		&h[stm][from][to],
		quietBonusBase+quietBonusScale*depth,
		maxHistory,
		maxHistory,
	)
}

func (h *quietHistory) applyMalus(stm board.Color, from, to board.Square, depth int32) {
	gravity(
		&h[stm][from][to],
		-(quietMalusBase + quietMalusScale*depth),
		maxHistory,
		maxHistory,
	)
}
