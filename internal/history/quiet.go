package history

import "github.com/masterstruct/Eunoia/internal/board"

const (
	quietBonusBase  = -250
	quietBonusScale = 300

	quietMalusBase  = -250
	quietMalusScale = 300
)

type quietHistory [2][64][64]int16 // [stm][from][to]

func (h *quietHistory) Get(stm board.Color, from, to board.Square) int32 {
	return int32(h[stm][from][to])
}

func (h *quietHistory) bonus(depth int32) int32 {
	return quietBonusBase + quietBonusScale*depth
}

func (h *quietHistory) malus(depth int32) int32 {
	return -(quietMalusBase + quietMalusScale*depth)
}

func (h *quietHistory) update(stm board.Color, from, to board.Square, bonus int32) {
	gravity(
		&h[stm][from][to],
		bonus,
		maxHistory,
		maxHistory,
	)
}
