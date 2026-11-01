package search

type Node struct {
	// add fields, update MakeMove and UnmakeMove
	staticEval int32
}

// headroom if accessing stack[ply+N]
type NodeStack [MaxPly + 8]Node

func newNodeStack() NodeStack {
	var stack NodeStack
	for i := range stack {
		stack[i] = Node{
			// set defaults
			staticEval: -INF,
		}
	}
	return stack
}
