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
