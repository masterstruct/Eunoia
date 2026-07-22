package board

import (
	"errors"
	"fmt"
	"strconv"
)

type Square int

const (
	A8, B8, C8, D8, E8, F8, G8, H8 Square = 56, 57, 58, 59, 60, 61, 62, 63
	A7, B7, C7, D7, E7, F7, G7, H7 Square = 48, 49, 50, 51, 52, 53, 54, 55
	A6, B6, C6, D6, E6, F6, G6, H6 Square = 40, 41, 42, 43, 44, 45, 46, 47
	A5, B5, C5, D5, E5, F5, G5, H5 Square = 32, 33, 34, 35, 36, 37, 38, 39
	A4, B4, C4, D4, E4, F4, G4, H4 Square = 24, 25, 26, 27, 28, 29, 30, 31
	A3, B3, C3, D3, E3, F3, G3, H3 Square = 16, 17, 18, 19, 20, 21, 22, 23
	A2, B2, C2, D2, E2, F2, G2, H2 Square = 8, 9, 10, 11, 12, 13, 14, 15
	A1, B1, C1, D1, E1, F1, G1, H1 Square = 0, 1, 2, 3, 4, 5, 6, 7

	NoSquare Square = -1
)

const (
	FileA = iota
	FileB
	FileC
	FileD
	FileE
	FileF
	FileG
	FileH
)

const (
	Rank1 = iota
	Rank2
	Rank3
	Rank4
	Rank5
	Rank6
	Rank7
	Rank8
)

var (
	ErrInvalidSquareLength = errors.New("square: string must be exactly 2 characters")
	ErrInvalidFile         = errors.New("square: file must be a letter between 'a' and 'h'")
	ErrInvalidRank         = errors.New("square: rank must be a digit between '1' and '8'")
)

func (sq Square) File() int {
	return int(sq) % 8
}

func (sq Square) Rank() int {
	return int(sq) / 8
}

func (sq Square) IsValid() bool {
	return sq >= A1 && sq <= H8
}

func (sq Square) String() string {
	f := sq.File()
	r := sq.Rank()
	// ascii manipulation
	letter := string('a' + byte(f))
	return fmt.Sprintf("%s%d", letter, r+1)
}

func ParseSquare(s string) (Square, error) {
	if len(s) != 2 {
		return NoSquare, fmt.Errorf("%w: %q", ErrInvalidSquareLength, s)
	}

	f := s[0] | 0x20 // fast lowercase conversion

	if f < 'a' || f > 'h' {
		return NoSquare, fmt.Errorf("%w: %q", ErrInvalidFile, s)
	}

	r, err := strconv.Atoi(string(s[1]))
	if err != nil || r < 1 || r > 8 {
		return NoSquare, fmt.Errorf("%w: %q", ErrInvalidRank, s)
	}

	return NewSquare(int(f-'a'), r-1), nil
}

func NewSquare(file, rank int) Square {
	return Square(rank*8 + file)
}
