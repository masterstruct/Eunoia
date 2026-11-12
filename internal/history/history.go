package history

const (
	maxHistory = 16384
)

type History struct {
	Quiet quietHistory
	Noisy noisyHistory
}

func (h *History) Clear() {
	clear(h.Quiet[:])
	clear(h.Noisy[:])
}

func gravity(entry *int16, bonus, maxBonus, maxValue int32) {
	gravityWithDecay(entry, int32(*entry), bonus, maxBonus, maxValue)
}

func gravityWithDecay(entry *int16, decay, bonus, maxBonus, maxValue int32) {
	bonus = clampInt32(bonus, -maxBonus, maxBonus)
	decay = decay * absInt32(bonus) / maxValue
	*entry += int16(bonus - decay)
}

func clampInt32(v, low, high int32) int32 {
	if v < low {
		return low
	}
	if v > high {
		return high
	}
	return v
}

func absInt32(v int32) int32 {
	if v < 0 {
		return -v
	}
	return v
}
