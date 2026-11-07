package board

func (pos *Position) Make(move Move) {
	from, to := move.From(), move.To()
	piece := pos.PieceOn(from)
	pieceType := piece.Type()
	capturedPiece := pos.PieceOn(to)
	stm := pos.SideToMove
	nstm := stm.Opponent()

	isPromo := move.IsPromo()
	isEnPassant := move.IsEnPassant()
	isCapture := move.IsCapture()
	isCastle := move.IsCastle()

	hash := pos.Hash

	// remove existing castling rights and re-apply them at the end
	hash ^= ZobristTable.CastlingKey(pos.Castling.ToIndex())

	pos.HalfmoveClock++
	oldEnPassantSq := pos.EnPassant
	if oldEnPassantSq != NoSquare {
		pos.EnPassant = NoSquare
		hash ^= ZobristTable.EnPassantKey(oldEnPassantSq.File())
	}

	pos.RemovePiece(from)
	hash ^= ZobristTable.PieceKey(piece, from)
	pos.RemovePiece(to)

	if !isCastle {
		epSquare := to - (16*Square(stm) - 8)
		if isEnPassant {
			pos.RemovePiece(epSquare)
			hash ^= ZobristTable.PieceKey(NewPiece(Pawn, nstm), epSquare)
		} else if isCapture {
			hash ^= ZobristTable.PieceKey(capturedPiece, to)
		}

		if isPromo {
			p := NewPiece(move.Promo(), stm)
			pos.PlacePiece(p, to)
			hash ^= ZobristTable.PieceKey(p, to)
		} else {
			pos.PlacePiece(piece, to)
			hash ^= ZobristTable.PieceKey(piece, to)
		}

		// rook move - remove castling rights
		if pieceType == Rook {
			ok, sq := pos.Castling.HasSquare(from)
			if ok {
				pos.Castling.Remove(sq)
			}
		}

		// rook captured - remove castling rights
		if isCapture && capturedPiece.Type() == Rook {
			ok, sq := pos.Castling.HasSquare(to)
			if ok {
				pos.Castling.Remove(sq)
			}
		}

		if isCapture || pieceType == Pawn {
			pos.HalfmoveClock = 0
		}

		if move.IsDoublePush() {
			themPawns := pos.Pawns(nstm)
			if (to.File() != FileA && themPawns.IsBitSet(to.Left())) ||
				(to.File() != FileH && themPawns.IsBitSet(to.Right())) {
				pos.EnPassant = epSquare
				hash ^= ZobristTable.EnPassantKey(epSquare.File())
			}
		}

		// king moved - remove castling rights
		if pieceType == King {
			pos.KingSq[stm] = to
			pos.Castling.Clear(stm)
		}
	} else {
		// castle - move pieces

		// remove rook
		hash ^= ZobristTable.PieceKey(NewPiece(Rook, stm), to)
		pos.Castling.Clear(stm)

		if stm == Black {
			if move.IsKingsideCastle() {
				pos.PlacePiece(BlackKing, G8)
				hash ^= ZobristTable.PieceKey(BlackKing, G8)
				pos.PlacePiece(BlackRook, F8)
				hash ^= ZobristTable.PieceKey(BlackRook, F8)
				pos.KingSq[Black] = G8
			} else {
				pos.PlacePiece(BlackKing, C8)
				hash ^= ZobristTable.PieceKey(BlackKing, C8)
				pos.PlacePiece(BlackRook, D8)
				hash ^= ZobristTable.PieceKey(BlackRook, D8)
				pos.KingSq[Black] = C8
			}
		} else {
			if move.IsKingsideCastle() {
				pos.PlacePiece(WhiteKing, G1)
				hash ^= ZobristTable.PieceKey(WhiteKing, G1)
				pos.PlacePiece(WhiteRook, F1)
				hash ^= ZobristTable.PieceKey(WhiteRook, F1)
				pos.KingSq[White] = G1
			} else {
				pos.PlacePiece(WhiteKing, C1)
				hash ^= ZobristTable.PieceKey(WhiteKing, C1)
				pos.PlacePiece(WhiteRook, D1)
				hash ^= ZobristTable.PieceKey(WhiteRook, D1)
				pos.KingSq[White] = C1
			}
		}
	}

	pos.SideToMove = nstm
	hash ^= ZobristTable.SideToMoveKey()
	pos.Ply++

	// re-apply castling rights
	hash ^= ZobristTable.CastlingKey(pos.Castling.ToIndex())
	pos.Hash = hash

	pos.Threats = pos.calculateThreats(nstm)
	pos.Checkers = pos.calculateCheckers(nstm)
	pos.Pinned = pos.calculateBothPinned()
}

func (pos *Position) MakeNullMove() {
	if pos.EnPassant != NoSquare {
		pos.Hash ^= ZobristTable.EnPassantKey(pos.EnPassant.File())
		pos.EnPassant = NoSquare
	}

	pos.SideToMove = pos.SideToMove.Opponent()
	pos.Hash ^= ZobristTable.SideToMoveKey()
	pos.Ply++

	pos.Threats = pos.calculateThreats(pos.SideToMove)
	pos.Checkers = pos.calculateCheckers(pos.SideToMove)
}
