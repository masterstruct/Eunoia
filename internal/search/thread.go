package search

type SearchState struct {
	Stop     bool
	Nodes    uint64
	MaxNodes uint64
}

func (ss *SearchState) Reset() {
	ss.Stop = false
	ss.Nodes = 0
	ss.MaxNodes = 0
}
