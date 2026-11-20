package history

import "github.com/masterstruct/Eunoia/internal/board"

const (
	noisyBonusBase  = -250
	noisyBonusScale = 300

	noisyMalusBase  = -250
	noisyMalusScale = 300
)

type noisyHistory [2][6][7][64]int16 // [stm][attacker][victim][to]

func (h *noisyHistory) get(stm board.Color, attacker, victim board.PieceType, to board.Square) int32 {
	return int32(h[stm][attacker][victim][to])
}

func (h *noisyHistory) update(stm board.Color, attacker, victim board.PieceType, to board.Square, bonus int32) {
	gravity(
		&h[stm][attacker][victim][to],
		bonus,
		maxHistory,
		maxHistory,
	)
}

func (h *noisyHistory) bonus(depth int32) int32 {
	return noisyBonusBase + noisyBonusScale*depth
}

func (h *noisyHistory) malus(depth int32) int32 {
	return -(noisyMalusBase + noisyMalusScale*depth)
}
