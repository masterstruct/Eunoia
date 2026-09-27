package board

import (
	"testing"
)

func TestIncrementalZobrist(t *testing.T) {
	pos, err := ParseFEN(KiwipeteFEN)
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

func walkAndVerifyHash(t *testing.T, pos Position, depth int, checked *int) {
	t.Helper()
	if depth <= 0 {
		return
	}

	var movelist Movelist
	GenerateLegalMoves(&pos, &movelist, All)

	for i := 0; i < movelist.Len; i++ {
		newPos := pos.MakeMove(movelist.Moves[i])

		want := ZobristTable.ComputeHash(&newPos)
		if newPos.Hash != want {
			t.Errorf("move %v depth %d: incremental hash %d, computed hash %d", movelist.Moves[i], depth, newPos.Hash, want)
		}
		*checked++

		walkAndVerifyHash(t, newPos, depth-1, checked)
	}
}
