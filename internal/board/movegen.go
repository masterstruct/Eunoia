package board

// Thank you Dan (Hobbes) for your movegen implementation this file is based on
// https://github.com/kelseyde/hobbes-chess-engine/blob/main/src/board/movegen.rs
// also inspired by Reckless

// All - all moves
// Quiets - non-captures and underpromotions
// Noisies - captures and queen promos
type MoveFilter uint8

const (
	All MoveFilter = iota
	Quiets
	Noisies
)

func (mf MoveFilter) genQuiets() bool {
	return mf == All || mf == Quiets
}

func (mf MoveFilter) genNoisies() bool {
	return mf == All || mf == Noisies
}

var (
	underPromos = [...]PieceType{Knight, Bishop, Rook}
	allPromos   = [...]PieceType{Knight, Bishop, Rook, Queen}
)

// returns true if the side to move is in check
func (pos *Position) InCheck() bool {
	return pos.Checkers != 0
}

func GenKnightMoves(pos *Position, filterMask Bitboard, movelist *Movelist) {
	stm := pos.SideToMove
	them := pos.Them()
	occ := pos.Occ()
	pinned := pos.Pinned[stm]

	// loop over each knight
	knights := pos.Knights(stm) &^ pinned
	for knights != 0 {
		from := knights.PopLSB()
		attacks := KnightAttacks[from] & filterMask
		collectAttackMoves(movelist, from, attacks, them, occ)
	}
}

func collectSliderMoves(pos *Position, pieces, filterMask Bitboard, attacks func(Square, Bitboard) Bitboard, movelist *Movelist) {
	stm := pos.SideToMove
	them := pos.Them()
	occ := pos.Occ()
	kingSq := pos.KingSq[stm]
	pinned := pos.Pinned[stm]

	free := pieces &^ pinned
	for free != 0 {
		from := free.PopLSB()
		attacks := attacks(from, occ) & filterMask
		collectAttackMoves(movelist, from, attacks, them, occ)
	}

	stuck := pieces & pinned
	for stuck != 0 {
		from := stuck.PopLSB()
		attacks := attacks(from, occ) & Extending(kingSq, from) & filterMask
		collectAttackMoves(movelist, from, attacks, them, occ)
	}
}

func collectAttackMoves(movelist *Movelist, from Square, attacks, them, occ Bitboard) {
	captures := attacks & them
	quiets := attacks &^ occ
	for captures != 0 {
		// capture
		movelist.Add(NewCapture(from, captures.PopLSB()))
	}
	for quiets != 0 {
		// quiet move
		movelist.Add(NewMove(from, quiets.PopLSB()))
	}
}

func GenBishopMoves(pos *Position, filterMask Bitboard, movelist *Movelist) {
	collectSliderMoves(pos, pos.Bishops(pos.SideToMove), filterMask, BishopAttacks, movelist)
}

func GenRookMoves(pos *Position, filterMask Bitboard, movelist *Movelist) {
	collectSliderMoves(pos, pos.Rooks(pos.SideToMove), filterMask, RookAttacks, movelist)
}

func GenQueenMoves(pos *Position, filterMask Bitboard, movelist *Movelist) {
	collectSliderMoves(pos, pos.Queens(pos.SideToMove), filterMask, QueenAttacks, movelist)
}

func GenPawnMoves(pos *Position, filter MoveFilter, evasionMask Bitboard, movelist *Movelist) {
	stm := pos.SideToMove
	pinned := pos.Pinned[stm]
	kingSq := pos.KingSq[stm]
	pawns := pos.Pawns(stm)

	if filter.genNoisies() {
		if pos.EnPassant != NoSquare {
			collectEnPassant(pos, pawns, movelist)
		}
		collectPawnCaptures(pos, pawns, pinned, kingSq, evasionMask, movelist)
	}

	up, promos := collectPawnPushes(pos, pawns, pinned, evasionMask, filter, movelist)
	if promos == EmptyBB {
		return
	}

	switch filter {
	case All:
		collectPawnPromos(movelist, promos, up, allPromos[:])
	case Quiets:
		collectPawnPromos(movelist, promos, up, underPromos[:])
	case Noisies:
		collectPawnPromos(movelist, promos, up, []PieceType{Queen})
	}
}

func collectEnPassant(pos *Position, pawns Bitboard, movelist *Movelist) {
	epSq := pos.EnPassant
	nstm := pos.SideToMove.Opponent()
	attackers := PawnAttacks[nstm][epSq] & pawns
	for attackers != 0 {
		move := NewEnPassant(attackers.PopLSB(), epSq)
		if pos.IsLegal(move) {
			movelist.Add(move)
		}
	}
}

