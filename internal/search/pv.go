package search

import (
	"fmt"
	"io"
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
	line := ss.pv.Line()

	nodes := ss.Nodes
	elapsed := max(time.Since(ss.StartTime).Milliseconds(), 1)
	nps := 1000 * nodes / uint64(elapsed)

	s := fmt.Sprint(line)
	s = s[1 : len(s)-1]

	fmt.Fprintf(w, "info depth %d score cp %d nodes %d nps %d time %d pv %s\n",
		depth, score, nodes, nps, elapsed, s)
}
