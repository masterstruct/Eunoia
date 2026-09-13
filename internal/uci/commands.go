package uci

import (
	"fmt"
	"io"
	"strconv"
	"strings"

	"github.com/masterstruct/Eunoia/internal/board"
	"github.com/masterstruct/Eunoia/internal/movegen"
	"github.com/masterstruct/Eunoia/internal/search"
)

func (e *engine) handleGo(w io.Writer, args []string) {
	e.mu.Lock()
	e.state.Stop = true
	e.mu.Unlock()

	e.running.Wait()

	e.mu.Lock()
	state := e.state
	state.PrepareForSearch()
	pos := e.pos
	state.SetHistory(e.gameHistory)
	e.mu.Unlock()

	limits := search.GoLimits{}

	// returns the integer following args[i] and ok bool
	intArg := func(i int) (int64, bool) {
		if i+1 >= len(args) {
			return 0, false
		}
		n, err := strconv.ParseInt(args[i+1], 10, 64)
		return n, err == nil
	}

	for i, arg := range args {
		switch arg {
		case "movetime":
			if v, ok := intArg(i); ok {
				limits.MoveTime = v
			}
		case "nodes":
			if v, ok := intArg(i); ok {
				limits.Nodes = v
			}
		case "wtime":
			if v, ok := intArg(i); ok {
				limits.WTime = v
			}
		case "btime":
			if v, ok := intArg(i); ok {
				limits.BTime = v
			}
		case "winc":
			if v, ok := intArg(i); ok {
				limits.WInc = v
			}
		case "binc":
			if v, ok := intArg(i); ok {
				limits.BInc = v
			}
		case "depth":
			if v, ok := intArg(i); ok {
				limits.Depth = v
			}
		case "infinite":
			limits.Infinite = true
		}
	}

	e.mu.Lock()
	e.state.SetLimits(limits, pos.SideToMove)
	e.mu.Unlock()

	e.running.Go(func() {
		move := state.SearchBestMove(pos)

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
func (e *engine) handlePosition(args []string) error {
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

	e.gameHistory = append(e.gameHistory[:0], e.pos.Hash)

	if len(rest) == 0 {
		return nil
	}
	if rest[0] != "moves" {
		return fmt.Errorf("uci: expected \"moves\", got %q", rest[0])
	}

	p, hashes, err := applyMoves(&e.pos, rest[1:])
	if err != nil {
		return err
	}
	e.pos = p
	e.gameHistory = append(e.gameHistory, hashes...)
	return nil
}

func applyMoves(pos *board.Position, moves []string) (board.Position, []uint64, error) {
	// TODO: if capture, don't generate non-captures with staged movegen

	newPos := *pos
	hashes := make([]uint64, 0, len(moves))

	for _, move := range moves {
		success := false
		n := len(move)
		if n < 4 || n > 5 {
			return *pos, nil, fmt.Errorf("uci: illegal move %q", move)
		}

		var movelist movegen.Movelist
		from, err := board.ParseSquare(move[:2])
		if err != nil {
			return *pos, nil, fmt.Errorf("uci: failed to parse move %q", move)
		}
		piece, ok := newPos.PieceOn(from)
		if !ok {
			return *pos, nil, fmt.Errorf("uci: illegal move %q", move)
		}

		switch piece.Type {
		case board.Pawn:
			movegen.GenPawnMoves(&newPos, &movelist)
		case board.Knight:
			movegen.GenKnightMoves(&newPos, &movelist)
		case board.Bishop:
			movegen.GenBishopMoves(&newPos, &movelist)
		case board.Rook:
			movegen.GenRookMoves(&newPos, &movelist)
		case board.Queen:
			movegen.GenQueenMoves(&newPos, &movelist)
		case board.King:
			movegen.GenKingMoves(&newPos, &movelist)
		}

		for i := range movelist.Len {
			m := movelist.Moves[i]
			if m.String() == move {
				newPos = newPos.MakeMove(m)
				hashes = append(hashes, newPos.Hash)
				success = true
				break
			}
		}
		if !success {
			return *pos, nil, fmt.Errorf("uci: illegal or unknown move %q", move)
		}
	}
	return newPos, hashes, nil
}

func (e *engine) handleSetOption(args []string) {
	n := len(args)
	if n < 2 || args[0] != "name" {
		return
	}

	end := n
	valueStart := -1
	for i := 1; i < n; i++ {
		if args[i] == "value" {
			end = i
			valueStart = i + 1
			break
		}
	}
	name := strings.Join(args[1:end], " ")

	value := ""
	if valueStart >= 0 && valueStart < n {
		value = strings.Join(args[valueStart:], " ")
	}

	switch name {
	case "UCI_Chess960":
		board.SetChess960(value == "true")
	case "Hash":
		mib, err := strconv.Atoi(value)
		if err == nil && mib > 0 {
			e.mu.Lock()
			e.state.Stop = true
			e.mu.Unlock()
			e.running.Wait()

			e.mu.Lock()
			e.state.ResizeTT(uint(mib))
			e.mu.Unlock()
		}
	case "Clear Hash":
		e.mu.Lock()
		e.state.Stop = true
		e.mu.Unlock()
		e.running.Wait()

		e.mu.Lock()
		e.state.ClearTT()
		e.mu.Unlock()
	}
}
