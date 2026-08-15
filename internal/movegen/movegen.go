package movegen

import (
	"github.com/masterstruct/Eunoia/internal/board"
)

func IsSquareAttacked(pos *board.Position, sq board.Square, byColor board.Color) bool {
	if byColor == board.NoColor {
		return false
	}

	if PawnAttacks[byColor.Opponent()][sq]&pos.PieceBB(board.Piece{Type: board.Pawn, Color: byColor}) != 0 {
		return true
	}
	if KnightAttacks[sq]&pos.PieceBB(board.Piece{Type: board.Knight, Color: byColor}) != 0 {
		return true
	}
	if KingAttacks[sq]&pos.PieceBB(board.Piece{Type: board.King, Color: byColor}) != 0 {
		return true
	}

	occupied := pos.Occupied()
	rooks := pos.PieceBB(board.Piece{Type: board.Rook, Color: byColor})
	queens := pos.PieceBB(board.Piece{Type: board.Queen, Color: byColor})

	rookAttacks := RookAttacks(sq, occupied)
	if rookAttacks&(rooks|queens) != 0 {
		return true
	}

	bishops := pos.PieceBB(board.Piece{Type: board.Bishop, Color: byColor})

	bishopAttacks := BishopAttacks(sq, occupied)
	if bishopAttacks&(bishops|queens) != 0 {
		return true
	}
	return false
}

func InCheck(pos *board.Position, color board.Color) bool {
	if color == board.Black {
		return IsSquareAttacked(
			pos,
			pos.BlackKing,
			color.Opponent(),
		)
	}
	return IsSquareAttacked(
		pos,
		pos.WhiteKing,
		color.Opponent(),
	)
}

func genKnightMoves(pos *board.Position) []board.Move {
	// a knight can make up to 8 moves
	movelist := make([]board.Move, 0, 16)
	color := pos.SideToMove
	opponents := pos.Colors[color.Opponent()]
	occupied := pos.Occupied()

	// loop over each knight
	knights := pos.PieceBB(board.Piece{Type: board.Knight, Color: color})
	for knights != 0 {
		from := knights.PopLSB()
		attacks := KnightAttacks[from]
		captures := attacks & opponents
		quiets := attacks &^ occupied

		for captures != 0 {
			// capture
			movelist = append(movelist, board.NewCapture(from, captures.PopLSB()))
		}

		for quiets != 0 {
			// quiet move
			movelist = append(movelist, board.NewMove(from, quiets.PopLSB()))
		}
	}
	return movelist
}

func genBishopMoves(pos *board.Position) []board.Move {
	// a bishop can make up to 13 moves
	movelist := make([]board.Move, 0, 13)
	color := pos.SideToMove
	opponents := pos.Colors[color.Opponent()]
	occupied := pos.Occupied()

	// loop over each bishop
	bishops := pos.PieceBB(board.Piece{Type: board.Bishop, Color: color})
	for bishops != 0 {
		from := bishops.PopLSB()
		attacks := BishopAttacks(from, occupied)
		captures := attacks & opponents
		quiets := attacks &^ occupied

		for captures != 0 {
			// capture
			movelist = append(movelist, board.NewCapture(from, captures.PopLSB()))
		}

		for quiets != 0 {
			// quiet move
			movelist = append(movelist, board.NewMove(from, quiets.PopLSB()))
		}
	}
	return movelist
}

func genRookMoves(pos *board.Position) []board.Move {
	// a rook can make up to 14 moves
	movelist := make([]board.Move, 0, 14)
	color := pos.SideToMove
	opponents := pos.Colors[color.Opponent()]
	occupied := pos.Occupied()

	// loop over each rook
	rooks := pos.PieceBB(board.Piece{Type: board.Rook, Color: color})
	for rooks != 0 {
		from := rooks.PopLSB()
		attacks := RookAttacks(from, occupied)
		captures := attacks & opponents
		quiets := attacks &^ occupied

		for captures != 0 {
			// capture
			movelist = append(movelist, board.NewCapture(from, captures.PopLSB()))
		}

		for quiets != 0 {
			// quiet move
			movelist = append(movelist, board.NewMove(from, quiets.PopLSB()))
		}
	}
	return movelist
}

