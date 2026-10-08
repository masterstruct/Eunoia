package search

import "github.com/masterstruct/Eunoia/internal/board"

const (
	MATE int32 = 30000
	INF  int32 = 32000
)

func isMateScore(score int32) bool {
	return isMating(score) || isMated(score)
}

func isMating(score int32) bool {
	return score >= MATE-int32(MaxPly)
}

func isMated(score int32) bool {
	return score <= -MATE+int32(MaxPly)
}

func scoreToTT(score int32, ply uint16) int16 {
	if isMating(score) {
		return int16(score) + int16(ply)
	}
	if isMated(score) {
		return int16(score) - int16(ply)
	}
	return int16(score)
}

func scoreFromTT(score int16, ply uint16) int32 {
	s := int32(score)
	if isMating(s) {
		return s - int32(ply)
	}
	if isMated(s) {
		return s + int32(ply)
	}
	return s
}

func mateInMoves(score int32) int {
	moves, _ := board.PlyToFullmoves(mateInPlies(score) - 1)
	if score < 0 {
		return -int(moves)
	}
	return int(moves)
}

func mateInPlies(score int32) uint16 {
	plies := MATE
	if score > 0 {
		plies -= score
	} else {
		plies += score
	}
	if plies <= 0 {
		plies = 1
	}
	return uint16(plies)
}
