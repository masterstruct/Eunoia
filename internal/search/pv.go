package search

import (
	"bytes"
	"io"
	"strconv"
	"time"

	"github.com/masterstruct/Eunoia/internal/board"
	"github.com/masterstruct/Eunoia/internal/tt"
)

const MaxPly uint16 = 128

type PVTable struct {
	length [MaxPly]uint16
	line   [MaxPly][MaxPly]board.Move
}

func NewPVTable() *PVTable {
	return &PVTable{}
}

func (pv *PVTable) Init(ply uint16) {
	if ply >= MaxPly {
		return
	}
	pv.length[ply] = ply
}

func (pv *PVTable) Store(ply uint16, move board.Move) {
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

func (ss *SearchState) printPV(w io.Writer, score int32, bound tt.Flag) {
	if ss.Quiet {
		return
	}

	pv := ss.pv.Line()
	if len(pv) == 0 {
		return
	}

	var buf bytes.Buffer

	nodes := ss.Nodes
	depth := ss.depth

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

	switch bound {
	case tt.Lower:
		buf.WriteString(" lowerbound")
	case tt.Upper:
		buf.WriteString(" upperbound")
	}

	buf.WriteString(" nodes ")
	buf.WriteString(strconv.FormatUint(nodes, 10))
	buf.WriteString(" nps ")
	buf.WriteString(strconv.FormatUint(nps, 10))
	buf.WriteString(" hashfull ")
	buf.WriteString(strconv.FormatUint(ss.tt.Hashfull(), 10))
	buf.WriteString(" time ")
	buf.WriteString(strconv.FormatInt(elapsed, 10))
	buf.WriteString(" pv")

	chess960 := board.IsChess960()
	for _, move := range pv {
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
