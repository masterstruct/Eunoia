package board

type Bitboards struct {
	pieces [6]Bitboard
	colors [2]Bitboard
}

func (pos Position) PieceBB(p Piece) Bitboard {
	return pos.pieces[p.Type] & pos.colors[p.Color]
}

func (pos Position) Occupied() Bitboard {
	return pos.colors[White] | pos.colors[Black]
}

func (pos *Position) PlacePiece(p Piece, sq Square) {
	pos.pieces[p.Type].SetBit(sq)
	pos.colors[p.Color].SetBit(sq)
	pos.Board[sq] = p
}

func (pos *Position) RemovePiece(sq Square) {
	// TODO: replace with mailbox lookup

	for _, pt := range PieceTypes() {
		pos.pieces[pt].ClearBit(sq)
	}
	pos.colors[Black].ClearBit(sq)
	pos.colors[White].ClearBit(sq)
}

func (pos Position) PieceOn(sq Square) (Piece, bool) {
	// TODO: replace with mailbox lookup
	sqBB := SquareBB[sq]

	if pos.Occupied()&sqBB == 0 {
		return NoPiece, false
	}

	color := White
	if pos.colors[Black]&sqBB != 0 {
		color = Black
	}

	switch {
	case pos.pieces[Pawn]&sqBB != 0:
		return NewPiece(Pawn, color), true
	case pos.pieces[Knight]&sqBB != 0:
		return NewPiece(Knight, color), true
	case pos.pieces[Bishop]&sqBB != 0:
		return NewPiece(Bishop, color), true
	case pos.pieces[Rook]&sqBB != 0:
		return NewPiece(Rook, color), true
	case pos.pieces[Queen]&sqBB != 0:
		return NewPiece(Queen, color), true
	case pos.pieces[King]&sqBB != 0:
		return NewPiece(King, color), true
	default:
		return NoPiece, false
	}
}

type Position struct {
	Bitboards
	Board      [64]Piece
	SideToMove Color
	// CastlingRights CastlingRights
	EnPassant     Square
	HalfmoveClock uint8
	Ply           uint16
	Hash          uint64
}

func NewBoard() [64]Piece {
	var board [64]Piece
	for sq := range board {
		board[sq] = NoPiece
	}
	return board
}

func NewPosition() Position {
	pos := Position{SideToMove: White, Ply: 1}
	pos.Board = NewBoard()
	return pos
}
