package board

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

	occupied := pos.Occupied()
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

// returns true if the side that just moved left its own king attacked.
// Call ONLY on a position returned by MakeMove to filter pseudolegal moves
func (pos *Position) IsIllegal() bool {
	return IsSquareAttacked(pos, pos.KingSq[pos.SideToMove.Opponent()], pos.SideToMove)
}

func GenKnightMoves(pos *Position, movelist *Movelist) {
	color := pos.SideToMove
	opponents := pos.Colors[color.Opponent()]
	occupied := pos.Occupied()

	// loop over each knight
	knights := pos.PieceBB(Piece{Type: Knight, Color: color})
	for knights != 0 {
		from := knights.PopLSB()
		attacks := KnightAttacks[from]
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
}

func GenBishopMoves(pos *Position, movelist *Movelist) {
	color := pos.SideToMove
	opponents := pos.Colors[color.Opponent()]
	occupied := pos.Occupied()

	// loop over each bishop
	bishops := pos.PieceBB(Piece{Type: Bishop, Color: color})
	for bishops != 0 {
		from := bishops.PopLSB()
		attacks := BishopAttacks(from, occupied)
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
}

func GenRookMoves(pos *Position, movelist *Movelist) {
	color := pos.SideToMove
	opponents := pos.Colors[color.Opponent()]
	occupied := pos.Occupied()

	// loop over each rook
	rooks := pos.PieceBB(Piece{Type: Rook, Color: color})
	for rooks != 0 {
		from := rooks.PopLSB()
		attacks := RookAttacks(from, occupied)
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
}

func GenQueenMoves(pos *Position, movelist *Movelist) {
	color := pos.SideToMove
	opponents := pos.Colors[color.Opponent()]
	occupied := pos.Occupied()

	// loop over each queen
	queens := pos.PieceBB(Piece{Type: Queen, Color: color})
	for queens != 0 {
		from := queens.PopLSB()
		attacks := QueenAttacks(from, occupied)
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
}

func GenPawnMoves(pos *Position, movelist *Movelist) {
	color := pos.SideToMove
	enemyColor := color.Opponent()
	occupied := pos.Occupied()
	enemyOcc := pos.Colors[enemyColor]

	pawns := pos.PieceBB(Piece{Type: Pawn, Color: color})

	// en passant
	epSq := pos.EnPassant
	if epSq != NoSquare {
		attackers := PawnAttacks[enemyColor][epSq] & pawns
		for attackers != 0 {
			movelist.Add(NewEnPassant(attackers.PopLSB(), epSq))
		}
	}

	for pawns != 0 {
		from := pawns.PopLSB()

		// captures
		captures := PawnAttacks[color][from] & enemyOcc
		for captures != 0 {
			to := captures.PopLSB()
			if to.Rank() == Rank1 || to.Rank() == Rank8 {
				movelist.Add(NewCapturePromo(from, to, Knight))
				movelist.Add(NewCapturePromo(from, to, Bishop))
				movelist.Add(NewCapturePromo(from, to, Rook))
				movelist.Add(NewCapturePromo(from, to, Queen))
			} else {
				movelist.Add(NewCapture(from, to))
			}
		}

		// pushes
		if color == White {
			to := from + 8
			if occupied.IsBitSet(to) {
				continue
			}

			if to.Rank() == Rank8 {
				movelist.Add(NewPromo(from, to, Knight))
				movelist.Add(NewPromo(from, to, Bishop))
				movelist.Add(NewPromo(from, to, Rook))
				movelist.Add(NewPromo(from, to, Queen))
				continue
			}

			movelist.Add(NewMove(from, to))

			if from.Rank() == Rank2 {
				to2 := from + 16
				if !occupied.IsBitSet(to2) {
					movelist.Add(NewDoublePush(from, to2))
				}
			}
		} else {
			to := from - 8
			if occupied.IsBitSet(to) {
				continue
			}

			if to.Rank() == Rank1 {
				movelist.Add(NewPromo(from, to, Knight))
				movelist.Add(NewPromo(from, to, Bishop))
				movelist.Add(NewPromo(from, to, Rook))
				movelist.Add(NewPromo(from, to, Queen))
				continue
			}

			movelist.Add(NewMove(from, to))

			if from.Rank() == Rank7 {
				to2 := from - 16
				if !occupied.IsBitSet(to2) {
					movelist.Add(NewDoublePush(from, to2))
				}
			}
		}
	}
}

func GenKingMoves(pos *Position, movelist *Movelist) {
	color := pos.SideToMove

	from := pos.KingSq[color]

	// temporarily reversed iteration order: kingside -> queenside

	// base := CastlingRights(color * 2)
	// for right := base; right < base+2; right++ {
	// 	rookSq := pos.Castling[right]
	// 	if pos.Castling.Has(right) && canCastle(pos, from, rookSq) {
	// 		movelist.Add(NewCastle(from, rookSq))
	// 	}
	// }

	base := int(color * 2)
	for right := base + 1; right >= base; right-- {
		rookSq := pos.Castling[right]
		if pos.Castling.Has(CastlingRights(right)) && canCastle(pos, from, rookSq) {
			movelist.Add(NewCastle(from, rookSq))
		}
	}

	attacks := KingAttacks[from]
	captures := attacks & pos.Colors[color.Opponent()]
	quiets := attacks &^ pos.Occupied()

	for captures != 0 {
		// capture
		movelist.Add(NewCapture(from, captures.PopLSB()))
	}

	for quiets != 0 {
		// quiet move
		movelist.Add(NewMove(from, quiets.PopLSB()))
	}
}

func canCastle(pos *Position, kingSq, rookSq Square) bool {
	// TODO: on ParseFEN(), precompute which bits to check
	// with occupied and IsSquareAttacked.
	// Compare (rookPath|kingPath)&occupied==0,
	// then just check if kingPath is attacked.

	rank := kingSq.Rank()
	kingFile := kingSq.File()
	rookFile := rookSq.File()
	occupied := pos.Occupied()
	occupied.ClearBit(kingSq)
	occupied.ClearBit(rookSq) // in chess960 king can jump over rook

	var toFileKing File
	var toFileRook File
	var kingDir File
	var rookDir File
	if rookSq > kingSq {
		// kingside
		kingDir = 1
		rookDir = -1
		toFileKing = FileG
		toFileRook = FileF
		if toFileRook > rookFile {
			// castling kingside (rook moves left), but because chess960,
			// rook can start on H1 and actually move left.
			rookDir = 1
		}
	} else {
		// queenside
		kingDir = -1
		rookDir = -1
		toFileKing = FileC
		toFileRook = FileD
		if toFileKing > kingFile {
			// castling queenside (left), but because chess960,
			// king can start on B1 and actually move right.
			// this cannot happen with kingside if position is legal.
			kingDir = 1
		}
		if toFileRook > rookFile {
			rookDir = 1
		}
	}

	for file := kingFile; file != toFileKing+kingDir; file += kingDir {
		sq := NewSquare(file, rank)
		if (occupied | pos.Threats).IsBitSet(sq) {
			return false
		}
	}
	for file := rookFile; file != toFileRook+rookDir; file += rookDir {
		sq := NewSquare(file, rank)
		if occupied.IsBitSet(sq) {
			return false
		}
	}
	return true
}

func GeneratePseudolegalMoves(pos *Position, movelist *Movelist) {
	GenKnightMoves(pos, movelist)
	GenBishopMoves(pos, movelist)
	GenRookMoves(pos, movelist)
	GenQueenMoves(pos, movelist)
	GenPawnMoves(pos, movelist)
	GenKingMoves(pos, movelist)
}

// squares attacked by opponent pieces
func (pos *Position) calculateThreats(myColor Color) Bitboard {
	occupied := pos.Occupied() &^ SquareBB[pos.KingSq[myColor]]
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
