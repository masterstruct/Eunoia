package history

import "github.com/masterstruct/Eunoia/internal/board"

// [stm][prev piece][prev to][piece][to]
type ContHist [2][6][64][6][64]int16

func (h *ContHist) Get(stm board.Color, prevPt board.PieceType,
	prevTo board.Square, pt board.PieceType, to board.Square) int32 {
	return int32(h[stm][prevPt][prevTo][pt][to])
}

func (h *ContHist) Update(stm board.Color, prevPt board.PieceType,
	prevTo board.Square, pt board.PieceType, to board.Square, bonus int32) {
	gravity(&h[stm][prevPt][prevTo][pt][to],
		bonus, maxHistory, maxHistory)
}
