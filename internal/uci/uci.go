package uci

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"strings"

	"github.com/masterstruct/Eunoia/internal/board"
	"github.com/masterstruct/Eunoia/internal/movegen"
	"github.com/masterstruct/Eunoia/internal/search"
)

type Engine struct {
	pos board.Position
}

func NewEngine() Engine {
	return Engine{pos: board.StartingPosition()}
}

func Loop(r io.Reader, w io.Writer) {
	scanner := bufio.NewScanner(r)
	eng := NewEngine()

	var searchCancel context.CancelFunc
	stopSearch := func() {
		if searchCancel != nil {
			searchCancel()
			searchCancel = nil
		}
	}

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
			eng.pos = board.StartingPosition()

		case "position":
			if err := eng.handlePosition(args); err != nil {
				fmt.Fprintf(w, "info string %v\n", err)
			}

		case "go":
			stopSearch()
			move := search.Search(eng.pos)
			if move == board.NullMove {
				fmt.Fprintln(w, "bestmove 0000")
				break
			}
			eng.pos = eng.pos.MakeMove(move)
			fmt.Fprintf(w, "bestmove %s\n", move.String())

		case "stop":
			stopSearch()

		case "quit":
			stopSearch()
			return

		case "d":
			fmt.Fprintln(w, eng.pos.String())
		}
	}
}

// position [fen <fenstring> | startpos] [moves <move1> ... <movei>]
func (e *Engine) handlePosition(args []string) error {
	if len(args) == 0 {
		return fmt.Errorf("uci: position requires arguments")
	}

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
