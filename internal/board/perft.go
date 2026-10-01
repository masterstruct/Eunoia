package board

import (
	"fmt"
	"time"
)

type PerftResult struct {
	Nodes uint64
	Time  time.Duration
	NPS   uint64
}

func Perft(pos *Position, depth int) PerftResult {
	start := time.Now()

	nodes := perft(*pos, depth)

	elapsed := time.Since(start)
	ns := max(uint64(elapsed.Nanoseconds()), 1)

	return PerftResult{
		Nodes: nodes,
		Time:  elapsed,
		NPS:   nodes * 1_000_000_000 / ns,
	}
}

func perft(pos Position, depth int) uint64 {
	var movelist Movelist
	GenerateLegalMoves(&pos, &movelist, All)

	if depth <= 1 {
		return uint64(movelist.Len)
	}

	var nodes uint64
	for i := 0; i < movelist.Len; i++ {
		newPos := pos.MakeMove(movelist.Moves[i])
		nodes += perft(newPos, depth-1)
	}

	return nodes
}

func SplitPerft(pos *Position, depth int) PerftResult {
	start := time.Now()

	var total uint64
	var movelist Movelist
	GenerateLegalMoves(pos, &movelist, All)

	for i := 0; i < movelist.Len; i++ {
		newPos := pos.MakeMove(movelist.Moves[i])
		nodes := perft(newPos, depth-1)
		fmt.Println(movelist.Moves[i], nodes)
		total += nodes
	}

	elapsed := time.Since(start)
	ns := max(uint64(elapsed.Nanoseconds()), 1)

	return PerftResult{
		Nodes: total,
		Time:  elapsed,
		NPS:   total * 1_000_000_000 / ns,
	}
}
