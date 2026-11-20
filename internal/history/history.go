package history

import (
	"github.com/masterstruct/Eunoia/internal/board"
)

const (
	maxHistory = 16384
)

type History struct {
	Quiet   quietHistory
	contOdd ContHist
}

func (h *History) Clear() {
	clear(h.Quiet[:])
	clear(h.contOdd[:])
}

func (h *History) ScoreQuiet(pos *board.Position, move board.Move, indices ContIndices) int32 {
	return h.Quiet.Get(pos.SideToMove, move.From(), move.To()) + h.cont(pos, indices, move)
}

func (h *History) Update(pos *board.Position, indices ContIndices,
	depth int32, bestMove board.Move, quietsTried *board.Movelist) {
	// TODO: replace movelist with a cheaper array
	if bestMove.IsNoisy() {
		// TODO: noisy history
	} else {
		h.applyQuietBonus(pos, indices, depth, bestMove)
		for i := range quietsTried.Len {
			// penalize quiets that didn't cause beta cutoff
			quiet := quietsTried.Moves[i].Move
			h.applyQuietMalus(pos, indices, depth, quiet)
		}
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
