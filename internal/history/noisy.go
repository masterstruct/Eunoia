package history

import "github.com/masterstruct/Eunoia/internal/board"

type noisyHistory [2][6][7][64]int16 // [stm][attacker][victim][to]

func (h *noisyHistory) Get(stm board.Color, attacker, victim board.PieceType, to board.Square) int32 {
	return int32(h[stm][attacker][victim][to])
}

func (h *noisyHistory) Update(stm board.Color, attacker, victim board.PieceType, to board.Square, bonus int32) {
	gravity(&h[stm][attacker][victim][to], bonus, maxHistory, maxHistory)
}
