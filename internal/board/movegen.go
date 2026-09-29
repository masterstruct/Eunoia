package board

// All - all moves
// Quiets - non-captures (including promos)
// Noisies - captures and promos
// Captures - only captures
type MoveFilter uint8

const (
	All MoveFilter = iota
	Quiets
	Noisies
	Captures
)

func (mf MoveFilter) genQuiets() bool {
	return mf == All || mf == Quiets
}

func (mf MoveFilter) genNoisies() bool {
	return mf == All || mf == Noisies
}

func (mf MoveFilter) genCaptures() bool {
	return mf == All || mf == Noisies || mf == Captures
}

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
	enemyColor := color.Opponent()
	occupied := pos.Occupied()
	enemyOcc := pos.Colors[enemyColor]
	pinned := pos.Pinned[color]
	kingSq := pos.KingSq[color]

	pawns := pos.PieceBB(Piece{Type: Pawn, Color: color})

	if filter.genCaptures() {
		epSq := pos.EnPassant
		if epSq != NoSquare {
			attackers := PawnAttacks[enemyColor][epSq] & pawns
			for attackers != 0 {
				mv := NewEnPassant(attackers.PopLSB(), epSq)
				if pos.IsLegal(mv) {
					movelist.Add(mv)
				}
			}
		}
	}

	for pawns != 0 {
		from := pawns.PopLSB()

		pinRay := FullBB
		if pinned.IsBitSet(from) {
			pinRay = Extending(kingSq, from)
		}

		if filter.genCaptures() {
			captures := PawnAttacks[color][from] & enemyOcc & pinRay & evasionMask
			for captures != 0 {
				addPawnMove(movelist, from, captures.PopLSB(), true)
			}
		}

		if !filter.genQuiets() && !filter.genNoisies() {
			continue
		}

		to := from.Up()
		startRank := Rank2
		promoRank := Rank8
		if color == Black {
			to = from.Down()
			startRank = Rank7
			promoRank = Rank1
		}

		if (occupied | ^pinRay).IsBitSet(to) {
			continue
		}

		if to.Rank() == promoRank {
			if filter.genNoisies() && evasionMask.IsBitSet(to) {
				addPawnMove(movelist, from, to, false)
			}
			continue
		}

		if filter.genQuiets() && evasionMask.IsBitSet(to) {
			movelist.Add(NewMove(from, to))
		}

		if from.Rank() == startRank && filter.genQuiets() {
			var to2 Square
			if color == White {
				to2 = from + 16
			} else {
				to2 = from - 16
			}
			if (^occupied & pinRay & evasionMask).IsBitSet(to2) {
				movelist.Add(NewDoublePush(from, to2))
			}
		}
	}
}

func addPawnMove(movelist *Movelist, from, to Square, capture bool) {
	if to.Rank() == Rank1 || to.Rank() == Rank8 {
		if capture {
			movelist.Add(NewCapturePromo(from, to, Knight))
			movelist.Add(NewCapturePromo(from, to, Bishop))
			movelist.Add(NewCapturePromo(from, to, Rook))
			movelist.Add(NewCapturePromo(from, to, Queen))
		} else {
			movelist.Add(NewPromo(from, to, Knight))
			movelist.Add(NewPromo(from, to, Bishop))
			movelist.Add(NewPromo(from, to, Rook))
			movelist.Add(NewPromo(from, to, Queen))
		}
		return
	}
	if capture {
		movelist.Add(NewCapture(from, to))
	} else {
		movelist.Add(NewMove(from, to))
	}
}

func GenKingMoves(pos *Position, filterMask Bitboard, movelist *Movelist) {
	color := pos.SideToMove
	from := pos.KingSq[color]

	attacks := KingAttacks[from] &^ pos.Threats & filterMask
	addAttackMoves(movelist, from, attacks, pos.Colors[color.Opponent()], pos.Occupied())
}

func genCastleMoves(pos *Position, color Color, movelist *Movelist) {
	base := int(color * 2)
	kingSq := pos.KingSq[color]

	for _, kingside := range []bool{true, false} {
		right := CastlingRights(base)
		if kingside {
			right++
		}
		if !pos.Castling.Has(right) {
			continue
		}
		rookSq := pos.Castling[right]
		if canCastle(pos, kingSq, rookSq) {
			movelist.Add(NewCastle(kingSq, rookSq))
		}
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
	if !pos.Chess960 {
		return safety&pos.Threats == 0
	}

	for safety != 0 {
		if isSquareAttacked(pos, safety.PopLSB(), pos.SideToMove.Opponent(), occupied) {
			return false
		}
	}
	return true
}

func castleTargets(color Color, kingside bool) (kingTo, rookTo Square) {
	rank := color.ExpectedKingRank()
	if kingside {
		return NewSquare(FileG, rank), NewSquare(FileF, rank)
	}
	return NewSquare(FileC, rank), NewSquare(FileD, rank)
}

func GenerateLegalMovesForPiece(pos *Position, pieceType PieceType, movelist *Movelist, filter MoveFilter) {
	color := pos.SideToMove
	us := pos.Colors[color]
	them := pos.Colors[color.Opponent()]
	kingSq := pos.KingSq[color]
	inCheck := pos.Checkers != 0

	filterMask := FullBB
	switch filter {
	case Quiets:
		filterMask = ^them
	case Noisies, Captures:
		filterMask = them
	}

	if pieceType == King {
		GenKingMoves(pos, filterMask, movelist)
		if filter.genQuiets() && !inCheck {
			genCastleMoves(pos, color, movelist)
		}

		return
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

	switch pieceType {
	case Pawn:
		GenPawnMoves(pos, filter, evasionMask, movelist)
	case Knight:
		GenKnightMoves(pos, filterMask, movelist)
	case Bishop:
		GenBishopMoves(pos, filterMask, movelist)
	case Rook:
		GenRookMoves(pos, filterMask, movelist)
	case Queen:
		GenQueenMoves(pos, filterMask, movelist)
	}
}

func GenerateLegalMoves(pos *Position, movelist *Movelist, filter MoveFilter) {
	GenerateLegalMovesForPiece(pos, King, movelist, filter)
	GenerateLegalMovesForPiece(pos, Knight, movelist, filter)
	GenerateLegalMovesForPiece(pos, Bishop, movelist, filter)
	GenerateLegalMovesForPiece(pos, Rook, movelist, filter)
	GenerateLegalMovesForPiece(pos, Queen, movelist, filter)
	GenerateLegalMovesForPiece(pos, Pawn, movelist, filter)
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

	return pawnAttacks(pawns, enemyColor) |
		knightAttacks(knights) |
		bishopAttacks(bishops|queens, occupied) |
		rookAttacks(rooks|queens, occupied) |
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
