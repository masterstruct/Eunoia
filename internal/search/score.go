package search

import "github.com/masterstruct/Eunoia/internal/board"

const (
	MATE int16 = 30000
	INF  int16 = 32000
)

func isMateScore(score int16) bool {
	return isMating(score) || isMated(score)
}

func isMating(score int16) bool {
	return score >= MATE-int16(MaxPly)
}

func isMated(score int16) bool {
	return score <= -MATE+int16(MaxPly)
}

func scoreToTT(score int16, ply uint16) int16 {
	if score >= MATE-int16(MaxPly) {
		return score + int16(ply)
	}
	if score <= -MATE+int16(MaxPly) {
		return score - int16(ply)
	}
	return score
}

func scoreFromTT(score int16, ply uint16) int16 {
	if score >= MATE-int16(MaxPly) {
		return score - int16(ply)
	}
	if score <= -MATE+int16(MaxPly) {
		return score + int16(ply)
	}
	return score
}

func mateInMoves(score int16) int {
	moves, _ := board.PlyToFullmoves(mateInPlies(score) - 1)
	if score < 0 {
		return -int(moves)
	}
	return int(moves)
}

func mateInPlies(score int16) uint16 {
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
