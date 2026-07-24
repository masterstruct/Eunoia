package board

type Bitboards struct {
	pieces [6]Bitboard
	colors [2]Bitboard
}

func (bb Bitboards) PieceBB(p Piece) Bitboard {
	return bb.pieces[p.Type] & bb.colors[p.Color]
}

func (bb Bitboards) Occupied() Bitboard {
	return bb.colors[White] | bb.colors[Black]
}

func (bb *Bitboards) PlacePiece(p Piece, sq Square) {
	bb.pieces[p.Type].SetBit(sq)
	bb.colors[p.Color].SetBit(sq)
}

func (bb *Bitboards) RemovePiece(sq Square) {
	// TODO: replace with mailbox lookup

	for _, pt := range PieceTypes() {
		bb.pieces[pt].ClearBit(sq)
	}
	bb.colors[Black].ClearBit(sq)
	bb.colors[White].ClearBit(sq)
}
