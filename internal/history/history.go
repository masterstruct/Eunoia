package history

import "github.com/masterstruct/Eunoia/internal/board"

const (
	maxHistory = 16384
)

type History struct {
	Quiet   quietHistory
	noisy   noisyHistory
	contOdd ContHist
}

func (h *History) Clear() {
	clear(h.Quiet[:])
	clear(h.noisy[:])
	clear(h.contOdd[:])
}

func (h *History) ScoreQuiet(pos *board.Position, move board.Move, indices ContIndices) int32 {
	return h.Quiet.Get(pos.SideToMove, move.From(), move.To()) + h.cont(pos, indices, move)
}

func (h *History) ScoreNoisy(pos *board.Position, move board.Move) int32 {
	return h.noisy.get(pos.SideToMove, pos.Attacker(move), pos.Victim(move), move.To())
}

func (h *History) Update(pos *board.Position, indices ContIndices,
	depth int32, move board.Move, quietsTried, noisiesTried *board.Movelist) {
	// TODO: replace movelist with a cheaper array
	stm := pos.SideToMove
	if move.IsNoisy() {
		h.noisy.update(stm, pos.Attacker(move),
			pos.Victim(move), move.To(), h.noisy.bonus(depth))
	} else {
		h.applyQuietBonus(pos, indices, depth, move)

		// penalize quiets that didn't cause beta cutoff
		for i := range quietsTried.Len {
			quiet := quietsTried.Moves[i].Move
			h.applyQuietMalus(pos, indices, depth, quiet)
		}
	}

	// penalize all noisies that didn't cause beta cutoff
	// regardless if the best move is quiet or noisy
	noisyMalus := h.noisy.malus(depth)
	for i := range noisiesTried.Len {
		noisyMove := noisiesTried.Moves[i].Move
		h.noisy.update(stm, pos.Attacker(noisyMove),
			pos.Victim(noisyMove), noisyMove.To(), noisyMalus)
	}
}

func (h *History) applyQuietBonus(pos *board.Position, indices ContIndices,
	depth int32, move board.Move) {
	stm := pos.SideToMove
	from := move.From()
	to := move.To()
	attacker := pos.Attacker(move)

	h.Quiet.update(stm, from, to, h.Quiet.bonus(depth))
	indices.Cont1.update(attacker, to, h.contOdd.bonus(depth))
}

func (h *History) applyQuietMalus(pos *board.Position, indices ContIndices,
	depth int32, move board.Move) {
	stm := pos.SideToMove
	from := move.From()
	to := move.To()
	attacker := pos.Attacker(move)

	h.Quiet.update(stm, from, to, h.Quiet.malus(depth))
	indices.Cont1.update(attacker, to, h.contOdd.malus(depth))
}

func (h *History) Cont1Row(stm board.Color, pt board.PieceType, to board.Square) *ContRow {
	return &h.contOdd[stm][pt][to]
}

func (h *History) cont(pos *board.Position, indices ContIndices, move board.Move) int32 {
	return indices.Cont1.get(pos.Attacker(move), move.To())
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
