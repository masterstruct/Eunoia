package history

import "github.com/masterstruct/Eunoia/internal/board"

const (
	corrSize = 16384

	corrMax      = 1024
	corrMaxBonus = corrMax / 4
)

// CorrHist buckets
const (
	corrPawn = iota
	corrWhite
	corrBlack
	corrMinor

	numCorrBuckets
)

const (
	corrPawnWeight    = 64
	corrNonPawnWeight = 64
	corrMinorWeight   = 64
)

type CorrHist [numCorrBuckets][2][corrSize]int16 // [bucket][stm][hash]

func (h *CorrHist) get(stm board.Color, hash uint64, bucket int) int32 {
	return int32(h[bucket][stm][hash%corrSize])
}

func (h *CorrHist) update(stm board.Color, bucket int, hash uint64, bonus int32) {
	entry := &h[bucket][stm][hash%corrSize]
	gravity(entry, bonus, corrMaxBonus, corrMax)
}

func (h *CorrHist) Clear() {
	clear(h[:])
}

func (h *CorrHist) Correction(pos *board.Position) int32 {
	stm := pos.SideToMove

	corr := h.get(stm, getPawnHash(pos), corrPawn) * corrPawnWeight
	corr += h.get(stm, getNonPawnHash(pos, board.White), corrWhite) * corrNonPawnWeight
	corr += h.get(stm, getNonPawnHash(pos, board.Black), corrBlack) * corrNonPawnWeight
	corr += h.get(stm, getMinorHash(pos), corrMinor) * corrMinorWeight

	return corr / corrMax
}

func (h *CorrHist) Update(pos *board.Position, depth int, score, staticEval int32) {
	stm := pos.SideToMove
	bonus := corrBonus(score-staticEval, depth)

	h.update(stm, corrPawn, getPawnHash(pos), bonus)
	h.update(stm, corrWhite, getNonPawnHash(pos, board.White), bonus)
	h.update(stm, corrBlack, getNonPawnHash(pos, board.Black), bonus)
	h.update(stm, corrMinor, getMinorHash(pos), bonus)
}

func corrBonus(diff int32, depth int) int32 {
	return (diff * int32(depth)) / 8
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
