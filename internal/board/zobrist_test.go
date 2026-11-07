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
