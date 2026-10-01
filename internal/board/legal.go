package board

func (pos *Position) IsLegal(move Move) bool {
	from := move.From()
	to := move.To()

	if move.IsEnPassant() {
		return isLegalEnPassant(pos, from, to)
	}

	kingSq := pos.KingSq[pos.SideToMove]

	if move.IsCastle() {
		return canCastle(pos, kingSq, to)
	}

	if from == kingSq {
		return !pos.Threats.IsBitSet(to)
	}

	if pos.Pinned[pos.SideToMove].IsBitSet(from) {
		return pos.Checkers == 0 && Extending(kingSq, from).IsBitSet(to)
	}

	if pos.Checkers.CountBits() > 1 {
		return false
	}

	if pos.Checkers == 0 {
		return true
	}

	checker := pos.Checkers.LSB()
	return (pos.Checkers | Between(kingSq, checker)).IsBitSet(to)
}

func isLegalEnPassant(pos *Position, from, epSq Square) bool {
	stm := pos.SideToMove
	them := pos.Colors[stm.Opponent()]
	kingSq := pos.KingSq[stm]

	capturedSq := epSq - (16*Square(stm) - 8)

	occ := pos.Occupied()
	occ.ClearBit(from)
	occ.ClearBit(capturedSq)
	occ.SetBit(epSq)

	bishops := pos.Pieces[Bishop] & them
	rooks := pos.Pieces[Rook] & them
	queens := pos.Pieces[Queen] & them

	return BishopAttacks(kingSq, occ)&(bishops|queens) == 0 &&
		RookAttacks(kingSq, occ)&(rooks|queens) == 0
}
