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
		entry := movegen.MagicEntry{Mask: mask, Magic: magic, Shift: 64 - indexBits}

		table, ok := movegen.TryMakeTable(isRook, sq, &entry)
		if ok {
			return entry, table
		}
	}
}

func writeMagics(name string, magics *[64]movegen.MagicEntry, w io.Writer) {
	fmt.Fprintf(w, "var %vMagics = [64]MagicEntry{\n", name)
	for sq := board.A1; sq <= board.H8; sq++ {
		entry := magics[sq]
		fmt.Fprintf(w, "    {Magic: 0x%016x, Shift: %v},\n", entry.Magic, entry.Shift)
	}
	fmt.Fprintf(w, "}\n\n")
}

func main() {
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))

	fmt.Println("Finding magics for bishops...")
	for sq := board.A1; sq <= board.H8; sq++ {
		mask := movegen.BishopMask(sq)

		entry, _ := findMagic(false, sq, uint8(mask.CountBits()), rng)

		movegen.BishopMagics[sq] = entry
	}

	fmt.Println("Finding magics for rooks...")
	for sq := board.A1; sq <= board.H8; sq++ {
		mask := movegen.RookMask(sq)

		entry, _ := findMagic(true, sq, uint8(mask.CountBits()), rng)
		movegen.RookMagics[sq] = entry
	}

	writeMagics("Rook", &movegen.RookMagics, os.Stdout)
	writeMagics("Bishop", &movegen.BishopMagics, os.Stdout)
}
