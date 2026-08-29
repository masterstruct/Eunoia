package search

import (
	"bytes"
	"io"
	"strconv"
	"time"

	"github.com/masterstruct/Eunoia/internal/board"
)

const MaxPly = 128

type PVTable struct {
	length [MaxPly]int
	line   [MaxPly][MaxPly]board.Move
}

func (pv *PVTable) Init(ply int) {
	if ply >= MaxPly {
		return
	}
	pv.length[ply] = ply
}

func (pv *PVTable) Store(ply int, move board.Move) {
	if ply >= MaxPly {
		return
	}
	pv.line[ply][ply] = move

	child := ply + 1
	if child >= MaxPly {
		pv.length[ply] = child
		return
	}
	for next := child; next < pv.length[child]; next++ {
		pv.line[ply][next] = pv.line[child][next]
	}
	pv.length[ply] = pv.length[child]
}

func (pv *PVTable) Line() []board.Move {
	return pv.line[0][:pv.length[0]]
}

func (ss *SearchState) printPV(w io.Writer, depth int, score int16) {
	if ss.Quiet {
		return
	}

	var buf bytes.Buffer

	nodes := ss.Nodes
	elapsed := max(time.Since(ss.StartTime).Milliseconds(), 1)
	nps := 1000 * nodes / uint64(elapsed)

	buf.WriteString("info depth ")
	buf.WriteString(strconv.Itoa(depth))
	if isMateScore(score) {
		buf.WriteString(" score mate ")
		buf.WriteString(strconv.Itoa(mateInMoves(score)))
	} else {
		buf.WriteString(" score cp ")
		buf.WriteString(strconv.Itoa(int(score)))
	}
	buf.WriteString(" nodes ")
	buf.WriteString(strconv.FormatUint(nodes, 10))
	buf.WriteString(" nps ")
	buf.WriteString(strconv.FormatUint(nps, 10))
	buf.WriteString(" time ")
	buf.WriteString(strconv.FormatInt(elapsed, 10))
	buf.WriteString(" pv")

	chess960 := board.IsChess960()
	for _, move := range ss.pv.Line() {
		buf.WriteByte(' ')
		writeMove(&buf, move, chess960)
	}

	buf.WriteByte('\n')
	w.Write(buf.Bytes())
}

func writeSquare(buf *bytes.Buffer, sq board.Square) {
	if sq == board.NoSquare {
		buf.WriteByte('-')
		return
	}
	buf.WriteByte('a' + byte(sq.File()))
	buf.WriteByte('1' + byte(sq.Rank()))
}

func writeMove(buf *bytes.Buffer, move board.Move, chess960 bool) {
	writeSquare(buf, move.From())
	to := move.To()
	if move.IsCastle() && !chess960 {
		to = board.FischerRandomToStandardCastling(move)
	}
	writeSquare(buf, to)
	if move.IsPromo() {
		buf.WriteByte(move.Promo().String())
	}
}

func isMateScore(score int16) bool {
	return score >= MATE-int16(MaxPly) || score <= -MATE+int16(MaxPly)
}

func mateInMoves(score int16) int {
	var plies int16
	if score > 0 {
		plies = MATE - score
	} else {
		plies = MATE + score
	}
	if plies <= 0 {
		plies = 1
	}
	moves, _ := board.PlyToFullmoves(uint16(plies - 1))
	if score < 0 {
		return -int(moves)
	}
	return int(moves)
}
