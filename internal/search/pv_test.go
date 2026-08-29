package search

import (
	"io"
	"testing"
	"time"

	"github.com/masterstruct/Eunoia/internal/board"
)

func BenchmarkPrintPV(b *testing.B) {
	moves := []board.Move{
		board.NewMove(board.E2, board.E4),
		board.NewMove(board.E7, board.E5),
		board.NewMove(board.G1, board.F3),
		board.NewMove(board.B8, board.C6),
		board.NewMove(board.B1, board.C3),
		board.NewMove(board.G8, board.F6),
		board.NewMove(board.D2, board.D4),
		board.NewCapture(board.E5, board.D4),
		board.NewMove(board.F3, board.D4),
	}

	var pv PVTable
	copy(pv.line[0][:], moves)
	pv.length[0] = len(moves)

	ss := &SearchState{pv: &pv}
	ss.StartTime = time.Now().Add(-5 * time.Second)
	ss.Nodes = 25461085

	for b.Loop() {
		ss.printPV(io.Discard, 8, 38, board.White)
	}
}
