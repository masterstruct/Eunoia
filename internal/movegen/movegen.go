package movegen

import (
	"github.com/masterstruct/Eunoia/internal/board"
)

const MaxMoves = 256

type Movelist struct {
	Moves [MaxMoves]board.Move
	Len   int
}

func (ml *Movelist) Add(m board.Move) {
	ml.Moves[ml.Len] = m
	ml.Len++
}

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
	return IsSquareAttacked(
		pos,
		pos.KingSq[color],
		color.Opponent(),
	)
}

func genKnightMoves(pos *board.Position, movelist *Movelist) {
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
			movelist.Add(board.NewCapture(from, captures.PopLSB()))
		}

		for quiets != 0 {
			// quiet move
			movelist.Add(board.NewMove(from, quiets.PopLSB()))
		}
	}
}

func genBishopMoves(pos *board.Position, movelist *Movelist) {
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
			movelist.Add(board.NewCapture(from, captures.PopLSB()))
		}

		for quiets != 0 {
			// quiet move
			movelist.Add(board.NewMove(from, quiets.PopLSB()))
		}
	}
}

func genRookMoves(pos *board.Position, movelist *Movelist) {
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
			movelist.Add(board.NewCapture(from, captures.PopLSB()))
		}

		for quiets != 0 {
			// quiet move
			movelist.Add(board.NewMove(from, quiets.PopLSB()))
		}
	}
}

func genQueenMoves(pos *board.Position, movelist *Movelist) {
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
			movelist.Add(board.NewCapture(from, captures.PopLSB()))
		}

		for quiets != 0 {
			// quiet move
			movelist.Add(board.NewMove(from, quiets.PopLSB()))
		}
	}
}

func genPawnMoves(pos *board.Position, movelist *Movelist) {
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
			movelist.Add(board.NewEnPassant(attackers.PopLSB(), epSq))
		}
	}

	for pawns != 0 {
		from := pawns.PopLSB()

		// captures
		captures := PawnAttacks[color][from] & enemyOcc
		for captures != 0 {
			to := captures.PopLSB()
			if to.Rank() == board.Rank1 || to.Rank() == board.Rank8 {
				movelist.Add(board.NewCapturePromo(from, to, board.Knight))
				movelist.Add(board.NewCapturePromo(from, to, board.Bishop))
				movelist.Add(board.NewCapturePromo(from, to, board.Rook))
				movelist.Add(board.NewCapturePromo(from, to, board.Queen))
			} else {
				movelist.Add(board.NewCapture(from, to))
			}
		}

		// pushes
		if color == board.White {
			to := from + 8
			if occupied.IsBitSet(to) {
				continue
			}

			if to.Rank() == board.Rank8 {
				movelist.Add(board.NewPromo(from, to, board.Knight))
				movelist.Add(board.NewPromo(from, to, board.Bishop))
				movelist.Add(board.NewPromo(from, to, board.Rook))
				movelist.Add(board.NewPromo(from, to, board.Queen))
				continue
			}

			movelist.Add(board.NewMove(from, to))

			if from.Rank() == board.Rank2 {
				to2 := from + 16
				if !occupied.IsBitSet(to2) {
					movelist.Add(board.NewDoublePush(from, to2))
				}
			}
		} else {
			to := from - 8
			if occupied.IsBitSet(to) {
				continue
			}

			if to.Rank() == board.Rank1 {
				movelist.Add(board.NewPromo(from, to, board.Knight))
				movelist.Add(board.NewPromo(from, to, board.Bishop))
				movelist.Add(board.NewPromo(from, to, board.Rook))
				movelist.Add(board.NewPromo(from, to, board.Queen))
				continue
			}

			movelist.Add(board.NewMove(from, to))

			if from.Rank() == board.Rank7 {
				to2 := from - 16
				if !occupied.IsBitSet(to2) {
					movelist.Add(board.NewDoublePush(from, to2))
				}
			}
		}
	}
}

func genKingMoves(pos *board.Position, movelist *Movelist) {
	color := pos.SideToMove

	castlingRights := pos.CastlingRights
	from := pos.KingSq[color]
	if castlingRights != board.NoCastling {
		rooks := pos.CastlingRookSq
		if color == board.Black {
			if castlingRights.Has(board.BlackKingside) && canCastle(pos, from, rooks.BlackKingside) {
				movelist.Add(board.NewCastle(from, rooks.BlackKingside))
			}
			if castlingRights.Has(board.BlackQueenside) && canCastle(pos, from, rooks.BlackQueenside) {
				movelist.Add(board.NewCastle(from, rooks.BlackQueenside))
			}
		} else {
			if castlingRights.Has(board.WhiteKingside) && canCastle(pos, from, rooks.WhiteKingside) {
				movelist.Add(board.NewCastle(from, rooks.WhiteKingside))
			}
			if castlingRights.Has(board.WhiteQueenside) && canCastle(pos, from, rooks.WhiteQueenside) {
				movelist.Add(board.NewCastle(from, rooks.WhiteQueenside))
			}
		}
	}

	attacks := KingAttacks[from]
	captures := attacks & pos.Colors[color.Opponent()]
	quiets := attacks &^ pos.Occupied()

	for captures != 0 {
		// capture
		movelist.Add(board.NewCapture(from, captures.PopLSB()))
	}

	for quiets != 0 {
		// quiet move
		movelist.Add(board.NewMove(from, quiets.PopLSB()))
	}
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

func GeneratePseudolegalMoves(pos *board.Position, movelist *Movelist) {
	genKnightMoves(pos, movelist)
	genBishopMoves(pos, movelist)
	genRookMoves(pos, movelist)
	genQueenMoves(pos, movelist)
	genPawnMoves(pos, movelist)
	genKingMoves(pos, movelist)
}
