package search

import "github.com/masterstruct/Eunoia/internal/board"

type CorrHist [2][corrHistSize]int32

const (
	corrHistMax  = 1024
	corrHistSize = 16384
	corrHistMask = corrHistSize - 1

	corrHistMaxBonus = corrHistMax / 4
)

func (hist *CorrHist) update(stm board.Color, key uint64, bonus int32) {
	index := key & corrHistMask
	entry := &hist[stm][index]
	bonus = min(max(bonus, -corrHistMaxBonus), corrHistMaxBonus)
	*entry = gravity(*entry, bonus, corrHistMax)
}

func (hist *CorrHist) updatePawnCorrHist(pos *board.Position, depth int, staticEval, bestScore int32) {
	stm := pos.SideToMove
	diff := bestScore - staticEval

	key := getPawnKey(pos)
	hist.update(stm, key, diff*int32(depth)/8)
}

func (hist *CorrHist) correction(pos *board.Position) int32 {
	index := getPawnKey(pos) & corrHistMask
	return 128 * hist[pos.SideToMove][index] / 2048
}

func getPawnKey(pos *board.Position) uint64 {
	var key uint64
	wp := pos.PieceBB(board.WhitePawn)
	bp := pos.PieceBB(board.BlackPawn)

	for wp != 0 {
		key ^= board.ZobristTable.PieceKey(board.White, board.Pawn, wp.PopLSB())
	}
	for bp != 0 {
		key ^= board.ZobristTable.PieceKey(board.Black, board.Pawn, bp.PopLSB())
	}
	return key
}
