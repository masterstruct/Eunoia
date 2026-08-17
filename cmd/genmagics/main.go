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
	for sq := board.A1; sq <= board.H8; sq += 2 {
		entry, entry2 := magics[sq], magics[sq+1]
		fmt.Fprintf(w, "    {Magic: 0x%016x}, {Magic: 0x%016x},\n", entry.Magic, entry2.Magic)
	}
	fmt.Fprintf(w, "}\n\n")
}

func main() {
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))

	fmt.Println("Finding magics for bishops...")
	var bishopSize uint64
	for sq := board.A1; sq <= board.H8; sq++ {
		mask := movegen.BishopMask(sq)

		indexBits := uint8(mask.CountBits())
		entry, _ := findMagic(false, sq, indexBits, rng)
		bishopSize += (1 << indexBits)

		movegen.BishopMagics[sq] = entry
	}

	var rookSize uint64
	fmt.Println("Finding magics for rooks...")
	for sq := board.A1; sq <= board.H8; sq++ {
		mask := movegen.RookMask(sq)

		indexBits := uint8(mask.CountBits())
		entry, _ := findMagic(true, sq, indexBits, rng)
		rookSize += (1 << indexBits)

		movegen.RookMagics[sq] = entry
	}

	fmt.Fprintf(os.Stdout, "const RookTableSize = %v\n", rookSize)
	fmt.Fprintf(os.Stdout, "const BishopTableSize = %v\n\n", bishopSize)
	writeMagics("Rook", &movegen.RookMagics, os.Stdout)
	writeMagics("Bishop", &movegen.BishopMagics, os.Stdout)
}
