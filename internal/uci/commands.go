package uci

import (
	"fmt"
	"io"
	"strconv"
	"strings"
	"time"

	"github.com/masterstruct/Eunoia/internal/board"
	"github.com/masterstruct/Eunoia/internal/movegen"
	"github.com/masterstruct/Eunoia/internal/search"
)

func (e *Engine) handleGo(w io.Writer, args []string) {
	e.mu.Lock()
	e.state.Stop = true
	e.mu.Unlock()

	e.running.Wait()

	state := &search.SearchState{}
	state.Reset() // TODO: once have TT, do not clear on go command

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

	p, err := applyMoves(&e.pos, rest[1:])
	if err != nil {
		return err
	} else {
		e.pos = p
	}
	return nil
}

func applyMoves(pos *board.Position, moves []string) (board.Position, error) {
	newPos := *pos

	for _, move := range moves {
		success := false

		var movelist movegen.Movelist
		// TODO: make this faster by only generating the moves
		// of pieces that are on the FROM square
		movegen.GeneratePseudolegalMoves(&newPos, &movelist)

		for i := range movelist.Len {
			m := movelist.Moves[i]
			if m.String() == move {
				newPos = newPos.MakeMove(m)
				success = true
				break
			}
		}
		if !success {
			return *pos, fmt.Errorf("uci: illegal or unknown move %q", move)
		}
	}
	return newPos, nil
}
