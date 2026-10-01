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
	var scores [board.MaxMoves]int
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
		movelist.Moves[j+1] = board.NewScoredMove(move, 0)
	}
}

// formula from https://asteri.sm/files/2023-02-20-viri-wiki#mvvlva
func mvvlvaScore(victim, attacker board.PieceType) int {
	return int(victim)*1000 + 60 - int(attacker)*10
}

func (ss *SearchState) updateButterflyHistory(stm board.Color, from, to board.Square, bonus int) {
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
