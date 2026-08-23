package board

import "testing"

func TestPRNGNext(t *testing.T) {
	p := &PRNG{state: 1234567}

	want := []uint64{
		6457827717110365317,
		3203168211198807973,
		9817491932198370423,
		4593380528125082431,
		16408922859458223821,
	}

	for i, w := range want {
		got := p.Next()
		if got != w {
			t.Fatalf("call %d: got %016x, want %016x", i+1, got, w)
		}
	}
}
