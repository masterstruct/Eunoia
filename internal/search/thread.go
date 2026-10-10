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

	rootPV         [MaxPly + 1]board.Move
	rootPVLen      int
	completedDepth int

	tt *tt.Table
	pv *PVTable

	keyHistory    []uint64 // history of position hashes for 3fold detection
	keyHistoryLen int

	butterflyHistory *[2][64][64]int32
}

func (ss *SearchState) Init(ttSizeMiB uint) {
	ss.depth = 1
	ss.tt = tt.NewTable(ttSizeMiB)
	ss.pv = NewPVTable()
	ss.butterflyHistory = &[2][64][64]int32{}
}

func (ss *SearchState) PrepareForSearch() {
	ss.Quiet = false
	ss.bestMove = board.NullMove
	ss.bestScore = 0
	ss.depth = 1
	ss.completedDepth = 0
	ss.rootPVLen = 0
	ss.resetTimeManager()
}

func (ss *SearchState) saveRootPV() {
	ss.rootPVLen = copy(ss.rootPV[:], ss.pv.Line())
}

func (ss *SearchState) ClearTables() {
	ss.ClearTT()
	*ss.pv = PVTable{}
	clear(ss.butterflyHistory[:])
}

// clears ONLY the TT - intended for `setoption name Clear Hash`.
// To clear everything use ss.ClearTables()
func (ss *SearchState) ClearTT() {
	ss.tt.Clear()
}

func (ss *SearchState) ResizeTT(sizeMiB uint) {
	ss.tt.Resize(sizeMiB)
}

func (ss *SearchState) SetHistory(gameHistory []uint64) {
	ss.keyHistory = append(ss.keyHistory[:0], gameHistory...)
	ss.keyHistoryLen = len(ss.keyHistory)
}
