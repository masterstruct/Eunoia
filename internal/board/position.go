package board

type Bitboards struct {
	Pieces [6]Bitboard
	Colors [2]Bitboard
}

type Position struct {
	Bitboards
	Board         [64]Piece
	SideToMove    Color
	Castling      Castling
	EnPassant     Square
	HalfmoveClock uint8
	Chess960      bool
	Ply           uint16
	Hash          uint64
	KingSq        [2]Square

	Threats  Bitboard
	Checkers Bitboard
	Pinned   [2]Bitboard
}

func (pos *Position) Us() Bitboard {
	return pos.Colors[pos.SideToMove]
}

func (pos *Position) Them() Bitboard {
	return pos.Colors[pos.SideToMove.Opponent()]
}

func (pos *Position) Our(pt PieceType) Bitboard {
	return pos.Pieces[pt] & pos.Us()
}

func (pos *Position) Their(pt PieceType) Bitboard {
	return pos.Pieces[pt] & pos.Them()
}

func (pos *Position) Pawns(c Color) Bitboard {
	return pos.Pieces[Pawn] & pos.Colors[c]
}

func (pos *Position) Knights(c Color) Bitboard {
	return pos.Pieces[Knight] & pos.Colors[c]
}

func (pos *Position) Bishops(c Color) Bitboard {
	return pos.Pieces[Bishop] & pos.Colors[c]
}

func (pos *Position) Rooks(c Color) Bitboard {
	return pos.Pieces[Rook] & pos.Colors[c]
}

func (pos *Position) Queens(c Color) Bitboard {
	return pos.Pieces[Queen] & pos.Colors[c]
}

func (pos *Position) Occ() Bitboard {
	return pos.Colors[White] | pos.Colors[Black]
}

func (pos *Position) PieceBB(piece Piece) Bitboard {
	return pos.Pieces[piece.Type()] & pos.Colors[piece.Color()]
}

func (pos *Position) PlacePiece(piece Piece, sq Square) {
	pos.Pieces[piece.Type()].SetBit(sq)
	pos.Colors[piece.Color()].SetBit(sq)
	pos.Board[sq] = piece
}

func (pos *Position) RemovePiece(sq Square) {
	piece, ok := pos.PieceOn(sq)
	if ok {
		pos.Pieces[piece.Type()].ClearBit(sq)
		pos.Colors[piece.Color()].ClearBit(sq)
		pos.Board[sq] = NoPiece
	}
}

func (pos *Position) PieceOn(sq Square) (Piece, bool) {
	piece := pos.Board[sq]
	return piece, piece != NoPiece
}

func (pos *Position) HasNonPawns() bool {
	return pos.Us()&^(pos.Pieces[King]|pos.Pieces[Pawn]) != 0
}

func newBoard() [64]Piece {
	var board [64]Piece
	for sq := range board {
		board[sq] = NoPiece
	}
	return board
}

func NewPosition() Position {
	pos := Position{
		SideToMove:    White,
		Castling:      noCastling,
		EnPassant:     NoSquare,
		HalfmoveClock: 0,
		Ply:           0,
		Hash:          69420,
		KingSq:        [2]Square{NoSquare, NoSquare},
		Chess960:      IsChess960(),
	}
	pos.Board = newBoard()
	return pos
}
