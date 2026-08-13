package movegen

import (
	"github.com/masterstruct/Eunoia/internal/board"
)

func IsSquareAttacked(pos board.Position, sq board.Square, byColor board.Color) bool {
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

func InCheck(pos board.Position) bool {
	// TODO: instead of creating new piece for
	// pieceBB lookup use stored king square values
	return IsSquareAttacked(
		pos,
		pos.PieceBB(board.NewPiece(board.King, pos.SideToMove)).LSB(),
		pos.SideToMove.Opponent(),
	)
}

func genKnightMoves(pos board.Position) []board.Move {
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

func genBishopMoves(pos board.Position) []board.Move {
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

func genRookMoves(pos board.Position) []board.Move {
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

func genQueenMoves(pos board.Position) []board.Move {
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

func genPawnMoves(pos board.Position) []board.Move {
	// TODO: replace slices. there can be multiple pieces so reallocation happens
	// a pawn can make up to 8 moves
	movelist := make([]board.Move, 0, 8)
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
		attackBB := PawnAttacks[color][from]
		fromRank := from.Rank()
		for to := range (attackBB & pos.Occupied()).Bits() {
			toPiece := pos.Board[to]
			if toPiece.Color != color &&
				((color == board.White && fromRank == board.Rank7) ||
					(color == board.Black && fromRank == board.Rank2)) {
				// capture promo
				movelist = append(movelist, board.NewCapturePromo(from, to, board.Knight))
				movelist = append(movelist, board.NewCapturePromo(from, to, board.Bishop))
				movelist = append(movelist, board.NewCapturePromo(from, to, board.Rook))
				movelist = append(movelist, board.NewCapturePromo(from, to, board.Queen))
				continue
			}
			if toPiece.Color != color {
				movelist = append(movelist, board.NewCapture(from, to))
				continue
			}
		}

		// regular pushes
		if color == board.White {
			if !pos.Occupied().IsBitSet(from + 8) {
				if fromRank == board.Rank7 {
					movelist = append(movelist, board.NewPromo(from, from+8, board.Knight))
					movelist = append(movelist, board.NewPromo(from, from+8, board.Bishop))
					movelist = append(movelist, board.NewPromo(from, from+8, board.Rook))
					movelist = append(movelist, board.NewPromo(from, from+8, board.Queen))
					continue
				}
				movelist = append(movelist, board.NewMove(from, from+8))
				if fromRank == board.Rank2 && !pos.Occupied().IsBitSet(from+16) {
					movelist = append(movelist, board.NewMove(from, from+16))
				}
			}
		} else {
			if !pos.Occupied().IsBitSet(from - 8) {
				if fromRank == board.Rank2 {
					movelist = append(movelist, board.NewPromo(from, from-8, board.Knight))
					movelist = append(movelist, board.NewPromo(from, from-8, board.Bishop))
					movelist = append(movelist, board.NewPromo(from, from-8, board.Rook))
					movelist = append(movelist, board.NewPromo(from, from-8, board.Queen))
					continue
				}
				movelist = append(movelist, board.NewMove(from, from-8))
				if fromRank == board.Rank7 && !pos.Occupied().IsBitSet(from-16) {
					movelist = append(movelist, board.NewMove(from, from-16))
				}
			}
		}
	}
	return movelist
}

// func genKingMoves(pos board.Position) []board.Move {
// 	// a king can make up to 8 moves
// 	movelist := make([]board.Move, 0, 8)
// 	color := pos.SideToMove

// 	// loop over each king
// 	// TODO: replace with king square lookup
// 	kingBB := pos.PieceBB(board.NewPiece(board.King, color))
// 	from := kingBB.LSB()

// 	attackBB := KingAttacks[from]

// 	for to := range attackBB.Bits() {
// 		toPiece := pos.Board[to]

// 		if toPiece == board.NoPiece {
// 			// quiet move
// 			movelist = append(movelist, board.NewMove(from, to))
// 			continue
// 		}
// 		// capture
// 		if toPiece.Color != color {
// 			movelist = append(movelist, board.NewCapture(from, to))
// 		}
// 	}

// 	return movelist
// }
