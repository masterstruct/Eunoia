package movegen

import (
	"fmt"

	"github.com/masterstruct/Eunoia/internal/board"
)

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

func init() {
	initMagics(true, &RookMagics, &RookMoves, "rook")
	initMagics(false, &BishopMagics, &BishopMoves, "bishop")
}

func initMagics(isRook bool, magics *[64]MagicEntry, moves *[64][]board.Bitboard, name string) {
	for sq := board.A1; sq <= board.H8; sq++ {
		entry := &magics[sq]

		table, ok := TryMakeTable(isRook, sq, entry)
		if !ok {
			panic(fmt.Sprintf("Bad %s magics! Square %v", name, sq))
		}
		moves[sq] = table
	}
}

func TryMakeTable(isRook bool, sq board.Square, entry *MagicEntry) ([]board.Bitboard, bool) {
	tableSize := 1 << entry.IndexBits
	table := make([]board.Bitboard, tableSize)
	used := make([]bool, tableSize)

	var moves board.Bitboard
	for blockers := range Subsets(entry.Mask) {
		if isRook {
			moves = RookAttacks(sq, blockers)
		} else {
			moves = BishopAttacks(sq, blockers)
		}

		index := MagicIndex(entry, blockers)
		if !used[index] {
			used[index] = true
			table[index] = moves
		} else if table[index] != moves {
			return []board.Bitboard{}, false
		}
	}
	return table, true
}
