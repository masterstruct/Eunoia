package history

import "github.com/masterstruct/Eunoia/internal/board"

type quietHistory [2][64][64]int16 // [stm][from][to]

func (h *quietHistory) Get(stm board.Color, from, to board.Square) int32 {
	return int32(h[stm][from][to])
}

func (h *quietHistory) Update(stm board.Color, from, to board.Square, bonus int32) {
	gravity(&h[stm][from][to], bonus, maxHistory, maxHistory)
}
