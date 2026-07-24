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

func (bb Bitboards) PieceOn(sq Square) (Piece, bool) {
	// TODO: replace with mailbox lookup
	sqBB := SquareBB[sq]

	if bb.Occupied()&sqBB == 0 {
		return NoPiece, false
	}

	color := White
	if bb.colors[Black]&sqBB != 0 {
		color = Black
	}

	switch {
	case bb.pieces[Pawn]&sqBB != 0:
		return NewPiece(Pawn, color), true
	case bb.pieces[Knight]&sqBB != 0:
		return NewPiece(Knight, color), true
	case bb.pieces[Bishop]&sqBB != 0:
		return NewPiece(Bishop, color), true
	case bb.pieces[Rook]&sqBB != 0:
		return NewPiece(Rook, color), true
	case bb.pieces[Queen]&sqBB != 0:
		return NewPiece(Queen, color), true
	case bb.pieces[King]&sqBB != 0:
		return NewPiece(King, color), true
	default:
		return NoPiece, false
	}
}
