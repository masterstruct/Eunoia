package uci

import (
	"testing"

	"github.com/masterstruct/Eunoia/internal/board"
)

var moves []string = []string{"e2e4", "c7c5", "c2c3", "d7d5", "e4d5", "d8d5", "d2d4", "g8f6", "g1f3", "g7g6", "d1b3", "c5d4", "f1c4", "d5e4", "e1f1", "e7e6", "c3d4", "e4f5", "b1c3", "a7a6", "c1g5", "b7b5", "a1e1", "b5c4", "b3c4", "b8d7", "g5f6", "f5f6", "c3d5", "f6f5", "d5c7", "e8d8", "c7a8", "f8d6", "f3e5", "d8e7", "c4c6", "d7e5", "d4e5", "d6e5", "a8b6", "e5d4", "b6d5", "e7f8", "d5e3", "f5b5", "c6b5", "a6b5", "e1e2", "e6e5", "e2d2", "c8e6", "b2b3", "f8g7", "f1e2", "h8c8", "d2c2", "c8a8", "h1d1", "h7h5", "d1d2", "h5h4", "h2h3", "b5b4", "c2c7", "d4c3", "d2c2", "g7f8", "c7b7", "f8g7", "f2f3", "g7f6", "e3d1", "e6f5", "d1c3", "b4c3", "c2c3", "a8a2", "e2f1", "a2a1", "f1f2", "a1a2", "f2g1", "a2a1", "g1h2", "a1e1", "c3c7", "f5e6", "b3b4", "e1e2", "b7b6", "f6g5", "c7c5", "g5f4", "b6b8", "e2b2", "b4b5", "f4f5", "b5b6", "f7f6", "b6b7", "g6g5", "c5a5", "b2b6", "h2g1", "b6b1", "g1f2", "b1b2", "f2e1", "b2b6", "a5c5", "e6f7", "c5c7", "f7d5", "c7d7", "d5e6", "d7e7", "e6d5", "e7e5", "f6e5", "b8f8", "d5f7", "b7b8q", "b6b8", "f8b8", "f7g6", "b8f8", "f5e6", "e1f2", "e6e7", "f8h8", "e7f6", "h8b8", "g6h7", "b8b6", "f6f5", "b6b7", "h7g6", "b7g7", "f5f6", "g7g8", "f6f7", "g8c8", "f7f6", "c8f8", "f6e6", "f8d8", "e6f5", "d8d6", "g6f7", "d6d7", "f7e8", "d7c7", "e8g6", "f2e3", "f5f6", "c7c6", "f6f5", "c6b6", "g6h7", "b6b7", "h7g6", "b7a7", "f5f6", "a7a6", "f6f7", "a6d6", "g6c2", "d6b6", "c2g6", "b6a6", "g6c2", "a6b6", "c2g6", "b6a6", "g6c2", "e3f2", "c2d3", "a6b6", "d3f5", "f2e2", "f7e7", "b6a6", "e7f7", "e2f2", "f5d3", "a6c6", "d3f5", "f2e3", "f5g6", "c6b6", "g6f5", "b6c6", "f5b1", "e3e2", "b1f5", "c6b6", "f5e6", "b6b8", "e6f5", "e2e3", "f5g6", "b8b6", "g6c2", "e3f2", "c2f5", "b6c6", "f7e7", "c6a6", "e7f7", "a6c6", "f5e6", "c6d6", "e6f5", "f2f1", "f5e6", "f1e1", "f7e7", "d6b6", "e7f7", "e1d2", "f7f6", "d2d3", "f6f5"}
var fen string = "8/8/1R2b3/4pkp1/7p/3K1P1P/6P1/8 w - - 99 113"

func TestApplyMoves(t *testing.T) {
	pos := board.StartingPosition()
	got, err := applyMoves(&pos, moves)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	got.CastlingRookSq = board.CastlingRookSquares{
		WhiteKingside:  board.NoSquare,
		WhiteQueenside: board.NoSquare,
		BlackKingside:  board.NoSquare,
		BlackQueenside: board.NoSquare,
	}

	want, err := board.ParseFEN(fen)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	if got != want {
		t.Errorf("position mismatch:\ngot:\n%v\nwant: \n%v\n", got, want)
	}
}

func BenchmarkApplyMoves(b *testing.B) {
	pos := board.StartingPosition()
	for b.Loop() {
		got, _ := applyMoves(&pos, moves)
		_ = got
	}
}
