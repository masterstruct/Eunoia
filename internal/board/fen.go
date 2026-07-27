package board

import (
	"strconv"
	"strings"
)

func ParseFEN(fen string) (Position, error) {
	pos := NewPosition()
	file := FileA
	rank := Rank8
	splits := strings.Split(fen, " ")
	n := len(splits)

	if n >= 1 {
		for _, char := range splits[0] {
			if char == '/' {
				file = 0
				rank -= 1
				continue
			}
			skip, err := strconv.Atoi(string(char))
			if err == nil && !(skip < 1 || skip > 8) {
				file += skip
				continue
			} // TODO: handle error

			pt := ParsePieceType(byte(char))
			color := NoColor
			if pt != NoPieceType {
				// char is a chess piece
				if char&0x20 != 0 {
					// uppercase => Black
					color = Black
				} else {
					color = White
				}
				pos.PlacePiece(NewPiece(pt, color), NewSquare(file, rank))
				file += 1
				continue
			}
		}
	}
	if n >= 2 {
		pos.SideToMove = ParseColor(splits[1][0])
	}
	if n >= 3 {
		// TODO: handle error
		rights, _ := ParseCastlingRights(splits[2])
		pos.CastlingRights = rights
	}
	if n >= 4 {
		// TODO: handle error
		sq, _ := ParseSquare(splits[3])
		pos.EnPassant = sq
	}
	if n >= 5 {
		halfMove, err := strconv.Atoi(splits[4])
		if err == nil {
			pos.HalfmoveClock = uint8(halfMove)
		} // TODO: handle error
	}
	if n >= 6 {
		ply, err := strconv.Atoi(splits[5])
		if err == nil {
			pos.Ply = uint16(ply)
		} // TODO: handle error
	}

	return pos, nil
}

func (pos Position) FEN() string {
	// TODO: use Bitboard.PopLSB() to loop faster
	var sb strings.Builder
	skip := 0

	for rank := Rank8; rank >= Rank1; rank -= 1 {
		for i := range 8 {
			sq := NewSquare(i, rank)

			piece, ok := pos.PieceOn(sq)
			if ok {
				// found piece
				if skip > 0 {
					sb.WriteString(strconv.Itoa(skip))
					skip = 0
				}
				sb.WriteByte(piece.String())
				continue
			}
			skip += 1
		}

		// reached end of rank
		if skip > 0 {
			sb.WriteString(strconv.Itoa(skip))
			skip = 0
		}
		if rank != Rank1 {
			sb.WriteByte('/')
		}
	}

	return sb.String()
}
