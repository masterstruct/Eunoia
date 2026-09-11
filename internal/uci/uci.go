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

func newEngine(ttSizeMiB uint) *engine {
	e := &engine{
		pos:   board.StartingPosition(),
		state: &search.SearchState{},
	}
	e.gameHistory = []uint64{e.pos.Hash}
	e.state.Init(ttSizeMiB)
	return e
}

func Loop(r io.Reader, w io.Writer) {
	if len(os.Args) > 1 && os.Args[1] == "bench" {
		Bench()
		return
	}

	scanner := bufio.NewScanner(r)

	eng := newEngine(tt.DefaultSizeMiB)

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
			eng.mu.Lock()
			eng.state.Stop = true
			eng.mu.Unlock()
			eng.running.Wait()
			return

		case "uci":
			fmt.Fprintln(w, "id name Eunoia")
			fmt.Fprintln(w, "id author Master Struct")
			fmt.Fprintln(w, "option name Threads type spin default 1 min 1 max 1")
			fmt.Fprintln(w, "option name Hash type spin default 64 min 1 max 33554432")
			fmt.Fprintln(w, "option name UCI_Chess960 type check default false")
			fmt.Fprintln(w, "uciok")

		case "setoption":
			eng.handleSetOption(args)

		case "position":
			if err := eng.handlePosition(args); err != nil {
				fmt.Fprintf(w, "info string %v\n", err)
			}

		case "ucinewgame":
			eng.mu.Lock()
			eng.state.Stop = true
			eng.mu.Unlock()
			eng.running.Wait()
			eng.mu.Lock()
			eng.state.Reset()
			eng.state.ClearButterflyHistory()
			eng.state.ClearTT()
			eng.mu.Unlock()

		case "isready":
			fmt.Fprintln(w, "readyok")

		case "go":
			eng.handleGo(w, args)

		case "stop":
			eng.mu.Lock()
			eng.state.Stop = true
			eng.mu.Unlock()
			eng.running.Wait()

		case "perft":
			depth := 0
			if len(args) > 0 {
				depth, _ = strconv.Atoi(args[0])
			}
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

		case "flip":
			eng.mu.Lock()
			eng.pos.SideToMove = eng.pos.SideToMove.Opponent()
			eng.mu.Unlock()
		}
	}
}
