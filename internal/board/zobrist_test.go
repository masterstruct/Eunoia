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
