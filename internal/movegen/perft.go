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

func Perft(pos *board.Position, depth int) PerftResult {
	start := time.Now()

	nodes := perft(*pos, depth)

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
	var movelist Movelist
	GeneratePseudolegalMoves(&pos, &movelist)
	mover := pos.SideToMove
	for i := 0; i < movelist.Len; i++ {
		newPos := pos.MakeMove(movelist.Moves[i])
		if InCheck(&newPos, mover) {
			continue
		}
		nodes += perft(newPos, depth-1)
	}

	return nodes
}

func SplitPerft(pos *board.Position, depth int) PerftResult {
	start := time.Now()

	var total uint64
	var movelist Movelist
	GeneratePseudolegalMoves(pos, &movelist)

	mover := pos.SideToMove
	for i := 0; i < movelist.Len; i++ {
		newPos := pos.MakeMove(movelist.Moves[i])
		if InCheck(&newPos, mover) {
			continue
		}
		nodes := perft(newPos, depth-1)
		fmt.Println(movelist.Moves[i], nodes)
		total += nodes
	}

	elapsed := time.Since(start)
	ns := uint64(elapsed.Nanoseconds())
	if ns == 0 {
		ns = 1
	}

	return PerftResult{
		Nodes: total,
		Time:  elapsed,
		NPS:   total * 1_000_000_000 / ns,
	}
}
