package search

import (
	"time"

	"github.com/masterstruct/Eunoia/internal/tt"
)

type SearchState struct {
	Quiet bool // avoid printing output

	Stop      bool
	Nodes     uint64
	MaxNodes  uint64
	SoftNodes uint64
	StartTime time.Time
	MaxTime   time.Time
	SoftTime  time.Time

	tt *tt.Table
	pv *PVTable

	keyHistory    []uint64 // history of position hashes for 3fold detection
	keyHistoryLen int

	butterflyHistory *[2][64][64]int
}

func (ss *SearchState) Init(ttSizeMiB uint) {
	ss.tt = tt.NewTable(ttSizeMiB)
	ss.pv = NewPVTable()
	ss.butterflyHistory = &[2][64][64]int{}
}

func (ss *SearchState) PrepareForSearch() {
	ss.Quiet = false
	ss.Stop = false
	ss.Nodes = 0
	ss.MaxNodes = 0
	ss.SoftNodes = 0
	ss.StartTime = time.Now()
	ss.MaxTime = time.Time{}
	ss.SoftTime = time.Time{}
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
