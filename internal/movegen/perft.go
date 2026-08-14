package movegen

import (
	"fmt"
	"time"

	"github.com/masterstruct/Eunoia/internal/board"
)

type PerftResult struct {
	Nodes uint64
	Time  time.Duration
	NPS   uint64
}

func Perft(pos board.Position, depth int) PerftResult {
	start := time.Now()

	nodes := perft(pos, depth)

	elapsed := time.Since(start)
	ns := max(uint64(elapsed.Nanoseconds()), 1)
	nps := nodes * 1_000_000_000 / ns

	return PerftResult{
		Nodes: nodes,
		Time:  elapsed,
		NPS:   nps,
	}
}

func perft(pos board.Position, depth int) uint64 {
	if depth == 0 {
		return 1
	}

	var nodes uint64
	movelist := GeneratePseudolegalMoves(pos)
	for _, move := range movelist {
		mover := pos.SideToMove
		newPos := pos.MakeMove(move)
		if InCheck(newPos, mover) {
			continue
		}
		nodes += perft(newPos, depth-1)
	}
	return nodes
}

func SplitPerft(pos board.Position, depth int) PerftResult {
	start := time.Now()

	var total uint64
	movelist := GeneratePseudolegalMoves(pos)
	for _, move := range movelist {
		mover := pos.SideToMove
		newPos := pos.MakeMove(move)
		if InCheck(newPos, mover) {
			continue
		}
		nodes := perft(newPos, depth-1)
		fmt.Println(move, nodes)
		total += nodes
	}

	elapsed := time.Since(start)
	nps := uint64(0)
	if elapsed > 0 {
		nps = total * 1_000_000_000 / uint64(elapsed.Nanoseconds())
	}

	return PerftResult{
		Nodes: total,
		Time:  elapsed,
		NPS:   nps,
	}
}