func collectPawnCaptures(pos *Position, pawns, pinned Bitboard, kingSq Square, evasionMask Bitboard, movelist *Movelist) {
	stm := pos.SideToMove
	them := pos.Them()

	freePawns := pawns &^ pinned
	pawnsLeft := freePawns
	pawnsRight := freePawns

	if stm == Black {
		if pinned != EmptyBB {
			pawnsLeft |= pawns & pinned & Diagonal(kingSq, SouthEast)
			pawnsRight |= pawns & pinned & Diagonal(kingSq, SouthWest)
		}
		capturesLeft := (pawnsLeft &^ FileH.Bits()).Shift(-7) & them & evasionMask
		capturesRight := (pawnsRight &^ FileA.Bits()).Shift(-9) & them & evasionMask
		collectPawnCaptureMoves(movelist, capturesLeft, -7, Rank1.Bits())
		collectPawnCaptureMoves(movelist, capturesRight, -9, Rank1.Bits())
	} else {
		if pinned != EmptyBB {
			pawnsLeft |= pawns & pinned & Diagonal(kingSq, NorthWest)
			pawnsRight |= pawns & pinned & Diagonal(kingSq, NorthEast)
		}
		capturesLeft := (pawnsLeft &^ FileA.Bits()).Shift(7) & them & evasionMask
		capturesRight := (pawnsRight &^ FileH.Bits()).Shift(9) & them & evasionMask
		collectPawnCaptureMoves(movelist, capturesLeft, 7, Rank8.Bits())
		collectPawnCaptureMoves(movelist, capturesRight, 9, Rank8.Bits())
	}
}

func collectPawnPushes(pos *Position, pawns, pinned Bitboard, evasionMask Bitboard, filter MoveFilter, movelist *Movelist) (int8, Bitboard) {
	stm := pos.SideToMove
	empty := ^pos.Occ()

	up := int8(8)
	thirdRank := Rank3
	promoRank := Rank8
	if stm == Black {
		up = -8
		thirdRank = Rank6
		promoRank = Rank1
	}

	pushPawns := pawns & (^pinned | pos.KingSq[stm].File().Bits())
	oneStep := pushPawns.Shift(up) & empty
	promoMask := promoRank.Bits()
	promos := oneStep & promoMask & evasionMask

	if filter.genQuiets() {
		singlePushes := oneStep &^ promoMask & evasionMask
		collectPawnQuiets(movelist, singlePushes, up)

		doublePushes := (oneStep & thirdRank.Bits()).Shift(up) & empty & evasionMask
		collectPawnDoublePushes(movelist, doublePushes, up*2)
	}
	return up, promos
}

func collectPawnPromos(movelist *Movelist, destinations Bitboard, offset int8, promos []PieceType) {
	for destinations != 0 {
		to := destinations.PopLSB()
		from := Square(int(to) - int(offset))
		for _, promo := range promos {
			movelist.Add(NewPromo(from, to, promo))
		}
	}
}

func collectPawnQuiets(movelist *Movelist, destinations Bitboard, offset int8) {
	for destinations != 0 {
		to := destinations.PopLSB()
		from := Square(int(to) - int(offset))
		movelist.Add(NewMove(from, to))
	}
}

func collectPawnDoublePushes(movelist *Movelist, destinations Bitboard, offset int8) {
	for destinations != 0 {
		to := destinations.PopLSB()
		from := Square(int(to) - int(offset))
		movelist.Add(NewDoublePush(from, to))
	}
}

func collectPawnCaptureMoves(movelist *Movelist, destinations Bitboard, offset int8, promoRank Bitboard) {
	promos := destinations & promoRank
	destinations &^= promoRank
	for destinations != 0 {
		to := destinations.PopLSB()
		from := Square(int(to) - int(offset))
		movelist.Add(NewCapture(from, to))
	}
	for promos != 0 {
		to := promos.PopLSB()
		from := Square(int(to) - int(offset))
		for _, promo := range allPromos {
			movelist.Add(NewCapturePromo(from, to, promo))
		}
	}
}

func GenKingMoves(pos *Position, filterMask Bitboard, movelist *Movelist) {
	stm := pos.SideToMove
	from := pos.KingSq[stm]

	attacks := KingAttacks[from] &^ pos.Threats & filterMask
	collectAttackMoves(movelist, from, attacks, pos.Them(), pos.Occ())
}

func collectCastleMoves(pos *Position, stm Color, movelist *Movelist) {
	base := stm * 2
	kingSq := pos.KingSq[stm]

	queensideRook := pos.Castling[base]
	kingsideRook := pos.Castling[base+1]

	if kingsideRook != NoSquare && canCastle(pos, kingSq, kingsideRook) {
		movelist.Add(NewCastle(kingSq, kingsideRook))
	}
	if queensideRook != NoSquare && canCastle(pos, kingSq, queensideRook) {
		movelist.Add(NewCastle(kingSq, queensideRook))
	}
}

