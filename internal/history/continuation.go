package history

import "github.com/masterstruct/Eunoia/internal/board"

type ContEntry struct {
	pt board.PieceType
	to board.Square
}

func NewContEntry() ContEntry {
	return ContEntry{
		pt: board.NoPieceType,
		to: board.NoSquare,
	}
}

// [stm][prev piece][prev to][piece][to]
type ContHist [2][6][64][6][64]int16

func (h *ContHist) Get(pos *board.Position, entry ContEntry, move board.Move) int32 {
	stm := pos.SideToMove
	piece := pos.Attacker(move)
	return int32(h[stm][entry.pt][entry.to][piece][move.To()])
}

func (h *ContHist) Update(pos *board.Position, entry ContEntry,
	move board.Move, bonus int32) {
	stm := pos.SideToMove
	piece := pos.Attacker(move)
	gravity(&h[stm][entry.pt][entry.to][piece][move.To()],
		bonus, maxHistory, maxHistory)
}
