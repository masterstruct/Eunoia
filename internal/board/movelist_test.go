package board

import "testing"

func TestMovelistSwap(t *testing.T) {
	var movelist Movelist
	first := NewMove(A1, A2).ScoredMove(10)
	second := NewMove(B1, B2).ScoredMove(20)
	third := NewMove(C1, C2).ScoredMove(30)
	movelist.AddScoredMove(first)
	movelist.AddScoredMove(second)
	movelist.AddScoredMove(third)

	movelist.Swap(0, 2)

	want := [3]ScoredMove{third, second, first}
	for i, move := range want {
		got := movelist.Moves[i]
		if got != move {
			t.Errorf("index %d got %v, want %v", i, got, move)
		}
	}
	if movelist.Len != 3 {
		t.Errorf("got len %d, want 3", movelist.Len)
	}
}

func TestMovelistRemove(t *testing.T) {
	moves := []ScoredMove{
		NewMove(A1, A2).ScoredMove(10),
		NewMove(B1, B2).ScoredMove(20),
		NewMove(C1, C2).ScoredMove(30),
		NewMove(D1, D2).ScoredMove(40),
		NewMove(E1, E2).ScoredMove(50),
	}
	tests := []struct {
		index int
		want  []ScoredMove
	}{
		{index: 1, want: []ScoredMove{moves[0], moves[4], moves[2], moves[3]}},
		{index: 3, want: []ScoredMove{moves[0], moves[4], moves[2]}},
		{index: 0, want: []ScoredMove{moves[2], moves[4]}},
		{index: 1, want: []ScoredMove{moves[2]}},
		{index: 0, want: []ScoredMove{}},
	}
	var movelist Movelist
	for _, move := range moves {
		movelist.AddScoredMove(move)
	}

	for _, tt := range tests {
		movelist.Remove(tt.index)

		if movelist.Len != len(tt.want) {
			t.Fatalf("after removing index %d, got len %d, want %d", tt.index, movelist.Len, len(tt.want))
		}
		for i, expected := range tt.want {
			got := movelist.Moves[i]
			if got != expected {
				t.Errorf("after removing index %d, Moves[%d] = %v, want %v", tt.index, i, got, expected)
			}
		}
	}
}
