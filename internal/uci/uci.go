package uci

import (
	"bufio"
	"fmt"
	"io"
	"strconv"
	"strings"

	"github.com/masterstruct/Eunoia/internal/board"
	"github.com/masterstruct/Eunoia/internal/movegen"
	"github.com/masterstruct/Eunoia/internal/search"
)

type Engine struct {
	pos   board.Position
	state search.SearchState
}

func NewEngine() Engine {
	return Engine{
		pos:   board.StartingPosition(),
		state: search.SearchState{},
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
			eng.pos = board.StartingPosition()

		case "position":
			if err := eng.handlePosition(args); err != nil {
				fmt.Fprintf(w, "info string %v\n", err)
			}

		case "go":
			eng.state.Reset()
			go func() {
				move := eng.state.SearchBestMove(eng.pos, 5)
				if move == board.NullMove {
					fmt.Fprintln(w, "bestmove 0000")
				} else {
					eng.pos = eng.pos.MakeMove(move)
					fmt.Fprintf(w, "bestmove %s\n", move.String())
				}
			}()

		case "perft":
			depth := 0
			if len(args) > 0 {
				depth, _ = strconv.Atoi(args[0])
			}
			perftRes := movegen.Perft(&eng.pos, depth)
			fmt.Fprintln(w, "total:", perftRes.Nodes)
			fmt.Fprintln(w, "time:", perftRes.Time)
			fmt.Fprintln(w, "nps:", perftRes.NPS)

		case "stop":
			eng.state.Stop = true

		case "quit":
			eng.state.Stop = true
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
