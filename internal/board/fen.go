package board

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

const startingFEN = "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1"

var (
	ErrInvalidFieldCount     = errors.New("fen: expected 4 to 6 space-separated fields (piece placement, side to move, castling, en passant, [halfmove clock], [fullmove number])")
	ErrInvalidRankCount      = errors.New("fen: piece placement must have exactly 8 ranks separated by '/'")
	ErrInvalidPieceChar      = errors.New("fen: piece placement contains an unrecognized character")
	ErrInvalidRankDigit      = errors.New("fen: empty-square digit must be between 1 and 8")
	ErrInvalidRankLength     = errors.New("fen: each rank must total exactly 8 squares")
	ErrInvalidSideToMove     = errors.New("fen: side to move must be 'w' or 'b'")
	ErrInvalidHalfmoveClock  = errors.New("fen: halfmove clock must be a non-negative integer between 0 and 255")
	ErrInvalidFullmoveNumber = errors.New("fen: fullmove number must be a positive integer starting from 1")
	ErrInvalidKingCount      = errors.New("fen: each side must have exactly one king")
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
		if len(r) > 8 {
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

		// detect numbers and skip that many squares
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
				// lowercase => Black
				color = Black
			} else {
				color = White
			}

			// place piece on the board
			pos.PlacePiece(NewPiece(pt, color), NewSquare(file, rank))

			file += 1
			continue
		}
		// character isn't a number or a valid piece
		return pos, fmt.Errorf("%w: %q", ErrInvalidPieceChar, fen)
	}

	// side to move
	sideToMove := ParseColor(splits[1][0])
	if sideToMove == NoColor {
		return pos, fmt.Errorf("%w: %q", ErrInvalidSideToMove, splits[1][0])
	}
	pos.SideToMove = sideToMove

	// king count
	whiteKingBB := pos.PieceBB(WhiteKing)
	blackKingBB := pos.PieceBB(BlackKing)
	if whiteKingBB.CountBits() != 1 || blackKingBB.CountBits() != 1 {
		return pos, fmt.Errorf("%w: %q", ErrInvalidKingCount, fen)
	}

	// castling rights
	rights, err := ParseCastlingRights(splits[2], whiteKingBB.LSB(), blackKingBB.LSB())
	if err != nil {
		return pos, err
	}
	pos.CastlingRights = rights

	// en passant
	sq, err := ParseSquare(splits[3])
	if err != nil {
		return pos, fmt.Errorf("%w: %q", err, splits[3])
	}
	pos.EnPassant = sq

	// halfmove clock
	if n >= 5 {
		halfMove, err := strconv.Atoi(splits[4])
		if err != nil || halfMove < 0 || halfMove > 100 {
			return pos, fmt.Errorf("%w: %q", ErrInvalidHalfmoveClock, splits[4])
		}
		pos.HalfmoveClock = uint8(halfMove)
	}

	// fullmove counter
	if n >= 6 {
		fullmoves, err := strconv.Atoi(splits[5])
		if err != nil || fullmoves <= 0 || fullmoves > 10000 {
			return pos, fmt.Errorf("%w: %q", ErrInvalidFullmoveNumber, splits[5])
		}
		pos.Ply = FullmovesToPly(uint16(fullmoves), pos.SideToMove)
	}

	return pos, nil
}

func (pos Position) FEN() string {
	var sb strings.Builder
	var rankBB Bitboard
	var skip int

	occupied := pos.Occupied()

	for rank := Rank8; rank >= Rank1; rank-- {
		rankBB = occupied & RankBB[rank]
		if rankBB == EmptyBB {
			// skip empty rank
			sb.WriteString("8")
			if rank != Rank1 {
				sb.WriteByte('/')
			}
			skip = 0
			continue
		}

		for file := FileA; file <= FileH; file++ {
			sq := NewSquare(file, rank)

			if !occupied.IsBitSet(sq) {
				skip++
				continue
			}
			if skip > 0 {
				sb.WriteByte(byte('0' + skip))
				skip = 0
			}

			piece, _ := pos.PieceOn(sq)
			sb.WriteByte(piece.String())
		}

		// reached end of rank
		if skip > 0 {
			sb.WriteByte(byte('0' + skip))
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

func StartingPosition() Position {
	pos, _ := ParseFEN(startingFEN)
	return pos
}
