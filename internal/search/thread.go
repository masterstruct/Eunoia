package search

import (
	"github.com/masterstruct/Eunoia/internal/board"
	"github.com/masterstruct/Eunoia/internal/tt"
)

type SearchState struct {
	Quiet bool // avoid printing output

	TimeManager

	bestMove  board.Move
	bestScore int32
	depth     int
	selDepth  uint16
	nmpMinPly uint16

	tt *tt.Table
	pv *PVTable

	Keys  Keys // history of position hashes for 3fold detection
	stack NodeStack

	butterflyHistory *[2][64][64]int32
	corrHist         *CorrHist
	killers          *[MaxPly + 1]board.Move
}

func (ss *SearchState) Init(ttSizeMiB uint) {
	ss.depth = 1
	ss.tt = tt.NewTable(ttSizeMiB)
	ss.pv = NewPVTable()
	ss.Keys = Keys{keys: make([]uint64, 0, MaxPly)}
	ss.stack = newNodeStack()
	ss.butterflyHistory = &[2][64][64]int32{}
	ss.corrHist = &CorrHist{}
	ss.killers = &[MaxPly + 1]board.Move{}
}

func (ss *SearchState) PrepareForSearch() {
	ss.Quiet = false
	ss.bestMove = board.NullMove
	ss.bestScore = 0
	ss.depth = 1
	ss.resetTimeManager()

	clear(ss.killers[:])
}

func (ss *SearchState) ClearTables() {
	ss.ClearTT()
	*ss.pv = PVTable{}
	clear(ss.butterflyHistory[:])
	ss.corrHist.Clear()
	clear(ss.killers[:])
}

// clears ONLY the TT - intended for `setoption name Clear Hash`.
// To clear everything use ss.ClearTables()
func (ss *SearchState) ClearTT() {
	ss.tt.Clear()
}

func (ss *SearchState) ResizeTT(sizeMiB uint) {
	ss.tt.Resize(sizeMiB)
}
