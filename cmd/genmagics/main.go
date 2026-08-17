package main

import (
	"fmt"
	"io"
	"math/rand"
	"os"
	"time"

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

func writeMagics(name string, magics *[64]movegen.MagicEntry, w io.Writer) {
	fmt.Fprintf(w, "var %vMagics = [64]MagicEntry{\n", name)
	for sq := board.A1; sq <= board.H8; sq++ {
		entry := magics[sq]
		fmt.Fprintf(w, "    {Mask: 0x%016x, Magic: 0x%016x, IndexBits: %v},\n", uint64(entry.Mask), entry.Magic, entry.IndexBits)
	}
	fmt.Fprintf(w, "}\n\n")
}

func main() {
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))

	fmt.Println("Finding magics for bishops...")
	for sq := board.A1; sq <= board.H8; sq++ {
		mask := movegen.RookMask(sq)

		entry, table := findMagic(false, sq, uint8(mask.CountBits()), rng)

		movegen.BishopMagics[sq] = entry
		movegen.BishopMoves[sq] = table
	}

	fmt.Println("Finding magics for rooks...")
	for sq := board.A1; sq <= board.H8; sq++ {
		mask := movegen.RookMask(sq)

		entry, table := findMagic(true, sq, uint8(mask.CountBits()), rng)

		movegen.RookMagics[sq] = entry
		movegen.RookMoves[sq] = table
	}

	writeMagics("Rook", &movegen.RookMagics, os.Stdout)
	writeMagics("Bishop", &movegen.BishopMagics, os.Stdout)
}
