package history

import (
	"github.com/masterstruct/Eunoia/internal/board"
)

const (
	cont1BonusBase  = -250
	cont1BonusScale = 300
	cont1MalusBase  = -250
	cont1MalusScale = 300
	cont1Max        = 2048
)

type ContRow [6][64]int16
type ContHist [2][6][64]ContRow // [stm][prev piece][prev to][piece][to]

func (r *ContRow) get(pt board.PieceType, to board.Square) int32 {
	if r == nil {
		return 0
	}
	return int32(r[pt][to])
}

func (r *ContRow) update(pt board.PieceType, to board.Square, bonus int32) {
	if r == nil {
		return
	}
	gravity(&r[pt][to], bonus, cont1Max, maxHistory)
}

func (h *ContHist) bonus(depth int32) int32 {
	return cont1BonusBase + cont1BonusScale*depth
}

func (h *ContHist) malus(depth int32) int32 {
	return -(cont1MalusBase + cont1MalusScale*depth)
}

type ContIndices struct {
	Cont1 *ContRow
}

func NewContIndices(cont1 *ContRow) ContIndices {
	return ContIndices{Cont1: cont1}
}
