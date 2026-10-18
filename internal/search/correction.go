package search

import "github.com/masterstruct/Eunoia/internal/board"

type CorrHist [2][16384]int32

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
