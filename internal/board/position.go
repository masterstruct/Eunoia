package board

type Bitboards struct {
	Pieces [6]Bitboard
	Colors [2]Bitboard
}

type Position struct {
	Bitboards
	Board          [64]Piece
	SideToMove     Color
	CastlingRights CastlingRights
	CastlingRookSq CastlingRookSquares // TODO: move out of Position, breaks equality checks
	EnPassant      Square
	HalfmoveClock  uint8
	Ply            uint16
	Hash           uint64
	KingSq         [2]Square
}

func (pos *Position) PieceBB(piece Piece) Bitboard {
	return pos.Pieces[piece.Type] & pos.Colors[piece.Color]
}

func (pos *Position) Occupied() Bitboard {
	return pos.Colors[White] | pos.Colors[Black]
}

func (pos *Position) PlacePiece(piece Piece, sq Square) {
	pos.Pieces[piece.Type].SetBit(sq)
	pos.Colors[piece.Color].SetBit(sq)
	pos.Board[sq] = piece
}

func (pos *Position) RemovePiece(sq Square) {
	piece, ok := pos.PieceOn(sq)
	if ok {
		pos.Pieces[piece.Type].ClearBit(sq)
		pos.Colors[piece.Color].ClearBit(sq)
		pos.Board[sq] = NoPiece
	}
}

func (pos *Position) PieceOn(sq Square) (Piece, bool) {
	piece := pos.Board[sq]
	return piece, piece != NoPiece
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
		SideToMove:     White,
		CastlingRights: AllCastling,
		EnPassant:      NoSquare,
		HalfmoveClock:  0,
		Ply:            0,
		Hash:           69420,
	}
	pos.Board = newBoard()
	return pos
}
