package search

import "github.com/masterstruct/Eunoia/internal/board"

// Based on Hobbes implementation
// https://github.com/kelseyde/hobbes-chess-engine/blob/main/src/search/see.rs

const (
	seePawnValue   = 100
	seeKnightValue = 300
	seeBishopValue = 330
	seeRookValue   = 500
	seeQueenValue  = 900

	seeQuietMultiplier = -80
	seeNoisyMultiplier = -100
)

func seeValue(pt board.PieceType) int32 {
	switch pt {
	case board.Pawn:
		return seePawnValue
	case board.Knight:
		return seeKnightValue
	case board.Bishop:
		return seeBishopValue
	case board.Rook:
		return seeRookValue
	case board.Queen:
		return seeQueenValue
	default:
		return 0
	}
}

func see(pos *board.Position, move board.Move, threshold int32) bool {
	from := move.From()
	to := move.To()

	balance := seeMoveValue(pos, move) - threshold

	if balance < 0 {
		return false
	}

	var nextVictim board.PieceType
	if move.IsPromo() {
		nextVictim = move.Promo()
	} else {
		nextVictim = pos.Board[from].Type
	}

	balance -= seeValue(nextVictim)

	if balance >= 0 {
		return true
	}

	occ := pos.Occupied() &^ from.Bit() &^ to.Bit()
	if move.IsEnPassant() {
		epSq := to - (16*board.Square(pos.SideToMove) - 8)
		occ ^= epSq.Bit()
	}

	attackers := seeSquareAttackers(pos, to, occ) & occ
	diagonal := pos.Pieces[board.Bishop] | pos.Pieces[board.Queen]
	orthogonal := pos.Pieces[board.Rook] | pos.Pieces[board.Queen]

	whitePinned := pos.Pinned[board.White]
	blackPinned := pos.Pinned[board.Black]
	pinned := whitePinned | blackPinned
	attackers &= ^pinned |
		(whitePinned & board.Extending(pos.KingSq[board.White], to)) |
		(blackPinned & board.Extending(pos.KingSq[board.Black], to))

	stm := pos.SideToMove.Opponent()

	for {
		ourAttackers := attackers & pos.Colors[stm]
		if ourAttackers == 0 {
			break
		}

		attacker := leastValuableAttacker(pos, ourAttackers)
		if attacker == board.King && (attackers&pos.Colors[stm.Opponent()]) != 0 {
			break
		}

		pieces := pos.Pieces[attacker] & ourAttackers
		sq := (ourAttackers & pieces).LSB()
		occ.ClearBit(sq)
		stm = stm.Opponent()

		balance = -balance - 1 - seeValue(attacker)
		if balance >= 0 {
			break
		}

		if attacker == board.Pawn || attacker == board.Bishop || attacker == board.Queen {
			attackers |= board.BishopAttacks(to, occ) & diagonal
		}
		if attacker == board.Rook || attacker == board.Queen {
			attackers |= board.RookAttacks(to, occ) & orthogonal
		}
		attackers &= occ
	}

	return stm != pos.SideToMove
}

func seeMoveValue(pos *board.Position, move board.Move) int32 {
	value := seeValue(pos.Board[move.To()].Type)
	if move.IsPromo() {
		value += seeValue(move.Promo())
	} else if move.IsEnPassant() {
		value = seeValue(board.Pawn)
	}
	return value
}

func seeSquareAttackers(pos *board.Position, sq board.Square, occ board.Bitboard) board.Bitboard {
	wPawns := board.PawnAttacks[board.Black][sq] & pos.PieceBB(board.WhitePawn)
	bPawns := board.PawnAttacks[board.White][sq] & pos.PieceBB(board.BlackPawn)
	knights := board.KnightAttacks[sq] & pos.Pieces[board.Knight]
	kings := board.KingAttacks[sq] & pos.Pieces[board.King]

	queens := pos.Pieces[board.Queen]
	bishops := pos.Pieces[board.Bishop]
	rooks := pos.Pieces[board.Rook]

	return wPawns | bPawns | knights | kings |
		board.BishopAttacks(sq, occ)&(bishops|queens) |
		board.RookAttacks(sq, occ)&(rooks|queens)
}

func leastValuableAttacker(pos *board.Position, ourAttackers board.Bitboard) board.PieceType {
	if (ourAttackers & pos.Pieces[board.Pawn]) != 0 {
		return board.Pawn
	}
	if (ourAttackers & pos.Pieces[board.Knight]) != 0 {
		return board.Knight
	}
	if (ourAttackers & pos.Pieces[board.Bishop]) != 0 {
		return board.Bishop
	}
	if (ourAttackers & pos.Pieces[board.Rook]) != 0 {
		return board.Rook
	}
	if (ourAttackers & pos.Pieces[board.Queen]) != 0 {
		return board.Queen
	}
	if (ourAttackers & pos.Pieces[board.King]) != 0 {
		return board.King
	}
	panic("No attackers found")
}
