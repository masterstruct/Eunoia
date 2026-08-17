package genmagics

import (
	"math/rand"

	"github.com/masterstruct/Eunoia/internal/board"
	"github.com/masterstruct/Eunoia/internal/movegen"
)

func findMagic(isRook bool, sq board.Square, indexBits uint8, rng *rand.Rand) (movegen.MagicEntry, []board.Bitboard) {
	var mask board.Bitboard
	if isRook {
		mask = movegen.RookMask(sq)
	} else {
		mask = movegen.BishopMask(sq)
	}

	for {
		magic := rng.Uint64() & rng.Uint64() & rng.Uint64()
		entry := movegen.MagicEntry{Mask: mask, Magic: magic, IndexBits: indexBits}

		table, ok := tryMakeTable(isRook, sq, &entry)
		if ok {
			return entry, table
		}
	}
}

func tryMakeTable(isRook bool, sq board.Square, entry *movegen.MagicEntry) ([]board.Bitboard, bool) {
	tableSize := 1 << entry.IndexBits
	table := make([]board.Bitboard, tableSize)
	used := make([]bool, tableSize)

	var moves board.Bitboard
	for blockers := range movegen.Subsets(entry.Mask) {
		if isRook {
			moves = movegen.RookAttacks(sq, blockers)
		} else {
			moves = movegen.BishopAttacks(sq, blockers)
		}

		index := movegen.MagicIndex(entry, blockers)
		if !used[index] {
			used[index] = true
			table[index] = moves
		} else if table[index] != moves {
			return []board.Bitboard{}, false
		}
	}
	return table, true
}
