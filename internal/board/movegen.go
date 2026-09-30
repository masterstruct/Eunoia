package board

// Thank you Dan, creator of Hobbes, for your movegen implementation this file is based on
// https://github.com/kelseyde/hobbes-chess-engine/blob/main/src/board/movegen.rs

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

const MaxMoves = 256

type Movelist struct {
	Moves [MaxMoves]Move
	Len   int
}

func (ml *Movelist) Add(m Move) {
	ml.Moves[ml.Len] = m
	ml.Len++
}

func IsSquareAttacked(pos *Position, sq Square, byColor Color) bool {
	return isSquareAttacked(pos, sq, byColor, pos.Occupied())
}

func isSquareAttacked(pos *Position, sq Square, byColor Color, occupied Bitboard) bool {
	if byColor == NoColor {
		return false
	}

	if PawnAttacks[byColor.Opponent()][sq]&pos.PieceBB(Piece{Type: Pawn, Color: byColor}) != 0 {
		return true
	}
	if KnightAttacks[sq]&pos.PieceBB(Piece{Type: Knight, Color: byColor}) != 0 {
		return true
	}
	if KingAttacks[sq]&pos.PieceBB(Piece{Type: King, Color: byColor}) != 0 {
		return true
	}

	rooks := pos.PieceBB(Piece{Type: Rook, Color: byColor})
	queens := pos.PieceBB(Piece{Type: Queen, Color: byColor})

	rookAttacks := RookAttacks(sq, occupied)
	if rookAttacks&(rooks|queens) != 0 {
		return true
	}

	bishops := pos.PieceBB(Piece{Type: Bishop, Color: byColor})

	bishopAttacks := BishopAttacks(sq, occupied)
	if bishopAttacks&(bishops|queens) != 0 {
		return true
	}
	return false
}

// returns true if the side to move is in check
func (pos *Position) InCheck() bool {
	return pos.Checkers != 0
}

func GenKnightMoves(pos *Position, filterMask Bitboard, movelist *Movelist) {
	color := pos.SideToMove
	opponents := pos.Colors[color.Opponent()]
	occupied := pos.Occupied()
	pinned := pos.Pinned[color]

	// loop over each knight
	knights := pos.PieceBB(Piece{Type: Knight, Color: color}) &^ pinned
	for knights != 0 {
		from := knights.PopLSB()
		attacks := KnightAttacks[from] & filterMask
		addAttackMoves(movelist, from, attacks, opponents, occupied)
	}
}

func genSliderMoves(pos *Position, pieces, filterMask Bitboard, attacks func(Square, Bitboard) Bitboard, movelist *Movelist) {
	color := pos.SideToMove
	opponents := pos.Colors[color.Opponent()]
	occupied := pos.Occupied()
	kingSq := pos.KingSq[color]
	pinned := pos.Pinned[color]

	free := pieces &^ pinned
	for free != 0 {
		from := free.PopLSB()
		attacks := attacks(from, occupied) & filterMask
		addAttackMoves(movelist, from, attacks, opponents, occupied)
	}

	stuck := pieces & pinned
	for stuck != 0 {
		from := stuck.PopLSB()
		attacks := attacks(from, occupied) & Extending(kingSq, from) & filterMask
		addAttackMoves(movelist, from, attacks, opponents, occupied)
	}
}

func addAttackMoves(movelist *Movelist, from Square, attacks, opponents, occupied Bitboard) {
	captures := attacks & opponents
	quiets := attacks &^ occupied
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
	bishops := pos.PieceBB(Piece{Type: Bishop, Color: pos.SideToMove})
	genSliderMoves(pos, bishops, filterMask, BishopAttacks, movelist)
}

func GenRookMoves(pos *Position, filterMask Bitboard, movelist *Movelist) {
	rooks := pos.PieceBB(Piece{Type: Rook, Color: pos.SideToMove})
	genSliderMoves(pos, rooks, filterMask, RookAttacks, movelist)
}

func GenQueenMoves(pos *Position, filterMask Bitboard, movelist *Movelist) {
	queens := pos.PieceBB(Piece{Type: Queen, Color: pos.SideToMove})
	genSliderMoves(pos, queens, filterMask, QueenAttacks, movelist)
}

func GenPawnMoves(pos *Position, filter MoveFilter, evasionMask Bitboard, movelist *Movelist) {
	color := pos.SideToMove
	pinned := pos.Pinned[color]
	kingSq := pos.KingSq[color]
	pawns := pos.PieceBB(Piece{Type: Pawn, Color: color})

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
	attackers := PawnAttacks[pos.SideToMove.Opponent()][epSq] & pawns
	for attackers != 0 {
		move := NewEnPassant(attackers.PopLSB(), epSq)
		if pos.IsLegal(move) {
			movelist.Add(move)
		}
	}
}

