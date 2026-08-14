package movegen

import (
	"github.com/masterstruct/Eunoia/internal/board"
)

func IsSquareAttacked(pos *board.Position, sq board.Square, byColor board.Color) bool {
	if byColor == board.NoColor {
		return false
	}
	if PawnAttacks[byColor.Opponent()][sq]&pos.PieceBB(board.NewPiece(board.Pawn, byColor)) != 0 {
		return true
	}
	if KnightAttacks[sq]&pos.PieceBB(board.NewPiece(board.Knight, byColor)) != 0 {
		return true
	}
	if KingAttacks[sq]&pos.PieceBB(board.NewPiece(board.King, byColor)) != 0 {
		return true
	}

	occupied := pos.Occupied()
	rooks := pos.PieceBB(board.NewPiece(board.Rook, byColor))
	queens := pos.PieceBB(board.NewPiece(board.Queen, byColor))

	rookAttacks := RookAttacks(sq, occupied)
	if rookAttacks&(rooks|queens) != 0 {
		return true
	}

	bishops := pos.PieceBB(board.NewPiece(board.Bishop, byColor))

	bishopAttacks := BishopAttacks(sq, occupied)
	if bishopAttacks&(bishops|queens) != 0 {
		return true
	}
	return false
}

func InCheck(pos *board.Position, color board.Color) bool {
	// TODO: instead of creating new piece for
	// pieceBB lookup use stored king square values
	return IsSquareAttacked(
		pos,
		pos.PieceBB(board.NewPiece(board.King, color)).LSB(),
		color.Opponent(),
	)
}

func genKnightMoves(pos *board.Position) []board.Move {
	// a knight can make up to 8 moves
	movelist := make([]board.Move, 0, 8)
	color := pos.SideToMove

	// loop over each knight
	knightBB := pos.PieceBB(board.NewPiece(board.Knight, color))
	for from := range knightBB.Bits() {
		attackBB := KnightAttacks[from]

		for to := range attackBB.Bits() {
			toPiece := pos.Board[to]

			if toPiece == board.NoPiece {
				// quiet move
				movelist = append(movelist, board.NewMove(from, to))
				continue
			}
			// capture
			if toPiece.Color != color {
				movelist = append(movelist, board.NewCapture(from, to))
			}
		}
	}
	return movelist
}

func genBishopMoves(pos *board.Position) []board.Move {
	// a bishop can make up to 13 moves
	movelist := make([]board.Move, 0, 13)
	color := pos.SideToMove
	occupied := pos.Occupied()

	// loop over each bishop
	bishopBB := pos.PieceBB(board.NewPiece(board.Bishop, color))
	for from := range bishopBB.Bits() {
		attackBB := BishopAttacks(from, occupied)

		for to := range (attackBB & occupied).Bits() {
			// capture
			toPiece := pos.Board[to]
			if toPiece.Color != color {
				movelist = append(movelist, board.NewCapture(from, to))
			}
		}

		attackBB &^= occupied
		for to := range attackBB.Bits() {
			// quiet move
			movelist = append(movelist, board.NewMove(from, to))
		}
	}
	return movelist
}

func genRookMoves(pos *board.Position) []board.Move {
	// a rook can make up to 14 moves
	movelist := make([]board.Move, 0, 14)
	color := pos.SideToMove
	occupied := pos.Occupied()

	// loop over each rook
	rookBB := pos.PieceBB(board.NewPiece(board.Rook, color))
	for from := range rookBB.Bits() {
		attackBB := RookAttacks(from, occupied)

		for to := range (attackBB & occupied).Bits() {
			// capture
			toPiece := pos.Board[to]
			if toPiece.Color != color {
				movelist = append(movelist, board.NewCapture(from, to))
			}
		}

		attackBB &^= occupied
		for to := range attackBB.Bits() {
			// quiet move
			movelist = append(movelist, board.NewMove(from, to))
		}
	}
	return movelist
}

func genQueenMoves(pos *board.Position) []board.Move {
	// a queen can make up to 27 moves
	movelist := make([]board.Move, 0, 27)
	color := pos.SideToMove
	occupied := pos.Occupied()

	// loop over each queen
	queenBB := pos.PieceBB(board.NewPiece(board.Queen, color))
	for from := range queenBB.Bits() {
		attackBB := QueenAttacks(from, occupied)

		for to := range (attackBB & occupied).Bits() {
			// capture
			toPiece := pos.Board[to]
			if toPiece.Color != color {
				movelist = append(movelist, board.NewCapture(from, to))
			}
		}

		attackBB &^= occupied
		for to := range attackBB.Bits() {
			// quiet move
			movelist = append(movelist, board.NewMove(from, to))
		}
	}
	return movelist
}

