package search

import "time"

type SearchState struct {
	Stop      bool
	Nodes     uint64
	MaxNodes  uint64
	StartTime time.Time
	MaxTime   time.Time
}

func (ss *SearchState) Reset() {
	ss.Stop = false
	ss.Nodes = 0
	ss.MaxNodes = 0
	ss.StartTime = time.Now()
	ss.MaxTime = time.Time{}
}
