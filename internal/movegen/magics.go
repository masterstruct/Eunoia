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
