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
	stm := newPos.SideToMove
	nstm := stm.Opponent()

	isPromo := move.IsPromo()
	isEnPassant := move.IsEnPassant()
	isCapture := move.IsCapture()
	isCastle := move.IsCastle()

	hash := newPos.Hash

	// remove existing castling rights and re-apply them at the end
	hash ^= ZobristTable.CastlingKey(newPos.Castling.ToIndex())

	newPos.HalfmoveClock++
	oldEnPassantSq := newPos.EnPassant
	if oldEnPassantSq != NoSquare {
		newPos.EnPassant = NoSquare
		hash ^= ZobristTable.EnPassantKey(oldEnPassantSq.File())
	}

	newPos.RemovePiece(from)
	hash ^= ZobristTable.PieceKey(stm, pieceType, from)
	newPos.RemovePiece(to)

	if !isCastle {
		epSquare := to - (16*Square(stm) - 8)
		if isEnPassant {
			newPos.RemovePiece(epSquare)
			hash ^= ZobristTable.PieceKey(nstm, Pawn, epSquare)
		} else if isCapture {
			hash ^= ZobristTable.PieceKey(nstm, capturedPiece.Type, to)
		}

		if isPromo {
			promo := move.Promo()
			newPos.PlacePiece(NewPiece(promo, stm), to)
			hash ^= ZobristTable.PieceKey(stm, promo, to)
		} else {
			newPos.PlacePiece(piece, to)
			hash ^= ZobristTable.PieceKey(stm, pieceType, to)
		}

		// rook move - remove castling rights
		if pieceType == Rook {
			ok, sq := pos.Castling.HasSquare(from)
			if ok {
				newPos.Castling.Remove(sq)
			}
		}

		// rook captured - remove castling rights
		if isCapture && capturedPiece.Type == Rook {
			ok, sq := pos.Castling.HasSquare(to)
			if ok {
				newPos.Castling.Remove(sq)
			}
		}

		if isCapture || pieceType == Pawn {
			newPos.HalfmoveClock = 0
		}

		if move.IsDoublePush() {
			themPawns := newPos.PieceBB(Piece{Type: Pawn, Color: nstm})
			if (to.File() != FileA && themPawns.IsBitSet(to.Left())) ||
				(to.File() != FileH && themPawns.IsBitSet(to.Right())) {
				newPos.EnPassant = epSquare
				hash ^= ZobristTable.EnPassantKey(epSquare.File())
			}
		}

		// king moved - remove castling rights
		if pieceType == King {
			newPos.KingSq[stm] = to
			newPos.Castling.Clear(stm)
		}
	} else {
		// castle - move pieces

		// remove rook
		hash ^= ZobristTable.PieceKey(stm, Rook, to)
		newPos.Castling.Clear(stm)

		if stm == Black {
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
		}
	}

	newPos.SideToMove = nstm
	hash ^= ZobristTable.SideToMoveKey()
	newPos.Ply++

	// re-apply castling rights
	hash ^= ZobristTable.CastlingKey(newPos.Castling.ToIndex())
	newPos.Hash = hash

	newPos.Threats = newPos.calculateThreats(nstm)
	newPos.Checkers = newPos.calculateCheckers(nstm)
	newPos.Pinned = newPos.calculateBothPinned()

	return newPos
}

func (pos *Position) MakeNullMove() Position {
	newPos := *pos
	if newPos.EnPassant != NoSquare {
		newPos.Hash ^= ZobristTable.EnPassantKey(newPos.EnPassant.File())
		newPos.EnPassant = NoSquare
	}

	newPos.SideToMove = pos.SideToMove.Opponent()
	newPos.Hash ^= ZobristTable.SideToMoveKey()
	newPos.Ply++

	newPos.Threats = newPos.calculateThreats(newPos.SideToMove)
	newPos.Checkers = newPos.calculateCheckers(newPos.SideToMove)

	return newPos
}
