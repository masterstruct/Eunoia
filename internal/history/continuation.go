package history

import (
	"github.com/masterstruct/Eunoia/internal/board"
)

const (
	cont1BonusBase  = 128
	cont1BonusScale = 128
	cont1MalusBase  = 128
	cont1MalusScale = 128
	cont1Max        = 2048
)

// [stm][prev piece][prev to][piece][to]
type ContHist [2][6][64][6][64]int16

func (h *ContHist) get(stm board.Color, entry contEntry,
	pt board.PieceType, to board.Square) int32 {
	if entry.pt == board.NoPieceType {
		return 0
	}
	return int32(h[stm][entry.pt][entry.to][pt][to])
}

func (h *ContHist) bonus(depth int32) int32 {
	return cont1BonusBase + cont1BonusScale*depth
}

func (h *ContHist) malus(depth int32) int32 {
	return -(cont1MalusBase + cont1MalusScale*depth)
}

func (h *ContHist) update(stm board.Color, entry contEntry,
	pt board.PieceType, to board.Square, bonus int32) {
	if entry.pt != board.NoPieceType {
		gravity(
			&h[stm][entry.pt][entry.to][pt][to],
			bonus,
			cont1Max,
			maxHistory,
		)
	}
}

type contEntry struct {
	pt board.PieceType
	to board.Square
}

func NewContEntry(piece board.Piece, move board.Move) contEntry {
	return contEntry{
		pt: piece.Type(),
		to: move.To(),
	}
}

type ContIndices struct {
	cont1 contEntry
}

func NewContIndices(cont1 contEntry) ContIndices {
	return ContIndices{
		cont1: cont1,
	}
}
