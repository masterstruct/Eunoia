package search

import "github.com/masterstruct/Eunoia/internal/board"

func (ss *SearchState) isDraw(pos *board.Position) bool {
	return pos.HalfmoveClock >= 100 || isInsufficientMaterial(pos) || ss.Keys.IsRepetition(pos.Hash, pos.HalfmoveClock)
}

func (k *Keys) IsRepetition(hash uint64, halfmoveClock uint8) bool {
	n := len(k.keys)
	if n == 0 {
		return false
	}
	cur := n - 1
	reps := 0
	for d := 4; d <= int(halfmoveClock) && d <= cur; d += 2 {
		idx := cur - d
		if k.keys[idx] != hash {
			continue
		}
		if idx >= k.len-1 {
			return true
		}
		reps++
		if reps >= 2 {
			return true
		}
	}
	return false
}

func isInsufficientMaterial(pos *board.Position) bool {
	pawns := pos.Pieces[board.Pawn]
	rooks := pos.Pieces[board.Rook]
	queens := pos.Pieces[board.Queen]
	if (pawns | rooks | queens) != 0 {
		return false
	}

	knights := pos.Pieces[board.Knight]
	bishops := pos.Pieces[board.Bishop]

	// impossible to mate with less than 2 minor pieces
	return (knights | bishops).CountBits() < 2
}

type Keys struct {
	keys []uint64
	len  int
}

func (k *Keys) Push(hash uint64) {
	k.keys = append(k.keys, hash)
}

func (k *Keys) Pop() {
	k.keys = k.keys[:len(k.keys)-1]
}

func (k *Keys) SetFromGame(gameHistory *Keys) {
	k.keys = append(k.keys[:0], gameHistory.keys[:gameHistory.len]...)
	k.len = gameHistory.len
}

func (k *Keys) SetGamePosition(hash uint64) {
	k.keys = append(k.keys[:0], hash)
	k.len = 1
}

func (k *Keys) AppendGame(hashes ...uint64) {
	k.keys = append(k.keys, hashes...)
	k.len += len(hashes)
}
