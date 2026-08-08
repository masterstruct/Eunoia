package board

func (pos Position) MakeMove(move Move) Position {
	newPos := pos
	from, to := move.From(), move.To()
	piece, _ := pos.PieceOn(from)
	capturedPiece, _ := pos.PieceOn(to)

	// move/capture/promote piece
	newPos.RemovePiece(from)
	newPos.RemovePiece(to)
	if move.IsPromo() {
		newPos.PlacePiece(NewPiece(move.Promo(), piece.Color), to)
	} else {
		newPos.PlacePiece(piece, to)
	}

	if move.IsEnPassant() {
		if piece == BlackPawn {
			newPos.RemovePiece(to + 8)
		} else {
			newPos.RemovePiece(to - 8)
		}
	}

	// castling rights

	if newPos.CastlingRights != NoCastling && (piece.Type == Rook) {
		switch from {
		case A1:
			newPos.CastlingRights.Remove(WhiteQueenside)
		case H1:
			newPos.CastlingRights.Remove(WhiteKingside)
		case A8:
			newPos.CastlingRights.Remove(BlackQueenside)
		case H8:
			newPos.CastlingRights.Remove(BlackKingside)
		}
	}

	if move.IsCapture() && capturedPiece.Type == Rook {
		switch to {
		case A1:
			newPos.CastlingRights.Remove(WhiteQueenside)
		case H1:
			newPos.CastlingRights.Remove(WhiteKingside)
		case A8:
			newPos.CastlingRights.Remove(BlackQueenside)
		case H8:
			newPos.CastlingRights.Remove(BlackKingside)
		}
	}

	if move.IsCastle() {
		if piece.Color == Black {
			newPos.CastlingRights.Remove(BlackKingside)
			newPos.CastlingRights.Remove(BlackQueenside)
			newPos.PlacePiece(piece, to)
			if move.IsKingsideCastle() {
				newPos.RemovePiece(H8)
				newPos.PlacePiece(BlackRook, F8)
			} else {
				newPos.RemovePiece(A8)
				newPos.PlacePiece(BlackRook, D8)
			}
		} else {
			newPos.CastlingRights.Remove(WhiteKingside)
			newPos.CastlingRights.Remove(WhiteQueenside)
			if move.IsKingsideCastle() {
				newPos.RemovePiece(H1)
				newPos.PlacePiece(WhiteRook, F1)
			} else {
				newPos.RemovePiece(A1)
				newPos.PlacePiece(WhiteRook, D1)
			}
		}
	} else {
		if piece.Type == King {
			newPos.CastlingRights = NoCastling
		}
	}

	// update
	newPos.SideToMove = pos.SideToMove.Opponent()
	newPos.Ply++
	if move.IsCapture() || piece.Type == Pawn {
		newPos.HalfmoveClock = 0
	} else {
		newPos.HalfmoveClock++
	}

	if move.IsDoublePush() {
		if piece == BlackPawn {
			newPos.EnPassant = to + 8
		} else {
			newPos.EnPassant = to - 8
		}
	} else {
		newPos.EnPassant = NoSquare
	}

	return newPos
}
