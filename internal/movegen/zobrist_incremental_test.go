package movegen

import (
	"testing"

	"github.com/masterstruct/Eunoia/internal/board"
)

func TestIncrementalZobrist(t *testing.T) {
	pos, err := board.ParseFEN("r3k2r/p1ppqpb1/bn2pnp1/3PN3/1p2P3/2N2Q1p/PPPBBPPP/R3K2R w KQkq -")
	if err != nil {
		t.Fatalf("bad FEN: %v", err)
	}

	checked := 0
	walkAndVerifyHash(t, pos, 4, &checked)

	if checked == 0 {
		t.Fatal("no positions checked")
	}
	t.Logf("checked %d positions", checked)
}

func walkAndVerifyHash(t *testing.T, pos board.Position, depth int, checked *int) {
	t.Helper()
	if depth <= 0 {
		return
	}

	var movelist Movelist
	GeneratePseudolegalMoves(&pos, &movelist)
	mover := pos.SideToMove

	for i := 0; i < movelist.Len; i++ {
		newPos := pos.MakeMove(movelist.Moves[i])
		if InCheck(&newPos, mover) {
			continue
		}

		want := board.ZobristTable.ComputeHash(&newPos)
		if newPos.Hash != want {
			t.Errorf("move %v depth %d: incremental hash %d, computed hash %d", movelist.Moves[i], depth, newPos.Hash, want)
		}
		*checked++

		walkAndVerifyHash(t, newPos, depth-1, checked)
	}
}
