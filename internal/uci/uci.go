package uci

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
	"sync"

	"github.com/masterstruct/Eunoia/internal/board"
	"github.com/masterstruct/Eunoia/internal/movegen"
	"github.com/masterstruct/Eunoia/internal/search"
	"github.com/masterstruct/Eunoia/internal/tt"
)

type engine struct {
	mu          sync.Mutex
	pos         board.Position
	gameHistory []uint64
	state       *search.SearchState
	running     sync.WaitGroup
}

func newEngine() *engine {
	e := &engine{
		pos:   board.StartingPosition(),
		state: &search.SearchState{},
	}
	e.gameHistory = []uint64{e.pos.Hash}
	e.state.UpdateMoveOverhead(search.DefaultMoveOverhead)
	return e
}

func Loop(r io.Reader, w io.Writer) {
	if len(os.Args) > 1 && os.Args[1] == "bench" {
		Bench()
		return
	}

	scanner := bufio.NewScanner(r)

	eng := newEngine()
	eng.state.Init(tt.DefaultSizeMiB)

	for scanner.Scan() {
		_ = scanner.Err()

		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}

		fields := strings.Fields(line)
		cmd, args := fields[0], fields[1:]

		switch cmd {
		case "quit":
			eng.cancelSearchAndWait()
			return

		case "uci":
			fmt.Fprintln(w, "id name Eunoia")
			fmt.Fprintln(w, "id author Master Struct")
			fmt.Fprintln(w, "option name Threads type spin default 1 min 1 max 1")
			fmt.Fprintln(w, "option name Hash type spin default", tt.DefaultSizeMiB, "min 1 max 33554432")
			fmt.Fprintln(w, "option name Clear Hash type button")
			fmt.Fprintln(w, "option name UCI_Chess960 type check default false")
			fmt.Fprintln(w, "option name Move Overhead type spin default", search.DefaultMoveOverhead, "min", search.MinMoveOverhead, "max", search.MaxMoveOverhead)
			fmt.Fprintln(w, "uciok")

		case "setoption":
			eng.cancelSearchAndWait()
			eng.handleSetOption(args)

		case "position":
			if err := eng.handlePosition(args); err != nil {
				fmt.Fprintf(w, "info string %v\n", err)
			}

		case "ucinewgame":
			eng.cancelSearchAndWait()
			eng.mu.Lock()
			eng.state.PrepareForSearch()
			eng.state.ClearTables()
			eng.mu.Unlock()

		case "isready":
			fmt.Fprintln(w, "readyok")

		case "go":
			eng.cancelSearchAndWait()
			eng.handleGo(w, args)

		case "stop":
			eng.cancelSearchAndWait()

		case "perft":
			if len(args) == 0 {
				break
			}

			depth, err := strconv.Atoi(args[0])
			if err != nil {
				break
			}

			eng.cancelSearchAndWait()

			eng.mu.Lock()
			pos := eng.pos
			eng.mu.Unlock()

			perftRes := movegen.Perft(&pos, depth)
			fmt.Fprintln(w, "total:", perftRes.Nodes)
			fmt.Fprintln(w, "time:", perftRes.Time)
			fmt.Fprintln(w, "nps:", perftRes.NPS)

		case "d":
			eng.mu.Lock()
			s := eng.pos.String()
			eng.mu.Unlock()
			fmt.Fprintln(w, s)
		}
	}
}

func (e *engine) cancelSearch() {
	e.mu.Lock()
	e.state.Stop = true
	e.mu.Unlock()
}

// cancels and waits until the engine quits the search
func (e *engine) cancelSearchAndWait() {
	e.cancelSearch()
	e.running.Wait()
}