func canCastle(pos *Position, kingSq, rookSq Square) bool {
	kingside := rookSq > kingSq
	kingTo, rookTo := castleTargets(pos.SideToMove, kingside)

	kingTravel := Between(kingSq, kingTo) | kingTo.Bit()
	rookTravel := Between(rookSq, rookTo) | rookTo.Bit()

	travel := (kingTravel | rookTravel) &^ kingSq.Bit()
	occ := pos.Occ() &^ rookSq.Bit()
	if travel&occ != 0 {
		return false
	}

	safety := kingTravel | kingSq.Bit()
	if safety&pos.Threats != 0 {
		return false
	}

	return !pos.Pinned[pos.SideToMove].IsBitSet(rookSq)
}

func castleTargets(stm Color, kingside bool) (kingTo, rookTo Square) {
	rank := stm.ExpectedKingRank()
	if kingside {
		return NewSquare(FileG, rank), NewSquare(FileF, rank)
	}
	return NewSquare(FileC, rank), NewSquare(FileD, rank)
}

func GenerateLegalMoves(pos *Position, movelist *Movelist, filter MoveFilter) {
	stm := pos.SideToMove
	us := pos.Us()
	them := pos.Them()
	kingSq := pos.KingSq[stm]
	inCheck := pos.Checkers != 0

	filterMask := FullBB
	switch filter {
	case Quiets:
		filterMask = ^them
	case Noisies:
		filterMask = them
	}

	GenKingMoves(pos, filterMask, movelist)
	if !inCheck && filter.genQuiets() {
		collectCastleMoves(pos, stm, movelist)
	}

	if pos.Checkers.CountBits() > 1 {
		return
	}

	evasionMask := FullBB
	if inCheck {
		evasionMask = pos.Checkers | Between(kingSq, pos.Checkers.LSB())
	}

	filterMask &= evasionMask
	if filter == Quiets {
		filterMask &^= them
	} else {
		filterMask &^= us
	}

	GenKnightMoves(pos, filterMask, movelist)
	GenBishopMoves(pos, filterMask, movelist)
	GenRookMoves(pos, filterMask, movelist)
	GenQueenMoves(pos, filterMask, movelist)
	GenPawnMoves(pos, filter, evasionMask, movelist)
}

// squares attacked by opponent pieces
func (pos *Position) calculateThreats(stm Color) Bitboard {
	occ := pos.Occ() &^ pos.KingSq[stm].Bit()
	nstm := stm.Opponent()

	pawns := pos.Pawns(nstm)
	knights := pos.Knights(nstm)
	bishops := pos.Bishops(nstm)
	rooks := pos.Rooks(nstm)
	queens := pos.Queens(nstm)

	return pawnAttacksSetwise(pawns, nstm) |
		knightAttacksSetwise(knights) |
		bishopAttacksSetwise(bishops|queens, occ) |
		rookAttacksSetwise(rooks|queens, occ) |
		KingAttacks[pos.KingSq[nstm]]
}

// pieces attacking the king
func (pos *Position) calculateCheckers(stm Color) Bitboard {
	occ := pos.Occ()
	kingSq := pos.KingSq[stm]
	nstm := stm.Opponent()

	pawns := pos.Pawns(nstm)
	knights := pos.Knights(nstm)
	bishops := pos.Bishops(nstm)
	rooks := pos.Rooks(nstm)
	queens := pos.Queens(nstm)

	return PawnAttacks[stm][kingSq]&pawns |
		KnightAttacks[kingSq]&knights |
		BishopAttacks(kingSq, occ)&(bishops|queens) |
		RookAttacks(kingSq, occ)&(rooks|queens)
}

func (pos *Position) calculateBothPinned() [2]Bitboard {
	return [2]Bitboard{pos.calculatePinned(Black), pos.calculatePinned(White)}
}

func (pos *Position) calculatePinned(stm Color) Bitboard {
	king := pos.KingSq[stm]
	us := pos.Colors[stm]
	them := pos.Colors[stm.Opponent()]

	diagonals := pos.Bishops(stm.Opponent()) | pos.Queens(stm.Opponent())
	orthogonals := pos.Rooks(stm.Opponent()) | pos.Queens(stm.Opponent())

	if diagonals|orthogonals == EmptyBB {
		return EmptyBB
	}

	potentialAttackers := BishopAttacks(king, them)&diagonals |
		RookAttacks(king, them)&orthogonals

	pinned := EmptyBB

	for potentialAttackers != 0 {
		attacker := potentialAttackers.PopLSB()
		between := Between(king, attacker)
		maybePinned := us & between
		if maybePinned.CountBits() == 1 {
			pinned |= maybePinned
		}
	}

	return pinned
}