func genPawnMoves(pos *board.Position) []board.Move {
	// TODO: replace slices. there can be multiple pieces so reallocation happens
	// a pawn can make up to 12 moves
	movelist := make([]board.Move, 0, 12)
	color := pos.SideToMove

	pawnBB := pos.PieceBB(board.NewPiece(board.Pawn, color))

	epSq := pos.EnPassant
	if epSq != board.NoSquare {
		attackers := PawnAttacks[color.Opponent()][epSq] & pawnBB
		for from := range attackers.Bits() {
			movelist = append(movelist, board.NewEnPassant(from, epSq))
		}
	}

	// loop over each pawn
	for from := range pawnBB.Bits() {
		fromRank := from.Rank()
		for to := range (PawnAttacks[color][from] & pos.Occupied()).Bits() {
			toPiece := pos.Board[to]
			if toPiece.Color == color {
				continue
			}
			if (color == board.White && fromRank == board.Rank7) ||
				(color == board.Black && fromRank == board.Rank2) {
				// capture promo
				movelist = append(movelist, board.NewCapturePromo(from, to, board.Knight))
				movelist = append(movelist, board.NewCapturePromo(from, to, board.Bishop))
				movelist = append(movelist, board.NewCapturePromo(from, to, board.Rook))
				movelist = append(movelist, board.NewCapturePromo(from, to, board.Queen))
				continue
			}
			movelist = append(movelist, board.NewCapture(from, to))
		}

		// regular pushes
		// TODO: bit shifts
		if color == board.White {
			up := from.Up()
			if !pos.Occupied().IsBitSet(up) {
				dup := up.Up()
				switch fromRank {
				case board.Rank7:
					movelist = append(movelist, board.NewPromo(from, up, board.Knight))
					movelist = append(movelist, board.NewPromo(from, up, board.Bishop))
					movelist = append(movelist, board.NewPromo(from, up, board.Rook))
					movelist = append(movelist, board.NewPromo(from, up, board.Queen))
					continue
				case board.Rank2:
					if !pos.Occupied().IsBitSet(dup) {
						movelist = append(movelist, board.NewDoublePush(from, dup))
					}
				}
				movelist = append(movelist, board.NewMove(from, up))
			}
		} else {
			down := from.Down()
			if !pos.Occupied().IsBitSet(down) {
				ddown := down.Down()
				switch fromRank {
				case board.Rank2:
					movelist = append(movelist, board.NewPromo(from, down, board.Knight))
					movelist = append(movelist, board.NewPromo(from, down, board.Bishop))
					movelist = append(movelist, board.NewPromo(from, down, board.Rook))
					movelist = append(movelist, board.NewPromo(from, down, board.Queen))
					continue
				case board.Rank7:
					if !pos.Occupied().IsBitSet(ddown) {
						movelist = append(movelist, board.NewDoublePush(from, ddown))
					}
				}
				movelist = append(movelist, board.NewMove(from, down))
			}
		}
	}
	return movelist
}

func genKingMoves(pos *board.Position) []board.Move {
	// a king can make up to 8 moves
	movelist := make([]board.Move, 0, 8)
	color := pos.SideToMove

	// loop over each king
	// TODO: replace with king square lookup
	kingBB := pos.PieceBB(board.NewPiece(board.King, color))
	from := kingBB.LSB()

	attackBB := KingAttacks[from]

	for to := range attackBB.Bits() {
		toPiece := pos.Board[to]

		if toPiece == board.NoPiece {
			// quiet move
			movelist = append(movelist, board.NewMove(from, to))
			continue
		}
		// capture
		if toPiece.Color != color {
			movelist = append(movelist, board.NewCapture(from, to))
		}
	}

	castlingRights := pos.CastlingRights
	if castlingRights != board.NoCastling {
		rooks := pos.CastlingRooks()
		if color == board.White {
			if castlingRights.Has(board.WhiteKingside) && canCastle(pos, from, rooks.WhiteKingside) {
				movelist = append(movelist, board.NewCastle(from, rooks.WhiteKingside))
			}
			if castlingRights.Has(board.WhiteQueenside) && canCastle(pos, from, rooks.WhiteQueenside) {
				movelist = append(movelist, board.NewCastle(from, rooks.WhiteQueenside))
			}
		} else {
			if castlingRights.Has(board.BlackKingside) && canCastle(pos, from, rooks.BlackKingside) {
				movelist = append(movelist, board.NewCastle(from, rooks.BlackKingside))
			}
			if castlingRights.Has(board.BlackQueenside) && canCastle(pos, from, rooks.BlackQueenside) {
				movelist = append(movelist, board.NewCastle(from, rooks.BlackQueenside))
			}
		}
	}

	return movelist
}

func canCastle(pos *board.Position, kingSq, rookSq board.Square) bool {
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
