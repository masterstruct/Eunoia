package search

import "github.com/masterstruct/Eunoia/internal/board"

func (ss *SearchState) isDraw(pos *board.Position) bool {
	return pos.HalfmoveClock >= 100 || ss.isRepetition(pos.Hash, pos.HalfmoveClock)
}

func (ss *SearchState) isRepetition(hash uint64, halfmoveClock uint8) bool {
	n := len(ss.keyHistory)
	if n == 0 {
		return false
	}
	cur := n - 1
	reps := 0
	for d := 4; d <= int(halfmoveClock) && d <= cur; d += 2 {
		idx := cur - d
		if ss.keyHistory[idx] != hash {
			continue
		}
		if idx >= ss.rootHistLen {
			return true
		}
		reps++
		if reps >= 2 {
			return true
		}
	}
	return false
}
