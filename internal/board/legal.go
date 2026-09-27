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
	color := pos.SideToMove
	enemyBB := pos.Colors[color.Opponent()]
	kingSq := pos.KingSq[color]

	capturedSq := epSq - (16*Square(color) - 8)

	occupied := pos.Occupied()
	occupied.ClearBit(from)
	occupied.ClearBit(capturedSq)
	occupied.SetBit(epSq)

	bishops := pos.Pieces[Bishop] & enemyBB
	rooks := pos.Pieces[Rook] & enemyBB
	queens := pos.Pieces[Queen] & enemyBB

	return BishopAttacks(kingSq, occupied)&(bishops|queens) == 0 &&
		RookAttacks(kingSq, occupied)&(rooks|queens) == 0
}