func genQueenMoves(pos *board.Position) []board.Move {
	// a queen can make up to 27 moves
	movelist := make([]board.Move, 0, 27)
	color := pos.SideToMove
	opponents := pos.Colors[color.Opponent()]
	occupied := pos.Occupied()

	// loop over each queen
	queens := pos.PieceBB(board.Piece{Type: board.Queen, Color: color})
	for queens != 0 {
		from := queens.PopLSB()
		attacks := QueenAttacks(from, occupied)
		captures := attacks & opponents
		quiets := attacks &^ occupied

		for captures != 0 {
			// capture
			movelist = append(movelist, board.NewCapture(from, captures.PopLSB()))
		}

		for quiets != 0 {
			// quiet move
			movelist = append(movelist, board.NewMove(from, quiets.PopLSB()))
		}
	}
	return movelist
}

func genPawnMoves(pos *board.Position) []board.Move {
	movelist := make([]board.Move, 0, 16)

	color := pos.SideToMove
	enemyColor := color.Opponent()
	occupied := pos.Occupied()
	enemyOcc := pos.Colors[enemyColor]

	pawns := pos.PieceBB(board.Piece{Type: board.Pawn, Color: color})

	// en passant
	epSq := pos.EnPassant
	if epSq != board.NoSquare {
		attackers := PawnAttacks[enemyColor][epSq] & pawns
		for attackers != 0 {
			movelist = append(movelist, board.NewEnPassant(attackers.PopLSB(), epSq))
		}
	}

	for pawns != 0 {
		from := pawns.PopLSB()

		// captures
		captures := PawnAttacks[color][from] & enemyOcc
		for captures != 0 {
			to := captures.PopLSB()
			if to.Rank() == board.Rank1 || to.Rank() == board.Rank8 {
				movelist = append(movelist, board.NewCapturePromo(from, to, board.Knight))
				movelist = append(movelist, board.NewCapturePromo(from, to, board.Bishop))
				movelist = append(movelist, board.NewCapturePromo(from, to, board.Rook))
				movelist = append(movelist, board.NewCapturePromo(from, to, board.Queen))
			} else {
				movelist = append(movelist, board.NewCapture(from, to))
			}
		}

		// pushes
		if color == board.White {
			to := from + 8
			if occupied.IsBitSet(to) {
				continue
			}

			if from.Rank() == board.Rank7 {
				movelist = append(movelist, board.NewPromo(from, to, board.Knight))
				movelist = append(movelist, board.NewPromo(from, to, board.Bishop))
				movelist = append(movelist, board.NewPromo(from, to, board.Rook))
				movelist = append(movelist, board.NewPromo(from, to, board.Queen))
				continue
			}

			movelist = append(movelist, board.NewMove(from, to))

			if from.Rank() == board.Rank2 {
				to2 := from + 16
				if !occupied.IsBitSet(to2) {
					movelist = append(movelist, board.NewDoublePush(from, to2))
				}
			}
		} else {
			to := from - 8
			if occupied.IsBitSet(to) {
				continue
			}

			if from.Rank() == board.Rank2 {
				movelist = append(movelist, board.NewPromo(from, to, board.Knight))
				movelist = append(movelist, board.NewPromo(from, to, board.Bishop))
				movelist = append(movelist, board.NewPromo(from, to, board.Rook))
				movelist = append(movelist, board.NewPromo(from, to, board.Queen))
				continue
			}

			movelist = append(movelist, board.NewMove(from, to))

			if from.Rank() == board.Rank7 {
				to2 := from - 16
				if !occupied.IsBitSet(to2) {
					movelist = append(movelist, board.NewDoublePush(from, to2))
				}
			}
		}
	}

	return movelist
}