func collectPawnCaptures(pos *Position, pawns, pinned Bitboard, kingSq Square, evasionMask Bitboard, movelist *Movelist) {
	color := pos.SideToMove
	enemyOcc := pos.Colors[color.Opponent()]

	freePawns := pawns &^ pinned
	pawnsLeft := freePawns
	pawnsRight := freePawns

	if color == Black {
		if pinned != EmptyBB {
			pawnsLeft |= pawns & pinned & Diagonal(kingSq, SouthEast)
			pawnsRight |= pawns & pinned & Diagonal(kingSq, SouthWest)
		}
		capturesLeft := (pawnsLeft &^ FileH.Bits()).Shift(-7) & enemyOcc & evasionMask
		capturesRight := (pawnsRight &^ FileA.Bits()).Shift(-9) & enemyOcc & evasionMask
		collectPawnCaptureMoves(movelist, capturesLeft, -7, Rank1.Bits())
		collectPawnCaptureMoves(movelist, capturesRight, -9, Rank1.Bits())
	} else {
		if pinned != EmptyBB {
			pawnsLeft |= pawns & pinned & Diagonal(kingSq, NorthWest)
			pawnsRight |= pawns & pinned & Diagonal(kingSq, NorthEast)
		}
		capturesLeft := (pawnsLeft &^ FileA.Bits()).Shift(7) & enemyOcc & evasionMask
		capturesRight := (pawnsRight &^ FileH.Bits()).Shift(9) & enemyOcc & evasionMask
		collectPawnCaptureMoves(movelist, capturesLeft, 7, Rank8.Bits())
		collectPawnCaptureMoves(movelist, capturesRight, 9, Rank8.Bits())
	}
}

func collectPawnPushes(pos *Position, pawns, pinned Bitboard, evasionMask Bitboard, filter MoveFilter, movelist *Movelist) (int8, Bitboard) {
	color := pos.SideToMove
	empty := ^pos.Occupied()

	up := int8(8)
	thirdRank := Rank3
	promoRank := Rank8
	if color == Black {
		up = -8
		thirdRank = Rank6
		promoRank = Rank1
	}

	pushPawns := pawns & (^pinned | pos.KingSq[color].File().Bits())
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
	color := pos.SideToMove
	from := pos.KingSq[color]

	attacks := KingAttacks[from] &^ pos.Threats & filterMask
	addAttackMoves(movelist, from, attacks, pos.Colors[color.Opponent()], pos.Occupied())
}

func genCastleMoves(pos *Position, color Color, movelist *Movelist) {
	base := color * 2
	kingSq := pos.KingSq[color]

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
	occupied := pos.Occupied() &^ rookSq.Bit()
	if travel&occupied != 0 {
		return false
	}

	safety := kingTravel | kingSq.Bit()
	if safety&pos.Threats != 0 {
		return false
	}

	return !pos.Pinned[pos.SideToMove].IsBitSet(rookSq)
}

func castleTargets(color Color, kingside bool) (kingTo, rookTo Square) {
	rank := color.ExpectedKingRank()
	if kingside {
		return NewSquare(FileG, rank), NewSquare(FileF, rank)
	}
	return NewSquare(FileC, rank), NewSquare(FileD, rank)
}

func GenerateLegalMoves(pos *Position, movelist *Movelist, filter MoveFilter) {
	color := pos.SideToMove
	us := pos.Colors[color]
	them := pos.Colors[color.Opponent()]
	kingSq := pos.KingSq[color]
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
		genCastleMoves(pos, color, movelist)
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
func (pos *Position) calculateThreats(myColor Color) Bitboard {
	occupied := pos.Occupied() &^ pos.KingSq[myColor].Bit()
	enemyColor := myColor.Opponent()
	enemyBB := pos.Colors[enemyColor]

	pawns := pos.Pieces[Pawn] & enemyBB
	knights := pos.Pieces[Knight] & enemyBB
	bishops := pos.Pieces[Bishop] & enemyBB
	rooks := pos.Pieces[Rook] & enemyBB
	queens := pos.Pieces[Queen] & enemyBB

	return pawnAttacksSetwise(pawns, enemyColor) |
		knightAttacksSetwise(knights) |
		bishopAttacksSetwise(bishops|queens, occupied) |
		rookAttacksSetwise(rooks|queens, occupied) |
		KingAttacks[pos.KingSq[enemyColor]]
}

// pieces attacking the king
func (pos *Position) calculateCheckers(myColor Color) Bitboard {
	occupied := pos.Occupied()
	kingSq := pos.KingSq[myColor]
	enemyBB := pos.Colors[myColor.Opponent()]

	pawns := pos.Pieces[Pawn] & enemyBB
	knights := pos.Pieces[Knight] & enemyBB
	bishops := pos.Pieces[Bishop] & enemyBB
	rooks := pos.Pieces[Rook] & enemyBB
	queens := pos.Pieces[Queen] & enemyBB

	return PawnAttacks[myColor][kingSq]&pawns |
		KnightAttacks[kingSq]&knights |
		BishopAttacks(kingSq, occupied)&(bishops|queens) |
		RookAttacks(kingSq, occupied)&(rooks|queens)
}

func (pos *Position) calculateBothPinned() [2]Bitboard {
	return [2]Bitboard{pos.calculatePinned(Black), pos.calculatePinned(White)}
}

func (pos *Position) calculatePinned(color Color) Bitboard {
	king := pos.KingSq[color]
	us := pos.Colors[color]
	them := pos.Colors[color.Opponent()]

	diagonals := (pos.Pieces[Bishop] | pos.Pieces[Queen]) & them
	orthogonals := (pos.Pieces[Rook] | pos.Pieces[Queen]) & them

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
