package board

func (pos *Position) MakeMove(move Move) Position {
	// To optimize this further, you can add 2 helper
	// functions: MakeBlackMove and MakeWhiteMove.
	// currently color checks add 3 code branches.

	newPos := *pos
	from, to := move.From(), move.To()
	piece, _ := newPos.PieceOn(from)
	pieceType := piece.Type
	capturedPiece, _ := newPos.PieceOn(to)
	color := newPos.SideToMove

	isPromo := move.IsPromo()
	isEnPassant := move.IsEnPassant()
	isCapture := move.IsCapture()
	isCastle := move.IsCastle()

	dir := 16*int(color) - 8
	epSquare := Square(int(to) - dir)

	newPos.HalfmoveClock++
	newPos.EnPassant = NoSquare

	newPos.RemovePiece(from)
	if !isCastle {
		// PlacePiece(sq) overrides the square, therefore if
		// I've done everything correctly, this is unnecessary:
		newPos.RemovePiece(to)

		if isPromo {
			newPos.PlacePiece(NewPiece(move.Promo(), color), to)
		} else {
			newPos.PlacePiece(piece, to)
		}

		if isEnPassant {
			newPos.RemovePiece(epSquare)
		}

		// rook move - remove castling rights
		if pieceType == Rook {
			rooks := pos.CastlingRookSq
			switch from {
			case rooks.WhiteQueenside:
				newPos.CastlingRights.Remove(WhiteQueenside)
			case rooks.WhiteKingside:
				newPos.CastlingRights.Remove(WhiteKingside)
			case rooks.BlackQueenside:
				newPos.CastlingRights.Remove(BlackQueenside)
			case rooks.BlackKingside:
				newPos.CastlingRights.Remove(BlackKingside)
			}
		}

		// rook captured - remove castling rights
		if isCapture && capturedPiece.Type == Rook {
			rooks := pos.CastlingRookSq
			switch to {
			case rooks.WhiteQueenside:
				newPos.CastlingRights.Remove(WhiteQueenside)
			case rooks.WhiteKingside:
				newPos.CastlingRights.Remove(WhiteKingside)
			case rooks.BlackQueenside:
				newPos.CastlingRights.Remove(BlackQueenside)
			case rooks.BlackKingside:
				newPos.CastlingRights.Remove(BlackKingside)
			}
		}

		if isCapture || pieceType == Pawn {
			newPos.HalfmoveClock = 0
		}
		if move.IsDoublePush() {
			newPos.EnPassant = epSquare
		}

		// king moved - remove castling rights
		if pieceType == King {
			newPos.KingSq[color] = to
			if color == Black {
				newPos.CastlingRights.Remove(BlackKingside | BlackQueenside)
			} else {
				newPos.CastlingRights.Remove(WhiteKingside | WhiteQueenside)
			}
		}
	} else {
		// castle - move pieces
		rooks := pos.CastlingRookSq

		if color == Black {
			if move.IsKingsideCastle() {
				newPos.RemovePiece(rooks.BlackKingside)
				newPos.PlacePiece(BlackKing, G8)
				newPos.PlacePiece(BlackRook, F8)
				newPos.KingSq[Black] = G8
			} else {
				newPos.RemovePiece(rooks.BlackQueenside)
				newPos.PlacePiece(BlackKing, C8)
				newPos.PlacePiece(BlackRook, D8)
				newPos.KingSq[Black] = C8
			}
			newPos.CastlingRights.Remove(BlackKingside | BlackQueenside)
		} else {
			if move.IsKingsideCastle() {
				newPos.RemovePiece(rooks.WhiteKingside)
				newPos.PlacePiece(WhiteKing, G1)
				newPos.PlacePiece(WhiteRook, F1)
				newPos.KingSq[White] = G1
			} else {
				newPos.RemovePiece(rooks.WhiteQueenside)
				newPos.PlacePiece(WhiteKing, C1)
				newPos.PlacePiece(WhiteRook, D1)
				newPos.KingSq[White] = C1
			}
			newPos.CastlingRights.Remove(WhiteKingside | WhiteQueenside)
		}
	}

	newPos.SideToMove = color.Opponent()
	newPos.Ply++

	return newPos
}
