package board

func (pos Position) MakeMove(move Move) Position {
	newPos := pos
	from, to := move.From(), move.To()
	piece, _ := pos.PieceOn(from)
	capturedPiece, _ := pos.PieceOn(to)
	color := pos.SideToMove

	isPromo := move.IsPromo()
	isEnPassant := move.IsEnPassant()
	isCapture := move.IsCapture()
	isCastle := move.IsCastle()

	pawnMoveDir := 8
	if color == Black {
		pawnMoveDir = -8
	}
	epSquare := Square(int(to) - pawnMoveDir)

	// move/capture/promote piece
	newPos.RemovePiece(from)
	if !isCastle {
		newPos.RemovePiece(to)
	}
	if isPromo {
		newPos.PlacePiece(NewPiece(move.Promo(), color), to)
	} else if !isCastle {
		newPos.PlacePiece(piece, to)
	}

	if isEnPassant {
		newPos.RemovePiece(epSquare)
	}

	// castling rights from rook moves/captures
	if piece.Type == Rook {
		rooks := pos.CastlingRooks()
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

	if isCapture && capturedPiece.Type == Rook {
		rooks := pos.CastlingRooks()
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

	if isCastle {
		if isCastle {
			rooks := pos.CastlingRooks()

			if color == Black {
				newPos.CastlingRights.Remove(BlackKingside | BlackQueenside)

				if move.IsKingsideCastle() {
					newPos.RemovePiece(rooks.BlackKingside)
					newPos.PlacePiece(BlackKing, G8)
					newPos.PlacePiece(BlackRook, F8)
				} else {
					newPos.RemovePiece(rooks.BlackQueenside)
					newPos.PlacePiece(BlackKing, C8)
					newPos.PlacePiece(BlackRook, D8)
				}
			} else {
				newPos.CastlingRights.Remove(WhiteKingside | WhiteQueenside)

				if move.IsKingsideCastle() {
					newPos.RemovePiece(rooks.WhiteKingside)
					newPos.PlacePiece(WhiteKing, G1)
					newPos.PlacePiece(WhiteRook, F1)
				} else {
					newPos.RemovePiece(rooks.WhiteQueenside)
					newPos.PlacePiece(WhiteKing, C1)
					newPos.PlacePiece(WhiteRook, D1)
				}
			}
		}

	} else if piece.Type == King {
		if color == Black {
			newPos.CastlingRights.Remove(BlackKingside)
			newPos.CastlingRights.Remove(BlackQueenside)
		} else {
			newPos.CastlingRights.Remove(WhiteKingside)
			newPos.CastlingRights.Remove(WhiteQueenside)
		}
	}

	// update state
	newPos.SideToMove = color.Opponent()
	newPos.Ply++
	if isCapture || piece.Type == Pawn {
		newPos.HalfmoveClock = 0
	} else {
		newPos.HalfmoveClock++
	}

	newPos.EnPassant = NoSquare
	if move.IsDoublePush() {
		newPos.EnPassant = epSquare
	}

	return newPos
}
