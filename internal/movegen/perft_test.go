package movegen

import (
	"fmt"
	"testing"

	"github.com/masterstruct/Eunoia/internal/board"
)

func TestPerft(t *testing.T) {
	runPerftTests(t, perftPositions, false)
}

func TestPerft_Chess960(t *testing.T) {
	withChess960(t)
	runPerftTests(t, perftPositionsChess960, false)
}

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
		fen:   board.KiwipeteFEN,
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

var perftPositionsChess960 = []struct {
	name  string
	fen   string
	depth int
	want  []uint64
}{
	{
		name:  "position 1",
		fen:   "bqnb1rkr/pp3ppp/3ppn2/2p5/5P2/P2P4/NPP1P1PP/BQ1BNRKR w HFhf - 2 9",
		depth: 5,
		want:  []uint64{21, 528, 12189, 326672, 8146062, 227689589},
	},
	{
		name:  "position 93",
		fen:   "1nrbkr1q/1pppp1pp/1n6/p4p2/N1b4P/8/PPPPPPPB/N1RBKR1Q w FCfc - 2 9",
		depth: 4,
		want:  []uint64{27, 862, 24141, 755171, 22027695, 696353497},
	},
	{
		name:  "position 175",
		fen:   "nrnk1rbb/p1p2ppp/3pq3/Qp2p3/1P1P4/8/P1P1PPPP/NRN1KRBB w fb - 2 9",
		depth: 4,
		want:  []uint64{28, 873, 25683, 791823, 23868737, 747991356},
	},
	{
		name:  "position 246",
		fen:   "1rbkqbr1/ppp1pppp/1n5n/3p4/3P4/1PP3P1/P3PP1P/NRBKQBNR w HBb - 1 9",
		depth: 5,
		want:  []uint64{27, 752, 20686, 606783, 16986290, 521817800},
	},
	{
		name:  "position 519",
		fen:   "r1bqk1rb/pppnpppp/5n2/3p4/2P3PP/2N5/PP1PPP2/R1BQKNRB w GAga - 1 9",
		depth: 4,
		want:  []uint64{32, 821, 27121, 733155, 24923473, 710765657},
	},
	{
		name:  "position 787",
		fen:   "1rqknrnb/2pp1ppp/p3p3/1p6/P2P4/5bP1/1PP1PP1P/BRQKNRNB w FBfb - 0 9",
		depth: 4,
		want:  []uint64{24, 737, 20052, 598439, 17948681, 536330341},
	},
	{
		name:  "position 889",
		fen:   "rqkb1rnn/1pp1pp1p/p5p1/1b1p4/3P4/P5P1/RPP1PP1P/1QKBBRNN w Ffa - 1 9",
		depth: 5,
		want:  []uint64{21, 505, 11592, 290897, 7147063, 188559137},
	},
}

func runPerftTests(t *testing.T, positions []struct {
	name  string
	fen   string
	depth int
	want  []uint64
}, splitperft bool) {
	t.Helper()
	var total uint64
	var n uint64
	for _, tt := range positions {
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

			var got PerftResult
			if splitperft {
				got = SplitPerft(&pos, depth)
			} else {
				got = Perft(&pos, depth)
			}
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

	if n > 0 {
		fmt.Println("Average nps:", total/n)
	}
}
