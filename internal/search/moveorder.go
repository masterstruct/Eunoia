package search

import (
	"github.com/masterstruct/Eunoia/internal/board"
	"github.com/masterstruct/Eunoia/internal/movegen"
)

func (ss *SearchState) orderMoves(pos *board.Position, movelist *movegen.Movelist) {
	n := movelist.Len
	if n == 0 {
		return
	}

	// TT lookup
	entry, ttHit := ss.tt.Probe(pos.Hash)
	var ttMove board.Move
	if ttHit {
		ttMove = entry.Move
	}

	// score moves
	var scores [movegen.MaxMoves]int16
	for i := range n {
		move := movelist.Moves[i]

		if ttHit && move == ttMove {
			scores[i] = 32000
			continue
		}

		if move.IsCapture() {
			attacker, _ := pos.PieceOn(move.From())
			victim, _ := pos.PieceOn(move.To())

			victimType := victim.Type
			if move.IsEnPassant() {
				victimType = board.Pawn
			}

			scores[i] = mvvlvaScore(victimType, attacker.Type)
		}
	}

	// reverse insertion sort
	for i := 1; i < n; i++ {
		score := scores[i]
		move := movelist.Moves[i]

		j := i - 1
		for j >= 0 && scores[j] < score {
			scores[j+1] = scores[j]
			movelist.Moves[j+1] = movelist.Moves[j]
			j--
		}

		scores[j+1] = score
		movelist.Moves[j+1] = move
	}
}

// formula from https://asteri.sm/files/2023-02-20-viri-wiki#mvvlva
func mvvlvaScore(victim, attacker board.PieceType) int16 {
	return int16(victim)*1000 + 60 - int16(attacker)*10
}
