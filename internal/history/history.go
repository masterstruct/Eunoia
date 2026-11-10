package history

import "github.com/masterstruct/Eunoia/internal/board"

const (
	maxHistory = 16384
)

type History struct {
	Quiet [2][64][64]int32
}

func (h *History) Clear() {
	clear(h.Quiet[:])
}

func gravity(current, update, max int32) int32 {
	return gravityBase(current, update, current, max)
}

func gravityBase(current, update, base, max int32) int32 {
	absUpdate := update
	if absUpdate < 0 {
		absUpdate = -absUpdate
	}
	return current + update - base*absUpdate/max
}

func (h *History) UpdateQuietHistory(stm board.Color, from, to board.Square, bonus int32) {
	clampedBonus := min(max(bonus, -maxHistory), maxHistory) // clamp
	h.Quiet[stm][from][to] =
		gravity(h.Quiet[stm][from][to], clampedBonus, maxHistory)
}
