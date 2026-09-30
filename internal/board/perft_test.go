package board

import (
	"fmt"
	"strconv"
	"strings"
	"sync"
	"testing"

	"github.com/masterstruct/Eunoia/resources"
)

const (
	standardMaxDepth = 5
	chess960MaxDepth = 4
)

func TestPerft(t *testing.T) {
	positions, err := loadPerftPositions("standard.epd", standardMaxDepth)
	if err != nil {
		t.Fatal(err)
	}
	runPerftTests(t, positions, false, false, true)
}

func TestPerft_Chess960(t *testing.T) {
	withChess960(t)
	positions, err := loadPerftPositions("frc.epd", chess960MaxDepth)
	if err != nil {
		t.Fatal(err)
	}
	runPerftTests(t, positions, false, true, false)
}

type perftPosition struct {
	name string
	fen  string
	want []uint64
}

func loadPerftPositions(filename string, maxDepth int) ([]perftPosition, error) {
	if maxDepth < 1 {
		return nil, fmt.Errorf("maximum depth must be positive: %d", maxDepth)
	}

	file, err := resources.Perft.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("open %s: %w", filename, err)
	}
	defer file.Close()

	var positions []perftPosition
	err = ParseEPD(file, func(epd EPD) error {
		positionIndex := len(positions) + 1
		position := perftPosition{
			name: fmt.Sprintf("position-%d", positionIndex),
			fen:  epd.FEN,
		}
		seenDepths := make(map[int]bool, len(epd.Operations))
		for _, operation := range epd.Operations {
			if !strings.HasPrefix(operation.Key, "D") {
				return fmt.Errorf("position %d: unexpected operation %q", positionIndex, operation.Key)
			}
			depth, err := strconv.Atoi(strings.TrimPrefix(operation.Key, "D"))
			if err != nil || depth < 1 {
				return fmt.Errorf("position %d: invalid depth %q", positionIndex, operation.Key)
			}
			if depth > maxDepth {
				continue
			}
			if seenDepths[depth] {
				return fmt.Errorf("position %d: duplicate depth %d", positionIndex, depth)
			}
			seenDepths[depth] = true
			if depth > len(position.want) {
				position.want = append(position.want, make([]uint64, depth-len(position.want))...)
			}
			nodes, err := strconv.ParseUint(operation.Value, 10, 64)
			if err != nil {
				return fmt.Errorf("position %d: invalid node count %q", positionIndex, operation.Value)
			}
			position.want[depth-1] = nodes
		}
		if len(position.want) == 0 {
			return fmt.Errorf("position %d: no perft counts", positionIndex)
		}
		for depth := range position.want {
			if !seenDepths[depth+1] {
				return fmt.Errorf("position %d: missing D%d count", positionIndex, depth+1)
			}
		}

		positions = append(positions, position)
		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("parse %s: %w", filename, err)
	}
	if len(positions) == 0 {
		return nil, fmt.Errorf("%s contains no perft positions", filename)
	}
	return positions, nil
}

func runPerftTests(t *testing.T, positions []perftPosition, splitperft bool, parallel bool, print bool) {
	t.Helper()

	var total uint64
	var n uint64
	var mu sync.Mutex

	for _, tt := range positions {
		t.Run(tt.name, func(t *testing.T) {
			if parallel {
				t.Parallel()
			}

			pos, err := ParseFEN(tt.fen)
			if err != nil {
				t.Fatalf("bad FEN: %v", err)
			}

			for depth, want := range tt.want {
				depth++
				var got PerftResult
				if splitperft {
					got = SplitPerft(&pos, depth)
				} else {
					got = Perft(&pos, depth)
				}
				if got.Nodes != want {
					t.Errorf("depth %d: got %d, want %d", depth, got.Nodes, want)
				}

				if depth == len(tt.want) {
					if print {
						fmt.Printf("%s depth %d: nodes %d, time %v, nps %d\n", tt.name, depth, got.Nodes, got.Time, got.NPS)
					}
					mu.Lock()
					total += got.NPS
					n++
					mu.Unlock()
				}
			}
		})
	}

	t.Cleanup(func() {
		if n > 0 {
			s := fmt.Sprint("Average nps: ", total/n)
			fmt.Println(strings.Repeat("~", len(s)))
			fmt.Println(s)
			fmt.Println(strings.Repeat("~", len(s)))
		}
	})
}
