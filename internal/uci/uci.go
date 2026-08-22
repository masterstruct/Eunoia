package uci

import (
	"bufio"
	"fmt"
	"io"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/masterstruct/Eunoia/internal/board"
	"github.com/masterstruct/Eunoia/internal/movegen"
	"github.com/masterstruct/Eunoia/internal/search"
)

type Engine struct {
	mu      sync.Mutex
	pos     board.Position
	state   *search.SearchState
	running sync.WaitGroup
}

func NewEngine() Engine {
	return Engine{
		pos:   board.StartingPosition(),
		state: &search.SearchState{},
	}
}

func Loop(r io.Reader, w io.Writer) {
	scanner := bufio.NewScanner(r)

	eng := NewEngine()

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}
		fields := strings.Fields(line)
		cmd, args := fields[0], fields[1:]

		switch cmd {
		case "uci":
			fmt.Fprintln(w, "id name Eunoia")
			fmt.Fprintln(w, "id author Master Struct")
			fmt.Fprintln(w, "option name Hash type spin default 1 min 1 max 1")
			fmt.Fprintln(w, "option name Threads type spin default 1 min 1 max 1")
			fmt.Fprintln(w, "uciok")

		case "isready":
			fmt.Fprintln(w, "readyok")

		case "ucinewgame":
			eng.mu.Lock()
			eng.pos = board.StartingPosition()
			eng.mu.Unlock()

		case "position":
			if err := eng.handlePosition(args); err != nil {
				fmt.Fprintf(w, "info string %v\n", err)
			}

		case "go":
			eng.handleGo(args, w)

		case "perft":
			depth := 0
			if len(args) > 0 {
				depth, _ = strconv.Atoi(args[0])
			}
			eng.mu.Lock()
			perftRes := movegen.Perft(&eng.pos, depth)
			eng.mu.Unlock()
			fmt.Fprintln(w, "total:", perftRes.Nodes)
			fmt.Fprintln(w, "time:", perftRes.Time)
			fmt.Fprintln(w, "nps:", perftRes.NPS)

		case "stop":
			eng.mu.Lock()
			eng.state.Stop = true
			eng.mu.Unlock()
			eng.running.Wait()

		case "quit":
			eng.mu.Lock()
			eng.state.Stop = true
			eng.mu.Unlock()
			eng.running.Wait()
			return

		case "d":
			eng.mu.Lock()
			s := eng.pos.String()
			eng.mu.Unlock()
			fmt.Fprintln(w, s)
		}
	}
}

// position [fen <fenstring> | startpos] [moves <move1> ... <movei>]
func (e *Engine) handlePosition(args []string) error {
	if len(args) == 0 {
		return fmt.Errorf("uci: position requires arguments")
	}

	e.mu.Lock()
	defer e.mu.Unlock()

	var rest []string

	switch args[0] {
	case "startpos":
		e.pos = board.StartingPosition()
		rest = args[1:]

	case "fen":
		end := len(args)
		for i := 1; i < len(args); i++ {
			if args[i] == "moves" {
				end = i
				break
			}
		}

		pos, err := board.ParseFEN(strings.Join(args[1:end], " "))
		if err != nil {
			return fmt.Errorf("uci: %w", err)
		}
		e.pos = pos
		rest = args[end:]

	default:
		return fmt.Errorf("uci: unknown position type %q", args[0])
	}

	if len(rest) == 0 {
		return nil
	}
	if rest[0] != "moves" {
		return fmt.Errorf("uci: expected \"moves\", got %q", rest[0])
	}

	for _, moveStr := range rest[1:] {
		if err := e.applyMove(moveStr); err != nil {
			return err
		}
	}
	return nil
}

// assumes e.mu is owned by caller (handlePosition)
func (e *Engine) applyMove(moveStr string) error {
	// TODO: make this better
	var ml movegen.Movelist
	movegen.GeneratePseudolegalMoves(&e.pos, &ml)

	for i := range ml.Len {
		m := ml.Moves[i]
		if m.String() == moveStr {
			e.pos = e.pos.MakeMove(m)
			return nil
		}
	}
	return fmt.Errorf("uci: illegal or unknown move %q", moveStr)
}

func (e *Engine) handleGo(args []string, w io.Writer) {
	e.mu.Lock()
	e.state.Stop = true
	e.mu.Unlock()

	e.running.Wait()

	state := &search.SearchState{}
	state.Reset()

	e.mu.Lock()
	pos := e.pos
	e.mu.Unlock()

	depth := 1000 // "infinite"

	timeLeft := 0
	moveTime := 0
	increment := 0

	// returns the integer following args[i] and ok bool
	intArg := func(i int) (int, bool) {
		if i+1 >= len(args) {
			return 0, false
		}
		n, err := strconv.Atoi(args[i+1])
		return n, err == nil
	}

	for i, arg := range args {
		switch arg {
		case "movetime":
			if v, ok := intArg(i); ok {
				moveTime = v
			}
		case "nodes":
			if v, ok := intArg(i); ok {
				state.MaxNodes = uint64(v)
			}
		case "wtime":
			if v, ok := intArg(i); ok && pos.SideToMove == board.White {
				timeLeft = v
			}
		case "btime":
			if v, ok := intArg(i); ok && pos.SideToMove == board.Black {
				timeLeft = v
			}
		case "winc":
			if v, ok := intArg(i); ok && pos.SideToMove == board.White {
				increment = v
			}
		case "binc":
			if v, ok := intArg(i); ok && pos.SideToMove == board.Black {
				increment = v
			}
		case "depth":
			if v, ok := intArg(i); ok {
				depth = v
			}
		case "infinite":
			depth = 1000
		}
	}

	switch {
	case moveTime > 0:
		state.MaxTime = state.StartTime.Add(time.Duration(moveTime) * time.Millisecond)
	case timeLeft > 0 || increment > 0:
		// TODO: add safety margin
		hard := timeLeft/3 + (increment*7)/10
		if hard > 0 {
			soft := timeLeft/30 + (increment*7)/10
			state.SoftTime = state.StartTime.Add(time.Duration(soft) * time.Millisecond)
			state.MaxTime = state.StartTime.Add(time.Duration(hard) * time.Millisecond)
		}
	}

	e.mu.Lock()
	e.state = state
	e.mu.Unlock()

	e.running.Go(func() {
		move := state.SearchBestMove(pos, depth)

		if move == board.NullMove {
			fmt.Fprintln(w, "bestmove 0000")
			return
		}

		e.mu.Lock()
		e.pos = pos.MakeMove(move)
		e.mu.Unlock()

		fmt.Fprintf(w, "bestmove %s\n", move.String())
	})
}
