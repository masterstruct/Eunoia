package search

import (
	"time"

	"github.com/masterstruct/Eunoia/internal/tt"
)

type SearchState struct {
	Stop      bool
	Nodes     uint64
	MaxNodes  uint64
	SoftNodes uint64
	StartTime time.Time
	MaxTime   time.Time
	SoftTime  time.Time

	tt *tt.Table
}

func (ss *SearchState) Reset() {
	ss.Stop = false
	ss.Nodes = 0
	ss.MaxNodes = 0
	ss.SoftNodes = 0
	ss.StartTime = time.Now()
	ss.MaxTime = time.Time{}
	ss.SoftTime = time.Time{}
}

func (ss *SearchState) ClearTT() {
	ss.tt.Clear()
}
