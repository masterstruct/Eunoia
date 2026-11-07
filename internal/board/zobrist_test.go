package board

import "testing"

func TestPRNGNext(t *testing.T) {
	p := NewPRNG(zobristSeed)

	want := []uint64{
		11504272154707947787,
		11487468225442647383,
		2122327086578615790,
		741898026479467641,
		10316960201484082974,
	}

	for i, w := range want {
		got := p.Next()
		if got != w {
			t.Fatalf("call %d: got %d, want %d", i+1, got, w)
		}
	}
}

func TestPRNGNext_Reproducible(t *testing.T) {
	p1 := NewPRNG(zobristSeed)
	p2 := NewPRNG(zobristSeed)

	for i := range 5 {
		got1 := p1.Next()
		got2 := p2.Next()
		if got1 != got2 {
			t.Fatalf("call %d: sequences don't match: %d vs %d", i+1, got1, got2)
		}
	}
}

func TestComputeHash(t *testing.T) {
	pos := StartingPosition()
	want := uint64(7484279522850932296)
	got := ZobristTable.ComputeHash(&pos)

	if got != want {
		t.Fatalf("got %d want %d", got, want)
	}
}

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
		newPos := pos
		newPos.Make(movelist.Moves[i].Move)

		want := ZobristTable.ComputeHash(&newPos)
		if newPos.Hash != want {
			t.Errorf("move %v depth %d: incremental hash %d, computed hash %d", movelist.Moves[i], depth, newPos.Hash, want)
		}
		*checked++

		walkAndVerifyHash(t, newPos, depth-1, checked)
	}
}
