package search

import (
	"github.com/masterstruct/Eunoia/internal/board"
)

// based on Hobbes implementation
// https://github.com/kelseyde/hobbes-chess-engine/blob/main/src/search/movepicker.rs

const (
	maxHistory = 2 << 13
)

type Stage uint8

const (
	TTMove Stage = iota
	GenerateNoisies
	Noisies
	Quiets
)

type MovePicker struct {
	movelist board.Movelist
	stage    Stage
	ttMove   board.Move
}

func NewMovePicker(ttMove board.Move) MovePicker {
	stage := TTMove
	if ttMove == board.NullMove {
		stage = GenerateNoisies
	}
	return MovePicker{
		movelist: board.Movelist{},
		stage:    stage,
		ttMove:   ttMove,
	}
}

func (mp *MovePicker) Next(pos *board.Position, ss *SearchState, skipQuiets bool) board.Move {
	if mp.stage == TTMove {
		mp.stage = GenerateNoisies
		if pos.IsLegal(mp.ttMove) {
			return mp.ttMove
		}
	}

	if mp.stage == GenerateNoisies {
		mp.stage = Noisies
		mp.genNoisies(pos, ss)
		mp.removeTT()
	}

	if mp.stage == Noisies {
		if !mp.movelist.IsEmpty() {
			return mp.pickBest()
		}
		if !skipQuiets {
			mp.stage = Quiets
			mp.genQuiets(pos, ss)
			mp.removeTT()
		}
	}

	if mp.stage == Quiets {
		if !skipQuiets && !mp.movelist.IsEmpty() {
			return mp.pickBest()
		}
	}

	return board.NullMove
}

func (mp *MovePicker) pickBest() board.Move {
	bestIndex := 0
	bestScore := mp.movelist.Moves[0].Score

	for i := 1; i < mp.movelist.Len; i++ {
		score := mp.movelist.Moves[i].Score
		if score > bestScore {
			bestScore = score
			bestIndex = i
		}
	}
	bestMove := mp.movelist.Moves[bestIndex].Move
	mp.movelist.Remove(bestIndex)
	return bestMove
}

func (mp *MovePicker) genNoisies(pos *board.Position, ss *SearchState) {
	var temp board.Movelist
	board.GenerateLegalMoves(pos, &temp, board.Noisies)

	for i := range temp.Len {
		move := temp.Moves[i].Move
		from := move.From()
		to := move.To()

		var score int32

		// TODO: add queen promo bonus

		// MVV-LVA
		if move.IsCapture() {
			victim := pos.Board[to].Type
			if move.IsEnPassant() {
				victim = board.Pawn
			}

			score += mvvlvaScore(victim, pos.Board[from].Type)
		}

		mp.movelist.AddScoredMove(move.ScoredMove(score))
	}
}

func (mp *MovePicker) genQuiets(pos *board.Position, ss *SearchState) {
	stm := pos.SideToMove

	var temp board.Movelist
	board.GenerateLegalMoves(pos, &temp, board.Quiets)

	for i := range temp.Len {
		move := temp.Moves[i]
		from := move.Move.From()
		to := move.Move.To()

		move.Score = ss.butterflyHistory[stm][from][to]
		mp.movelist.AddScoredMove(move)
	}
}

func (mp *MovePicker) removeTT() {
	for i := range mp.movelist.Len {
		if mp.movelist.Moves[i].Move == mp.ttMove {
			mp.movelist.Remove(i)
			return
		}
	}
}

// formula from https://asteri.sm/files/2023-02-20-viri-wiki#mvvlva
func mvvlvaScore(victim, attacker board.PieceType) int32 {
	return int32(victim)*1000 + 60 - int32(attacker)*10
}

func (ss *SearchState) updateButterflyHistory(stm board.Color, from, to board.Square, bonus int32) {
	// history gravity
	// https://chessprogramming.org/History_Heuristic#history-bonuses

	// 	_____________________
	// /        scaler       \
	// | 25k: -1.61 +- 3.63  |
	// | STC: -0.97 +- 9.20  |
	// \ LTC: 35.39 +- 12.74 /
	//  ---------------------
	//         \   ^__^
	//          \  (oo)\_______
	//             (__)\       )\/\
	//                 ||----w |
	//                 ||     ||

	clampedBonus := bonus
	if clampedBonus > maxHistory {
		clampedBonus = maxHistory
	} else if clampedBonus < -maxHistory {
		clampedBonus = -maxHistory
	}

	absBonus := clampedBonus
	if absBonus < 0 {
		absBonus = -absBonus
	}

	ss.butterflyHistory[stm][from][to] += clampedBonus - ss.butterflyHistory[stm][from][to]*absBonus/maxHistory
}
