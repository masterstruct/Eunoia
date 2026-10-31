package search

type Node struct {
	// TODO: add fields, update MakeMove and UnmakeMove
}

type NodeStack struct {
	// headroom if accessing stack[ply+N]
	data [MaxPly + 8]Node
}

func newNodeStack() NodeStack {
	var stack NodeStack
	for i := range stack.data {
		stack.data[i] = Node{
			// TODO: set defaults
		}
	}
	return stack
}
