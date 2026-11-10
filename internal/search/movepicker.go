package search

import (
	"github.com/masterstruct/Eunoia/internal/board"
)

// based on Hobbes implementation
// https://github.com/kelseyde/hobbes-chess-engine/blob/main/src/search/movepicker.rs

const (
	killerBonus = 10_000_000
)

type Stage uint8

const (
	TTMove Stage = iota
	GenerateNoisies
	GoodNoisies
	Quiets
	BadNoisies
)

type MovePicker struct {
	movelist   board.Movelist
	stage      Stage
	ttMove     board.Move
	skipQuiets bool
	ply        uint16
	badNoisies board.Movelist
}

func NewMovePicker(ply uint16, ttMove board.Move, skipQuiets bool) MovePicker {
	stage := TTMove
	if ttMove == board.NullMove {
		stage = GenerateNoisies
	}
	return MovePicker{
		movelist:   board.Movelist{},
		stage:      stage,
		ttMove:     ttMove,
		skipQuiets: skipQuiets,
		ply:        ply,
	}
}

func (mp *MovePicker) Next(pos *board.Position, ss *SearchState) board.Move {
	if mp.stage == TTMove {
		mp.stage = GenerateNoisies
		if pos.IsLegal(mp.ttMove) {
			return mp.ttMove
		}
	}

	if mp.stage == GenerateNoisies {
		mp.stage = GoodNoisies
		mp.genNoisies(pos, ss)
	}

	if mp.stage == GoodNoisies {
		if !mp.movelist.IsEmpty() {
			return pickBest(&mp.movelist)
		}
		if mp.skipQuiets {
			mp.stage = BadNoisies
		} else {
			mp.stage = Quiets
			mp.genQuiets(pos, ss)
		}
	}

	if mp.stage == Quiets {
		if !mp.skipQuiets && !mp.movelist.IsEmpty() {
			return pickBest(&mp.movelist)
		}
		mp.stage = BadNoisies
	}

	if mp.stage == BadNoisies {
		if !mp.badNoisies.IsEmpty() {
			return pickBest(&mp.badNoisies)
		}
	}

	return board.NullMove
}

func pickBest(movelist *board.Movelist) board.Move {
	bestIndex := 0
	bestScore := movelist.Moves[0].Score

	for i := 1; i < movelist.Len; i++ {
		score := movelist.Moves[i].Score
		if score > bestScore {
			bestScore = score
			bestIndex = i
		}
	}
	bestMove := movelist.Moves[bestIndex].Move
	movelist.Remove(bestIndex)
	return bestMove
}

func (mp *MovePicker) genNoisies(pos *board.Position, ss *SearchState) {
	var temp board.Movelist
	board.GenerateLegalMoves(pos, &temp, board.Noisies)
	mp.removeTT(&temp)

	for i := range temp.Len {
		var score int32
		move := temp.Moves[i].Move

		// TODO: add queen promo bonus

		// MVV-LVA
		if move.IsCapture() {
			score += mvvlvaScore(pos.Victim(move), pos.Attacker(move))
		}

		scoredMove := move.Score(score)

		if !isGoodNoisy(pos, move) {
			mp.badNoisies.AddScoredMove(scoredMove)
			continue
		}

		mp.movelist.AddScoredMove(scoredMove)
	}
}

func (mp *MovePicker) genQuiets(pos *board.Position, ss *SearchState) {
	stm := pos.SideToMove

	var temp board.Movelist
	board.GenerateLegalMoves(pos, &temp, board.Quiets)
	mp.removeTT(&temp)

	for i := range temp.Len {
		move := temp.Moves[i]
		from := move.Move.From()
		to := move.Move.To()

		move.Score = ss.history.Quiet[stm][from][to]

		// killer move bonus
		if ss.stack[mp.ply].killer == move.Move {
			move.Score += killerBonus
		}
		mp.movelist.AddScoredMove(move)
	}
}

func (mp *MovePicker) removeTT(movelist *board.Movelist) {
	for i := range movelist.Len {
		if movelist.Moves[i].Move == mp.ttMove {
			movelist.Remove(i)
			return
		}
	}
}

// formula from https://asteri.sm/files/2023-02-20-viri-wiki#mvvlva
func mvvlvaScore(victim, attacker board.PieceType) int32 {
	return int32(victim)*1000 + 60 - int32(attacker)*10
}

func isGoodNoisy(pos *board.Position, move board.Move) bool {
	if move.IsPromo() && move.Promo() == board.Queen {
		return true
	}
	threshold := int32(-100)
	return see(pos, move, threshold)
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
