package movegen

import "github.com/masterstruct/Eunoia/internal/board"

func RookMask(sq board.Square) board.Bitboard {
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

func BishopMask(sq board.Square) board.Bitboard {
	var mask board.Bitboard
	mask = BishopAttacks(sq, board.EmptyBB)
	return mask &^ board.EdgesBB
}

var RookMagics [64]MagicEntry
var BishopMagics [64]MagicEntry
var RookMoves [64][]board.Bitboard
var BishopMoves [64][]board.Bitboard

type MagicEntry struct {
	Mask      board.Bitboard
	Magic     uint64
	IndexBits uint8
}

func MagicIndex(entry *MagicEntry, occupied board.Bitboard) int {
	occupied &= entry.Mask
	hash := uint64(occupied) * entry.Magic
	return int(hash >> (64 - entry.IndexBits))
}

// iterate over all subsets of a bitboard
func Subsets(mask board.Bitboard) func(yield func(board.Bitboard) bool) {
	return func(yield func(board.Bitboard) bool) {
		subset := mask
		for {
			if !yield(subset) {
				return
			}
			if subset == 0 {
				return
			}
			subset = (subset - 1) & mask
		}
	}
}
