package search

import "github.com/masterstruct/Eunoia/internal/board"

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
