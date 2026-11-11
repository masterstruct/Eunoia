package history

import "github.com/masterstruct/Eunoia/internal/board"

const (
	corrSize = 16384

	corrGrain    = 16 // entry units per centipawn
	corrMaxCP    = 64
	corrMax      = corrMaxCP * corrGrain
	corrMaxBonus = corrMax / 4

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

type CorrHist [numCorrBuckets][2][corrSize]int16 // [bucket][stm][hash]

func (h *CorrHist) get(stm board.Color, hash uint64, bucket int) int32 {
	return int32(h[bucket][stm][hash%corrSize])
}

func (h *CorrHist) update(stm board.Color, hash uint64, bucket int, bonus int32) {
	entry := &h[bucket][stm][hash%corrSize]
	gravity(entry, bonus, corrMaxBonus, corrMax)
}

func (h *CorrHist) Clear() {
	clear(h[:])
}

func (h *CorrHist) Correction(pos *board.Position) int32 {
	stm := pos.SideToMove

	sum := h.get(stm, getPawnHash(pos), corrPawn) * corrPawnWeight
	sum += h.get(stm, getNonPawnHash(pos, board.White), corrWhite) * corrNonPawnWeight
	sum += h.get(stm, getNonPawnHash(pos, board.Black), corrBlack) * corrNonPawnWeight
	sum += h.get(stm, getMinorHash(pos), corrMinor) * corrMinorWeight

	cp := sum / (corrGrain * corrWeightScale)
	return min(max(cp, -corrMaxCP), corrMaxCP)
}

func (h *CorrHist) Update(pos *board.Position, depth int, score, staticEval int32) {
	stm := pos.SideToMove
	diff := score - staticEval

	h.update(stm, getPawnHash(pos), corrPawn, corrBonus(diff, depth, corrPawnBonusMult, corrPawnBonusDiv))
	h.update(stm, getNonPawnHash(pos, board.White), corrWhite, corrBonus(diff, depth, corrNonPawnBonusMult, corrNonPawnBonusDiv))
	h.update(stm, getNonPawnHash(pos, board.Black), corrBlack, corrBonus(diff, depth, corrNonPawnBonusMult, corrNonPawnBonusDiv))
	h.update(stm, getMinorHash(pos), corrMinor, corrBonus(diff, depth, corrMinorBonusMult, corrMinorBonusDiv))
}

func corrBonus(diff int32, depth int, mult, div int32) int32 {
	bonus := diff * int32(depth) * mult / div
	return min(max(bonus, -corrMaxBonus), corrMaxBonus)
}

// TODO: update incrementally in makemove
func getPawnHash(pos *board.Position) uint64 {
	var key uint64
	wp := pos.Pawns(board.White)
	bp := pos.Pawns(board.Black)

	for !wp.IsEmpty() {
		key ^= board.ZobristTable.PieceKey(board.WhitePawn, wp.PopLSB())
	}
	for !bp.IsEmpty() {
		key ^= board.ZobristTable.PieceKey(board.BlackPawn, bp.PopLSB())
	}
	return key
}

func getNonPawnHash(pos *board.Position, stm board.Color) uint64 {
	var key uint64
	pieces := pos.Colors[stm] &^ pos.Pieces[board.Pawn]
	for !pieces.IsEmpty() {
		sq := pieces.PopLSB()
		key ^= board.ZobristTable.PieceKey(pos.PieceOn(sq), sq)
	}
	return key
}

func getMinorHash(pos *board.Position) uint64 {
	var key uint64
	pieces := pos.Pieces[board.Knight] | pos.Pieces[board.Bishop] | pos.Pieces[board.King]
	for !pieces.IsEmpty() {
		sq := pieces.PopLSB()
		key ^= board.ZobristTable.PieceKey(pos.PieceOn(sq), sq)
	}
	return key
}
