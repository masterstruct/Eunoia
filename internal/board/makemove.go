package board

func (pos Position) MakeMove(move Move) Position {
	newPos := pos
	from, to := move.From(), move.To()
	piece, _ := pos.PieceOn(from)

	// move/capture piece
	newPos.RemovePiece(to)
	newPos.PlacePiece(pos.Board[from], to)
	newPos.RemovePiece(from)

	// update
	newPos.SideToMove = pos.SideToMove.Opponent()
	newPos.Ply++
	if move.IsCapture() || piece.Type == Pawn {
		newPos.HalfmoveClock = 0
	}

	return newPos
}
