package search

import (
	"github.com/masterstruct/Eunoia/internal/board"
	"github.com/masterstruct/Eunoia/internal/history"
)

type Node struct {
	// add fields, update MakeMove and UnmakeMove
	staticEval int32
	killer     board.Move
	excluded   board.Move
	cont       *history.ContRow
}

// headroom if accessing stack[ply+N]
type NodeStack [MaxPly + 8]Node

func (stack *NodeStack) reset() {
	for i := range stack {
		stack[i] = Node{
			// set defaults
			staticEval: -INF,
			killer:     board.NullMove,
			excluded:   board.NullMove,
		}
	}
}
