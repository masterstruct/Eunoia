package board

import (
	"fmt"
	"strconv"
	"strings"
)

const resetColor = "\033[0m"

func (pos *Position) String() string {
	// If using VSCode, go to settings
	// and set this setting:
	// "terminal.integrated.minimumContrastRatio": 1,
	// this will fix the colors of the pretty board

	var sb strings.Builder

	// fen
	fmt.Fprintln(&sb, pos.FEN())

	for rank := Rank8; rank >= Rank1; rank-- {
		// ranks
		sb.WriteString(fgRGB(235, 160, 172))
		sb.WriteString(strconv.Itoa(int(rank + 1)))
		sb.WriteByte(' ')

		for file := FileA; file <= FileH; file++ {
			sq := NewSquare(file, rank)
			sqColor := sq.Color()

			sb.WriteString(squareBG(sqColor))

			piece, ok := pos.PieceOn(sq)
			if ok {
				sb.WriteString(pieceFG(piece.Color, sqColor))
				sb.WriteString(piece.PrettyString())
			} else {
				sb.WriteByte(' ')
			}

			sb.WriteByte(' ')
			sb.WriteString(resetColor)
		}

		sb.WriteByte('\n')
	}

	// files
	sb.WriteString(fgRGB(116, 199, 236))
	sb.WriteString("  a b c d e f g h\n")
	sb.WriteString(resetColor)
	sb.WriteByte('\n')

	// board state
	side := "white"
	if pos.SideToMove == Black {
		side = "black"
	}
	fullmoves, _ := PlyToFullmoves(pos.Ply)

	fmt.Fprintf(&sb, "%-13s%9s\n", "Side to move:", side)
	fmt.Fprintf(&sb, "%-16s%6s\n", "Castling rights:", pos.CastlingRights)
	fmt.Fprintf(&sb, "%-18s%4s\n", "En passant square:", pos.EnPassant)
	fmt.Fprintf(&sb, "%-13s%9d\n", "50 move rule:", pos.HalfmoveClock)
	fmt.Fprintf(&sb, "%-13s%9d\n", "Fullmoves:", fullmoves)
	fmt.Fprintf(&sb, "%-5s%17x\n", "Hash:", pos.Hash)

	return sb.String()
}

func squareBG(sqColor Color) string {
	if sqColor == Black {
		return bgRGB(152, 124, 180)
	}
	return bgRGB(229, 218, 241)
}

func pieceFG(pieceColor Color, sqColor Color) string {
	if pieceColor == Black {
		return fgRGB(49, 50, 68)
	}
	if sqColor == Black {
		// light piece 1
		return fgRGB(211, 183, 146)
	}
	// light piece 2
	return fgRGB(196, 173, 146)
}

func fgRGB(r, g, b int) string {
	// print pretty RGB colors for text foreground
	return fmt.Sprintf("\033[1;38;2;%d;%d;%dm", r, g, b)
}

func bgRGB(r, g, b int) string {
	// print pretty RBG colors for text background
	return fmt.Sprintf("\033[48;2;%d;%d;%dm", r, g, b)
}
