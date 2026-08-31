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
	var scores [movegen.MaxMoves]int
	for i := range n {
		move := movelist.Moves[i]

		var score int

		if ttHit && move == ttMove {
			score += 32000
		}

		if move.IsCapture() {
			attacker, _ := pos.PieceOn(move.From())
			victim, _ := pos.PieceOn(move.To())

			victimType := victim.Type
			if move.IsEnPassant() {
				victimType = board.Pawn
			}

			score += mvvlvaScore(victimType, attacker.Type)
		}

		scores[i] = score
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
func mvvlvaScore(victim, attacker board.PieceType) int {
	return int(victim)*1000 + 60 - int(attacker)*10
}
