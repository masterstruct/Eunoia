package search

import (
	"github.com/masterstruct/Eunoia/internal/board"
)

const (
	ttMoveBonus  = 1_000_000
	captureBonus = 100_000

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
	}
}

func (mp *MovePicker) genNoisies(pos *board.Position, ss *SearchState) {
	stm := pos.SideToMove

	var temp board.Movelist
	board.GenerateLegalMoves(pos, &temp, board.Noisies)

	for i := range temp.Len {
		move := temp.Moves[i].Move
		if move == mp.ttMove {
			continue
		}

		from := move.From()
		to := move.To()

		// butterfly history
		score := ss.butterflyHistory[stm][from][to]

		// TODO: add queen promo bonus

		// capture bonus and MVV-LVA
		if move.IsCapture() {
			victim := pos.Board[to].Type
			if move.IsEnPassant() {
				victim = board.Pawn
			}

			// TODO: remove capture bonus - no point because no longer mixing noisies/quiets
			score += captureBonus + mvvlvaScore(victim, pos.Board[from].Type)
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
		if move.Move == mp.ttMove {
			continue
		}

		from := move.Move.From()
		to := move.Move.To()

		move.Score = ss.butterflyHistory[stm][from][to]
		mp.movelist.AddScoredMove(move)
	}
}

func (ss *SearchState) orderMoves(pos *board.Position, movelist *board.Movelist) {
	n := movelist.Len
	if n == 0 {
		return
	}

	stm := pos.SideToMove

	// TT lookup
	entry, ttHit := ss.tt.Probe(pos.Hash)
	var ttMove board.Move
	if ttHit {
		ttMove = entry.Move
	}

	// score moves
	var scores [board.MaxMoves]int32
	for i := range n {
		move := movelist.Moves[i].Move
		from := move.From()
		to := move.To()

		score := ss.butterflyHistory[stm][from][to]

		if ttHit && move == ttMove {
			score += ttMoveBonus
		}

		if move.IsCapture() {
			attacker, _ := pos.PieceOn(from)
			victim, _ := pos.PieceOn(to)

			victimType := victim.Type
			if move.IsEnPassant() {
				victimType = board.Pawn
			}

			score += captureBonus + mvvlvaScore(victimType, attacker.Type)
		}

		scores[i] = score
	}

	// reverse insertion sort
	for i := 1; i < n; i++ {
		score := scores[i]
		move := movelist.Moves[i].Move

		j := i - 1
		for j >= 0 && scores[j] < score {
			scores[j+1] = scores[j]
			movelist.Moves[j+1] = movelist.Moves[j]
			j--
		}

		scores[j+1] = score
		movelist.Moves[j+1] = move.ScoredMove(0)
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