func genKingMoves(pos *board.Position) []board.Move {
	// a king can make up to 8 moves
	movelist := make([]board.Move, 0, 8)
	color := pos.SideToMove

	var from board.Square
	castlingRights := pos.CastlingRights
	if color == board.Black {
		from = pos.BlackKing
		if castlingRights != board.NoCastling {
			rooks := board.CastlingRooks
			if castlingRights.Has(board.BlackKingside) && canCastle(pos, from, rooks.BlackKingside) {
				movelist = append(movelist, board.NewCastle(from, rooks.BlackKingside))
			}
			if castlingRights.Has(board.BlackQueenside) && canCastle(pos, from, rooks.BlackQueenside) {
				movelist = append(movelist, board.NewCastle(from, rooks.BlackQueenside))
			}

		}
	} else {
		from = pos.WhiteKing
		if castlingRights != board.NoCastling {
			rooks := board.CastlingRooks
			if castlingRights.Has(board.WhiteKingside) && canCastle(pos, from, rooks.WhiteKingside) {
				movelist = append(movelist, board.NewCastle(from, rooks.WhiteKingside))
			}
			if castlingRights.Has(board.WhiteQueenside) && canCastle(pos, from, rooks.WhiteQueenside) {
				movelist = append(movelist, board.NewCastle(from, rooks.WhiteQueenside))
			}
		}
	}

	occupied := pos.Occupied()

	attacks := KingAttacks[from]
	captures := attacks & pos.Colors[color.Opponent()]
	quiets := attacks &^ occupied

	for captures != 0 {
		// capture
		to := captures.PopLSB()
		if pos.Board[to].Color != color {
			movelist = append(movelist, board.NewCapture(from, to))
		}
	}

	for quiets != 0 {
		// quiet move
		to := quiets.PopLSB()
		movelist = append(movelist, board.NewMove(from, to))
	}

	return movelist
}

func canCastle(pos *board.Position, kingSq, rookSq board.Square) bool {
	// TODO: on ParseFEN(), precompute which bits to check
	// with occupied and IsSquareAttacked.
	// Compare (rookPath|kingPath)&occupied==0,
	// then just check if kingPath is attacked.

	rank := kingSq.Rank()
	kingFile := kingSq.File()
	rookFile := rookSq.File()
	occupied := pos.Occupied()
	occupied.ClearBit(kingSq)
	occupied.ClearBit(rookSq) // in chess960 king can jump over rook

	var toFileKing board.File
	var toFileRook board.File
	var kingDir board.File
	var rookDir board.File
	if rookSq > kingSq {
		// kingside
		kingDir = 1
		rookDir = -1
		toFileKing = board.FileG
		toFileRook = board.FileF
		if toFileRook > rookFile {
			// castling kingside (rook moves left), but because chess960,
			// rook can start on H1 and actually move left.
			rookDir = 1
		}
	} else {
		// queenside
		kingDir = -1
		rookDir = -1
		toFileKing = board.FileC
		toFileRook = board.FileD
		if toFileKing > kingFile {
			// castling queenside (left), but because chess960,
			// king can start on B1 and actually move right.
			// this cannot happen with kingside if position is legal.
			kingDir = 1
		}
		if toFileRook > rookFile {
			rookDir = 1
		}
	}

	color := pos.Board[kingSq].Color
	for file := kingFile; file != toFileKing+kingDir; file += kingDir {
		sq := board.NewSquare(file, rank)
		if occupied.IsBitSet(sq) || IsSquareAttacked(pos, sq, color.Opponent()) {
			return false
		}
	}
	for file := rookFile; file != toFileRook+rookDir; file += rookDir {
		sq := board.NewSquare(file, rank)
		if occupied.IsBitSet(sq) {
			return false
		}
	}
	return true
}

func GeneratePseudolegalMoves(pos *board.Position) []board.Move {
	// a chess position can have up to 218 legal moves
	movelist := make([]board.Move, 0, 218)

	movelist = append(movelist, genKnightMoves(pos)...)
	movelist = append(movelist, genBishopMoves(pos)...)
	movelist = append(movelist, genRookMoves(pos)...)
	movelist = append(movelist, genQueenMoves(pos)...)
	movelist = append(movelist, genPawnMoves(pos)...)
	movelist = append(movelist, genKingMoves(pos)...)

	return movelist
}
