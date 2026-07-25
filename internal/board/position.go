package board

type CastlingRights uint8

const (
	BlackKingside CastlingRights = 1 << iota
	BlackQueenside
	WhiteKingside
	WhiteQueenside

	NoCastling  CastlingRights = 0
	AnyCastling CastlingRights = BlackKingside | BlackQueenside | WhiteKingside | WhiteQueenside
)

type Bitboards struct {
	pieces [6]Bitboard
	colors [2]Bitboard
}

type Position struct {
	Bitboards
	Board          [64]Piece
	SideToMove     Color
	CastlingRights CastlingRights
	EnPassant      Square
	HalfmoveClock  uint8
	Ply            uint16
	Hash           uint64
}

func (pos *Position) PieceBB(piece Piece) Bitboard {
	return pos.pieces[piece.Type] & pos.colors[piece.Color]
}

func (pos *Position) Occupied() Bitboard {
	return pos.colors[White] | pos.colors[Black]
}

func (pos *Position) PlacePiece(piece Piece, sq Square) {
	pos.pieces[piece.Type].SetBit(sq)
	pos.colors[piece.Color].SetBit(sq)
	pos.Board[sq] = piece
}

func (pos *Position) RemovePiece(sq Square) {
	piece, ok := pos.PieceOn(sq)
	if ok {
		pos.pieces[piece.Type].ClearBit(sq)
		pos.colors[piece.Color].ClearBit(sq)
		pos.Board[sq] = NoPiece
	}
}

func (pos *Position) PieceOn(sq Square) (Piece, bool) {
	piece := pos.Board[sq]
	return piece, piece != NoPiece
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
