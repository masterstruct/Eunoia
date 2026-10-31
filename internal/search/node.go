package search

import "github.com/masterstruct/Eunoia/internal/board"

type Node struct {
	move     board.Move
	piece    board.Piece
	captured board.Piece
}

type NodeStack struct {
	data [MaxPly + 8]Node
}

func newNodeStack() NodeStack {
	var stack NodeStack
	for i := range stack.data {
		stack.data[i] = Node{
			move:     board.NullMove,
			piece:    board.NoPiece,
			captured: board.NoPiece,
		}
	}
	return stack
}
