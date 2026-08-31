package uci

import (
	"fmt"
	"time"

	"github.com/masterstruct/Eunoia/internal/board"
)

const benchDepth = 5

type benchResult struct {
	nodes uint64
	time  time.Duration
}

func Bench() {
	eng := newEngine()

	var totalNodes uint64
	var totalTime time.Duration

	fmt.Println("Benchmarking...")

	for _, benchStruct := range benchPositions {
		eng.state.Reset()
		eng.state.ClearButterflyHistory()
		eng.state.Init()
		eng.state.Quiet = true

		board.SetChess960(benchStruct.frc)
		pos, err := board.ParseFEN(benchStruct.fen)
		if err != nil {
			panic("Bad bench fen: " + benchStruct.fen)
		}
		eng.pos = pos

		res := eng.searchBenchPos()

		totalNodes += res.nodes
		totalTime += res.time
	}
	fmt.Println("\n-=- Benchmark Results -=-")
	fmt.Println("Total time (ms):", totalTime.Milliseconds())
	fmt.Println("Nodes searched:", totalNodes)
	fmt.Println("Nodes/second:", totalNodes*1_000_000_000/uint64(totalTime.Nanoseconds()))
	fmt.Println("-==-==-==-==-==-==-==-==-")
}

func (e *engine) searchBenchPos() benchResult {
	start := time.Now()

	e.state.SearchBestMove(e.pos, benchDepth)

	return benchResult{
		nodes: e.state.Nodes,
		time:  time.Since(start),
	}
}

var benchPositions = []struct {
	fen string
	frc bool
}{
	// Standard positions
	{"r2q1r2/2pn1pbk/bp1p1np1/3Pp2p/1PP1P3/N5P1/2NBQPBP/1R3RK1 b - - 1 19", false},
	{"r2qnrk1/pbnp2bp/1pp5/3Ppp2/2P5/1PN3P1/PBN2PBP/1R1Q1RK1 b - - 2 15", false},
	{"rnr3k1/ppq2p2/2pb2p1/3p3p/2PPp3/1PN1P1P1/P4PPN/R2Q1RK1 b - - 1 18", false},
	{"r3n1k1/5pp1/1q5p/8/2Qn4/3B4/1P3PPP/2KNR3 w - - 14 44", false},
	{"r1bqk2r/3n1pbp/p1p1p1p1/3pP3/5P2/3BB3/PPPN2PP/R2Q1RK1 w kq - 1 12", false},
	{"r1b2rk1/4qpb1/p5p1/3pP3/2p2Pn1/4BN2/PPB3P1/R3QRK1 w - - 2 20", false},
	{"2q3k1/5rbN/5Bp1/3p4/p1p2P1Q/P7/1P4P1/3R3K b - - 0 34", false},
	{"rnb1k2r/p3ppb1/1ppp4/5PBp/3P2n1/q1N4B/P1PQN2P/1R2K2R w Kkq h6 0 13", false},
	{"1r3rk1/4bpp1/p2p3p/2qPP3/P5PP/1p6/1PP1Q2R/1K1R1B2 w - - 1 24", false},
	{"5K2/3r4/8/2n3k1/4n1p1/8/8/5R2 b - - 9 77", false},
	{"6k1/2R3p1/2n5/3bnP2/6Pp/7P/p4K2/B7 b - - 68 114", false},
	{"8/6p1/p1kNp3/1bp1P2P/R2p2P1/P4PK1/8/2b5 w - - 2 60", false},
	{"4k2r/1p1qbp1p/p3p1pB/3pP3/1PnPN1QP/P7/5PP1/2R3K1 b k - 2 27", false},
	{"3n1r1k/4r2p/1p1q4/p3p3/1PQpR3/P5P1/5PKP/2R1N3 w - - 2 68", false},
	{"6k1/1R3Rp1/7p/6r1/2N5/8/5K2/3r1B2 w - - 8 64", false},
	{"1r3kb1/3nrpp1/7p/1p2PP1P/pP3KP1/4R3/5N2/4RB2 b - - 0 39", false},
	{"r4rk1/pp1qbp1p/3p2p1/3Pp3/1P2P1n1/3P1N2/3B1PPP/R2QK2R w KQ - 5 15", false},
	{"8/4P2k/3q4/P5Q1/5P2/4K3/8/8 b - - 4 139", false},
	{"6k1/8/8/2pQ2p1/P7/2q1P3/4KP2/8 b - - 1 113", false},
	{"5bk1/p2q2p1/6Q1/1P1r2P1/1n1p1p1p/1P3N1P/1B2rP1K/R4R2 w - - 0 28", false},
	{"r2q1rk1/pppn2bp/2P5/3Pp1pb/4Pp2/2N2N1P/PP2QBP1/R4RK1 b - - 0 17", false},
	{"8/r4pk1/1n1p2p1/3B2Pp/1P1R1P1P/2P2R2/2K5/7r w - - 7 46", false},
	{"1r2r3/4ppkp/b1pq2p1/8/p2PP2P/4QN2/P1R2PP1/4R1K1 b - - 0 28", false},
	{"k1r2b2/p1rqnp2/1p2p1p1/1P1pP1Pp/1B1P1P1N/1Q1Bn1KP/R1P5/R7 w - - 73 61", false},
	{"8/4bpk1/pB4p1/8/PnB1K1P1/7P/8/8 b - - 7 44", false},

	// FRC positions
	{"nr2krbq/pppp2pp/2n2b2/2P1p3/4p3/1N5P/PP1P1PPB/1NQBRRK1 b fb - 2 7", true},
	{"bbr2krq/pp6/4nn2/P1p1pppp/2Np4/1N1PP1P1/1PP2P1P/2KRBRQB b gc - 1 11", true},
	{"b2bqrkr/p4p1p/1p4nP/2pp2PN/4pQ1n/1P5B/P1PPPP2/1RNK2R1 w GB - 9 14", true},
	{"3r4/p2rnkp1/2n2p2/bpp4p/5P2/1PPPNNP1/PBKR3P/4R3 b - - 30 35", true},
	{"8/3p4/kpn1b3/pNp3B1/PnP1P3/1P6/3P3r/1K1B2R1 b - - 40 48", true},
	{"r1kb3r/1p1q4/2pBp3/2Pp1p1p/1Q1P2p1/2R5/P3PPPP/5RK1 b - - 4 19", true},
	{"rbbn1rkq/pp1n1pp1/2p5/3pP3/3N4/3N1PPp/PPPQ3P/RK3BBR w AH - 1 12", true},
	{"r1k1r1q1/1p1n1n2/7p/PPp1pb2/2P2p1P/2B2N2/2P1B2R/Q1K3R1 b - - 5 26", true},
	{"n2rkr1b/pp2pp2/3q2p1/3p1b2/3Pn3/1N3N2/PPB1PPP1/1QB2RKR w HFfd - 6 10", true},
}
