package search

import (
	"github.com/masterstruct/Eunoia/internal/board"
)

type Node struct {
	// add fields, update MakeMove and UnmakeMove
	staticEval int32
	piece      board.Piece
	move       board.Move
	killer     board.Move
	excluded   board.Move
}

// headroom if accessing stack[ply+N]
type NodeStack [MaxPly + 8]Node

func (stack *NodeStack) reset() {
	for i := range stack {
		stack[i] = Node{
			// set defaults
			staticEval: -INF,
			piece:      board.NoPiece,
			move:       board.NullMove,
			killer:     board.NullMove,
			excluded:   board.NullMove,
		}
	}
}
