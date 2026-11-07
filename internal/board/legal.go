package board

func (pos *Position) IsLegal(move Move) bool {
	stm := pos.SideToMove
	from := move.From()
	to := move.To()

	if !pos.Us().IsBitSet(from) {
		return false
	}

	if move.IsEnPassant() {
		return isLegalEnPassant(pos, from, to)
	}

	kingSq := pos.KingSq[stm]

	if move.IsCastle() {
		return canCastle(pos, kingSq, to)
	}

	if from == kingSq {
		return !pos.Threats.IsBitSet(to)
	}

	if pos.Pinned[stm].IsBitSet(from) {
		return pos.Checkers.IsEmpty() && Extending(kingSq, from).IsBitSet(to)
	}

	if pos.Checkers.CountBits() > 1 {
		return false
	}

	if pos.Checkers.IsEmpty() {
		return true
	}

	checker := pos.Checkers.LSB()
	return (pos.Checkers | Between(kingSq, checker)).IsBitSet(to)
}

func isLegalEnPassant(pos *Position, from, epSq Square) bool {
	stm := pos.SideToMove
	kingSq := pos.KingSq[stm]

	capturedSq := epSq - (16*Square(stm) - 8)

	occ := pos.Occ()
	occ.ClearBit(from)
	occ.ClearBit(capturedSq)
	occ.SetBit(epSq)

	bishops := pos.Their(Bishop)
	rooks := pos.Their(Rook)
	queens := pos.Their(Queen)

	return (BishopAttacks(kingSq, occ) & (bishops | queens)).IsEmpty() &&
		(RookAttacks(kingSq, occ) & (rooks | queens)).IsEmpty()
}
