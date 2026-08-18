package movegen

import (
	"fmt"
	"testing"

	"github.com/masterstruct/Eunoia/internal/board"
)

var perftPositions = []struct {
	name  string
	fen   string
	depth int
	want  []uint64
}{
	{
		name:  "starting position",
		fen:   board.StartingFEN,
		depth: 5,
		want:  []uint64{20, 400, 8902, 197281, 4865609, 119060324, 3195901860, 84998978956, 2439530234167},
	},
	{
		name:  "kiwipete",
		fen:   "r3k2r/p1ppqpb1/bn2pnp1/3PN3/1p2P3/2N2Q1p/PPPBBPPP/R3K2R w KQkq -",
		depth: 4,
		want:  []uint64{48, 2039, 97862, 4085603, 193690690, 8031647685},
	},
	{
		name:  "position 3",
		fen:   "8/2p5/3p4/KP5r/1R3p1k/8/4P1P1/8 w - - 0 1",
		depth: 6,
		want:  []uint64{14, 191, 2812, 43238, 674624, 11030083, 178633661, 3009794393},
	},
	{
		name:  "position 4",
		fen:   "r3k2r/Pppp1ppp/1b3nbN/nP6/BBP1P3/q4N2/Pp1P2PP/R2Q1RK1 w kq - 0 1",
		depth: 5,
		want:  []uint64{6, 264, 9467, 422333, 15833292, 706045033},
	},
	{
		name:  "position 5",
		fen:   "rnbq1k1r/pp1Pbppp/2p5/8/2B5/8/PPP1NnPP/RNBQK2R w KQ - 1 8",
		depth: 4,
		want:  []uint64{44, 1486, 62379, 2103487, 89941194},
	},
	{
		name:  "position 6",
		fen:   "r4rk1/1pp1qppp/p1np1n2/2b1p1B1/2B1P1b1/P1NP1N2/1PP1QPPP/R4RK1 w - - 0 10",
		depth: 4,
		want:  []uint64{46, 2079, 89890, 3894594, 164075551, 6923051137, 287188994746, 11923589843526, 490154852788714},
	},
	{
		name:  "position 7",
		fen:   "r3k2r/Pppp1ppp/1b3nbN/nP6/BBP1P3/q4N2/P2P2PP/q2Q1R1K w kq - 0 2",
		depth: 4,
		want:  []uint64{37, 1766, 67665, 3251340, 126325537, 6084758848, 238196369383},
	},
}

func TestPerft(t *testing.T) {
	var total uint64
	var n uint64
	for _, tt := range perftPositions {
		depth := tt.depth
		t.Run(tt.name, func(t *testing.T) {
			pos, err := board.ParseFEN(tt.fen)
			if err != nil {
				t.Fatalf("bad FEN: %v", err)
			}

			if depth < 1 {
				t.Fatalf("depth must be >= 1: %d", depth)
			}

			if depth > len(tt.want) {
				depth = len(tt.want)
			}

			got := Perft(&pos, depth)
			want := tt.want[depth-1]

			if got.Nodes != want {
				t.Errorf("depth %d: got %d, want %d", depth, got.Nodes, want)
			}
			fmt.Println("total:", got.Nodes)
			fmt.Println("time:", got.Time)
			fmt.Println("nps:", got.NPS)
			total += got.NPS
			n++
		})
	}
	fmt.Println("Average nps:", total/n)
}
