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
	oppColor := color.Opponent()

	isPromo := move.IsPromo()
	isEnPassant := move.IsEnPassant()
	isCapture := move.IsCapture()
	isCastle := move.IsCastle()

	hash := newPos.Hash

	// remove existing castling rights and re-apply them at the end
	hash ^= ZobristTable.CastlingKey(newPos.CastlingRights)

	newPos.HalfmoveClock++
	oldEnPassantSq := newPos.EnPassant
	if oldEnPassantSq != NoSquare {
		newPos.EnPassant = NoSquare
		hash ^= ZobristTable.EnPassantKey(oldEnPassantSq.File())
	}

	newPos.RemovePiece(from)
	hash ^= ZobristTable.PieceKey(color, pieceType, from)
	newPos.RemovePiece(to)

	if !isCastle {
		epSquare := to - (16*Square(color) - 8)
		if isEnPassant {
			newPos.RemovePiece(epSquare)
			hash ^= ZobristTable.PieceKey(oppColor, Pawn, epSquare)
		} else if isCapture {
			hash ^= ZobristTable.PieceKey(oppColor, capturedPiece.Type, to)
		}

		if isPromo {
			promo := move.Promo()
			newPos.PlacePiece(NewPiece(promo, color), to)
			hash ^= ZobristTable.PieceKey(color, promo, to)
		} else {
			newPos.PlacePiece(piece, to)
			hash ^= ZobristTable.PieceKey(color, pieceType, to)
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
			oppPawns := newPos.PieceBB(Piece{Type: Pawn, Color: oppColor})
			if (to.File() != FileA && oppPawns.IsBitSet(to.Left())) ||
				(to.File() != FileH && oppPawns.IsBitSet(to.Right())) {
				newPos.EnPassant = epSquare
				hash ^= ZobristTable.EnPassantKey(epSquare.File())
			}
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

		// remove rook
		hash ^= ZobristTable.PieceKey(color, Rook, to)

		if color == Black {
			if move.IsKingsideCastle() {
				newPos.PlacePiece(BlackKing, G8)
				hash ^= ZobristTable.PieceKey(Black, King, G8)
				newPos.PlacePiece(BlackRook, F8)
				hash ^= ZobristTable.PieceKey(Black, Rook, F8)
				newPos.KingSq[Black] = G8
			} else {
				newPos.PlacePiece(BlackKing, C8)
				hash ^= ZobristTable.PieceKey(Black, King, C8)
				newPos.PlacePiece(BlackRook, D8)
				hash ^= ZobristTable.PieceKey(Black, Rook, D8)
				newPos.KingSq[Black] = C8
			}
			newPos.CastlingRights.Remove(BlackKingside | BlackQueenside)
		} else {
			if move.IsKingsideCastle() {
				newPos.PlacePiece(WhiteKing, G1)
				hash ^= ZobristTable.PieceKey(White, King, G1)
				newPos.PlacePiece(WhiteRook, F1)
				hash ^= ZobristTable.PieceKey(White, Rook, F1)
				newPos.KingSq[White] = G1
			} else {
				newPos.PlacePiece(WhiteKing, C1)
				hash ^= ZobristTable.PieceKey(White, King, C1)
				newPos.PlacePiece(WhiteRook, D1)
				hash ^= ZobristTable.PieceKey(White, Rook, D1)
				newPos.KingSq[White] = C1
			}
			newPos.CastlingRights.Remove(WhiteKingside | WhiteQueenside)
		}
	}

	newPos.SideToMove = oppColor
	hash ^= ZobristTable.SideToMoveKey()
	newPos.Ply++

	// re-apply castling rights
	hash ^= ZobristTable.CastlingKey(newPos.CastlingRights)

	newPos.Hash = hash

	return newPos
}
