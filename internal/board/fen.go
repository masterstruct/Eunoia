package board

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"unicode"
)

const startingFEN = "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1"

var (
	ErrInvalidFieldCount = errors.New("fen: expected 4 to 6 space-separated fields (piece placement, side to move, castling, en passant, [halfmove clock], [fullmove number])")
	ErrInvalidRankCount  = errors.New("fen: piece placement must have exactly 8 ranks separated by '/'")
	ErrInvalidPieceChar  = errors.New("fen: piece placement contains an unrecognized character")
	ErrInvalidRankDigit  = errors.New("fen: empty-square digit must be between 1 and 8")
	ErrInvalidRankLength = errors.New("fen: each rank must total exactly 8 squares")
)

func ParseFEN(fen string) (Position, error) {
	splits := strings.Fields(fen)
	n := len(splits)
	pos := NewPosition()

	if n < 4 || n > 6 {
		return pos, fmt.Errorf("%w: %q", ErrInvalidFieldCount, fen)
	}

	ranks := strings.Split(splits[0], "/")
	if len(ranks) != 8 {
		return pos, fmt.Errorf("%w: %q", ErrInvalidRankCount, fen)
	}
	for _, r := range ranks {
		if isWhitespaceOnly(r) || strings.ContainsRune(r, ' ') || len(r) > 8 {
			return pos, fmt.Errorf("%w: %q", ErrInvalidRankLength, fen)
		}
	}

	file := FileA
	rank := Rank8

	for _, char := range splits[0] {
		if char == '/' {
			// end of rank
			// check how many files were processed
			// in this rank -- must be exactly 8
			if file != 8 {
				return pos, fmt.Errorf("%w: %q", ErrInvalidRankLength, fen)
			}

			file = 0
			rank -= 1
			continue
		}
		skip, err := strconv.Atoi(string(char))
		if err == nil {
			// number
			if skip < 1 || skip > 8 {
				return pos, fmt.Errorf("%w: %q", ErrInvalidRankDigit, fen)
			}
			file += skip
			continue
		}

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
		} else {
			return pos, fmt.Errorf("%w: %q", ErrInvalidPieceChar, fen)
		}
	}

	pos.SideToMove = ParseColor(splits[1][0])

	// TODO: handle error
	rights, _ := ParseCastlingRights(splits[2])
	pos.CastlingRights = rights

	// TODO: handle error
	sq, _ := ParseSquare(splits[3])
	pos.EnPassant = sq

	if n >= 5 {
		halfMove, err := strconv.Atoi(splits[4])
		if err == nil {
			pos.HalfmoveClock = uint8(halfMove)
		} // TODO: handle error
	}
	if n >= 6 {
		fullmoves, err := strconv.Atoi(splits[5])
		if err == nil {
			pos.Ply = FullmovesToPly(uint16(fullmoves), pos.SideToMove)
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

	sb.WriteByte(' ')
	sb.WriteString(pos.SideToMove.String())
	sb.WriteByte(' ')
	sb.WriteString(pos.CastlingRights.String())
	sb.WriteByte(' ')
	sb.WriteString(pos.EnPassant.String())
	sb.WriteByte(' ')
	sb.WriteString(strconv.Itoa(int(pos.HalfmoveClock)))
	sb.WriteByte(' ')
	fullmoves, _ := PlyToFullmoves(pos.Ply)
	sb.WriteString(strconv.Itoa(int(fullmoves)))
	return sb.String()
}

func PlyToFullmoves(ply uint16) (fullmoves uint16, sideToMove Color) {
	if ply%2 == 0 {
		return ply/2 + 1, White
	}
	return (ply-1)/2 + 1, Black
}

func FullmovesToPly(fullmoves uint16, sideToMove Color) uint16 {
	if fullmoves == 0 {
		return 0
	}
	if sideToMove == White {
		return (fullmoves - 1) * 2
	}
	return (fullmoves-1)*2 + 1
}

func isWhitespaceOnly(s string) bool {
	if len(s) == 0 {
		return true
	}
	for _, r := range s {
		if !unicode.IsSpace(r) {
			return false
		}
	}
	return true
}

func StartingPosition() Position {
	pos, _ := ParseFEN(startingFEN)
	return pos
}
