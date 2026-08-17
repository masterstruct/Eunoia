package movegen

import "github.com/masterstruct/Eunoia/internal/board"

func rookMask(sq board.Square) board.Bitboard {
	file := sq.File()
	rank := sq.Rank()
	var mask board.Bitboard

	for f := board.FileB; f < board.FileH; f++ {
		mask.SetBit(board.NewSquare(f, rank))
	}

	for r := board.Rank2; r < board.Rank8; r++ {
		mask.SetBit(board.NewSquare(file, r))
	}

	mask.ClearBit(sq)
	return mask
}

func bishopMask(sq board.Square) board.Bitboard {
	var mask board.Bitboard
	mask = BishopAttacks(sq, board.EmptyBB)
	return mask &^ board.EdgesBB
}

var rookMagics [64]MagicEntry
var bishopMagics [64]MagicEntry
var rookMoves [64][]board.Bitboard
var bishopMoves [64][]board.Bitboard

type MagicEntry struct {
	Mask  board.Bitboard
	Magic uint64
	Shift uint8
}

func MagicIndex(entry *MagicEntry, occupied board.Bitboard) int {
	occupied &= entry.Mask
	hash := uint64(occupied) * entry.Magic
	return int(hash >> entry.Shift)
}
