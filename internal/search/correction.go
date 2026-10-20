package search

import "github.com/masterstruct/Eunoia/internal/board"

const (
	corrHistSize = 16384
	corrHistMask = corrHistSize - 1

	corrHistGrain    = 16 // entry units per centipawn
	corrMaxCP        = 64
	corrHistMax      = corrMaxCP * corrHistGrain
	corrHistMaxBonus = corrHistMax / 4

	corrWeightScale = 256
)

// CorrHist buckets
const (
	corrPawn = iota
	corrWhite
	corrBlack
	corrMinor

	numCorrBuckets
)

// bonus = diff * depth * mult / div
const (
	corrPawnBonusMult = 16
	corrPawnBonusDiv  = 128
	corrPawnWeight    = 256

	corrNonPawnBonusMult = 16
	corrNonPawnBonusDiv  = 128
	corrNonPawnWeight    = 256

	corrMinorBonusMult = 16
	corrMinorBonusDiv  = 128
	corrMinorWeight    = 256
)

type CorrHist struct {
	entries [numCorrBuckets][2][corrHistSize]int16 // [bucket][stm][key]
}

func (ch *CorrHist) get(stm board.Color, key uint64, bucket int) int32 {
	return int32(ch.entries[bucket][stm][key&corrHistMask])
}

func (ch *CorrHist) update(stm board.Color, key uint64, bucket int, bonus int32) {
	entry := &ch.entries[bucket][stm][key&corrHistMask]
	*entry = int16(gravity(int32(*entry), bonus, corrHistMax))
}

func (ch *CorrHist) Clear() {
	clear(ch.entries[:])
}

func (ch *CorrHist) Correction(pos *board.Position) int32 {
	stm := pos.SideToMove

	sum := ch.get(stm, getPawnKey(pos), corrPawn) * corrPawnWeight
	sum += ch.get(stm, getNonPawnKey(pos, board.White), corrWhite) * corrNonPawnWeight
	sum += ch.get(stm, getNonPawnKey(pos, board.Black), corrBlack) * corrNonPawnWeight
	sum += ch.get(stm, getMinorKey(pos), corrMinor) * corrMinorWeight

	cp := sum / (corrHistGrain * corrWeightScale)
	return min(max(cp, -corrMaxCP), corrMaxCP)
}

func (ch *CorrHist) Update(pos *board.Position, depth int, staticEval, bestScore int32) {
	stm := pos.SideToMove
	diff := bestScore - staticEval

	ch.update(stm, getPawnKey(pos), corrPawn, corrBonus(diff, depth, corrPawnBonusMult, corrPawnBonusDiv))
	ch.update(stm, getNonPawnKey(pos, board.White), corrWhite, corrBonus(diff, depth, corrNonPawnBonusMult, corrNonPawnBonusDiv))
	ch.update(stm, getNonPawnKey(pos, board.Black), corrBlack, corrBonus(diff, depth, corrNonPawnBonusMult, corrNonPawnBonusDiv))
	ch.update(stm, getMinorKey(pos), corrMinor, corrBonus(diff, depth, corrMinorBonusMult, corrMinorBonusDiv))
}

func corrBonus(diff int32, depth int, mult, div int32) int32 {
	bonus := diff * int32(depth) * mult / div
	return min(max(bonus, -corrHistMaxBonus), corrHistMaxBonus)
}

// TODO: update incrementally in makemove
func getPawnKey(pos *board.Position) uint64 {
	var key uint64
	wp := pos.PieceBB(board.WhitePawn)
	bp := pos.PieceBB(board.BlackPawn)

	for wp != 0 {
		key ^= board.ZobristTable.PieceKey(board.White, board.Pawn, wp.PopLSB())
	}
	for bp != 0 {
		key ^= board.ZobristTable.PieceKey(board.Black, board.Pawn, bp.PopLSB())
	}
	return key
}

func getNonPawnKey(pos *board.Position, stm board.Color) uint64 {
	var key uint64
	pieces := pos.Colors[stm] &^ pos.Pieces[board.Pawn]
	for pieces != 0 {
		sq := pieces.PopLSB()
		key ^= board.ZobristTable.PieceKey(stm, pos.Board[sq].Type, sq)
	}
	return key
}

func getMinorKey(pos *board.Position) uint64 {
	var key uint64
	pieces := pos.Pieces[board.Knight] | pos.Pieces[board.Bishop] | pos.Pieces[board.King]
	for pieces != 0 {
		sq := pieces.PopLSB()
		p := pos.Board[sq]
		key ^= board.ZobristTable.PieceKey(p.Color, p.Type, sq)
	}
	return key
}
