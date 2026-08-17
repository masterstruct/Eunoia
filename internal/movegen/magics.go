package movegen

import (
	"fmt"

	"github.com/masterstruct/Eunoia/internal/board"
)

func RookMask(sq board.Square) board.Bitboard {
	file := sq.File()
	rank := sq.Rank()
	var mask board.Bitboard

	for f := board.FileB; f < board.FileH; f++ {
		mask.SetBit(board.NewSquare(f, rank))
	}

	for r := board.Rank2; r < board.Rank8; r++ {
		mask.SetBit(board.NewSquare(file, r))
	}

	mask.ClearBit(sq)
	return mask
}

func BishopMask(sq board.Square) board.Bitboard {
	var mask board.Bitboard
	mask = BishopAttacksSlow(sq, board.EmptyBB)
	return mask &^ board.EdgesBB
}

type MagicEntry struct {
	Mask   board.Bitboard
	Magic  uint64
	Shift  uint8
	Offset int
}

func MagicIndex(entry *MagicEntry, occupied board.Bitboard) int {
	hash := uint64(occupied&entry.Mask) * entry.Magic
	return int(hash>>entry.Shift) + entry.Offset
}

// iterate over all subsets of a bitboard
func Subsets(mask board.Bitboard) func(yield func(board.Bitboard) bool) {
	return func(yield func(board.Bitboard) bool) {
		subset := mask
		for {
			if !yield(subset) {
				return
			}
			if subset == 0 {
				return
			}
			subset = (subset - 1) & mask
		}
	}
}

func init() {
	rookOffset := 0
	bishopOffset := 0

	for sq := board.A1; sq <= board.H8; sq++ {
		rookEntry := &RookMagics[sq]
		rookEntry.Mask = RookMask(sq)
		rookTable, ok := TryMakeTable(true, sq, rookEntry)
		if !ok {
			panic(fmt.Sprintf("Bad rook magics! Square %v", sq))
		}
		rookEntry.Offset = rookOffset
		if rookOffset+len(rookTable) > len(RookMoves) {
			panic(fmt.Sprintf("Rook table overflow at square %v", sq))
		}
		copy(RookMoves[rookOffset:], rookTable)
		rookOffset += len(rookTable)

		bishopEntry := &BishopMagics[sq]
		bishopEntry.Mask = BishopMask(sq)
		bishopTable, ok := TryMakeTable(false, sq, bishopEntry)
		if !ok {
			panic(fmt.Sprintf("Bad bishop magics! Square %v", sq))
		}
		bishopEntry.Offset = bishopOffset
		if bishopOffset+len(bishopTable) > len(BishopMoves) {
			panic(fmt.Sprintf("Bishop table overflow at square %v", sq))
		}
		copy(BishopMoves[bishopOffset:], bishopTable)
		bishopOffset += len(bishopTable)
	}
}

func TryMakeTable(isRook bool, sq board.Square, entry *MagicEntry) ([]board.Bitboard, bool) {
	tableSize := 1 << (64 - entry.Shift)
	table := make([]board.Bitboard, tableSize)
	used := make([]bool, tableSize)

	var moves board.Bitboard
	for blockers := range Subsets(entry.Mask) {
		if isRook {
			moves = RookAttacksSlow(sq, blockers)
		} else {
			moves = BishopAttacksSlow(sq, blockers)
		}

		index := MagicIndex(entry, blockers)
		if !used[index] {
			used[index] = true
			table[index] = moves
		} else if table[index] != moves {
			return []board.Bitboard{}, false
		}
	}
	return table, true
}

const RookTableSize = 102400
const BishopTableSize = 5248

var RookMoves [RookTableSize]board.Bitboard
var BishopMoves [BishopTableSize]board.Bitboard

var RookMagics = [64]MagicEntry{
	{Magic: 0x0180008118400122, Shift: 52},
	{Magic: 0x0380200140001284, Shift: 53},
	{Magic: 0x0500081302200040, Shift: 53},
	{Magic: 0x1880142800900180, Shift: 53},
	{Magic: 0x0200020020100804, Shift: 53},
	{Magic: 0x0200020008100104, Shift: 53},
	{Magic: 0x2280028005000200, Shift: 53},
	{Magic: 0x01000100002c4982, Shift: 52},
	{Magic: 0x3280800020400080, Shift: 53},
	{Magic: 0x0000804000200280, Shift: 54},
	{Magic: 0x00110010c1042008, Shift: 54},
	{Magic: 0x04120040101a2200, Shift: 54},
	{Magic: 0x0000800802808c00, Shift: 54},
	{Magic: 0x0081000401000882, Shift: 54},
	{Magic: 0x0304000102080410, Shift: 54},
	{Magic: 0x0182800041000080, Shift: 53},
	{Magic: 0x0040078000402981, Shift: 53},
	{Magic: 0xc03002400040a008, Shift: 54},
	{Magic: 0x0010008018200180, Shift: 54},
	{Magic: 0x04010100200a1003, Shift: 54},
	{Magic: 0x0800808004000800, Shift: 54},
	{Magic: 0x3006008004008280, Shift: 54},
	{Magic: 0x1802140050020108, Shift: 54},
	{Magic: 0x00100e0024044881, Shift: 53},
	{Magic: 0x00800040c0002000, Shift: 53},
	{Magic: 0x4001002100400080, Shift: 54},
	{Magic: 0x0046100080200080, Shift: 54},
	{Magic: 0x8200090100500020, Shift: 54},
	{Magic: 0x0210080100250030, Shift: 54},
	{Magic: 0x0002000280800400, Shift: 54},
	{Magic: 0x0118980400100102, Shift: 54},
	{Magic: 0x08020c8200004924, Shift: 53},
	{Magic: 0x1044804000800035, Shift: 53},
	{Magic: 0x5104802000804001, Shift: 54},
	{Magic: 0x0000802000801004, Shift: 54},
	{Magic: 0x0088008408801000, Shift: 54},
	{Magic: 0x0100040180801801, Shift: 54},
	{Magic: 0x1d04004100400200, Shift: 54},
	{Magic: 0x0008011024008802, Shift: 54},
	{Magic: 0x008010804a000405, Shift: 53},
	{Magic: 0x0000204010808000, Shift: 53},
	{Magic: 0x040a201000424000, Shift: 54},
	{Magic: 0x0280824200220010, Shift: 54},
	{Magic: 0xc109016030010018, Shift: 54},
	{Magic: 0x4630100801010024, Shift: 54},
	{Magic: 0xa54c008002008004, Shift: 54},
	{Magic: 0x0000302803040002, Shift: 54},
	{Magic: 0x0880508400420011, Shift: 53},
	{Magic: 0x4002608041060a00, Shift: 53},
	{Magic: 0x0004400020810100, Shift: 54},
	{Magic: 0x0000200080100080, Shift: 54},
	{Magic: 0x0003411122016a00, Shift: 54},
	{Magic: 0x01000c0080080080, Shift: 54},
	{Magic: 0x1202010810048200, Shift: 54},
	{Magic: 0x641011108a280400, Shift: 54},
	{Magic: 0x0002090402448600, Shift: 53},
	{Magic: 0x0848804122010032, Shift: 52},
	{Magic: 0x0800102100804007, Shift: 53},
	{Magic: 0x0020002100440831, Shift: 53},
	{Magic: 0x0000200501900029, Shift: 53},
	{Magic: 0x0191000800044251, Shift: 53},
	{Magic: 0x5012002450110832, Shift: 53},
	{Magic: 0x0601108108100204, Shift: 53},
	{Magic: 0x2800908100a404c2, Shift: 52},
}

var BishopMagics = [64]MagicEntry{
	{Magic: 0x0408011002044101, Shift: 58},
	{Magic: 0x004210150a028024, Shift: 59},
	{Magic: 0x8090010071042000, Shift: 59},
	{Magic: 0xc1080a0028003200, Shift: 59},
	{Magic: 0x0044042000040800, Shift: 59},
	{Magic: 0x0822081404008000, Shift: 59},
	{Magic: 0x1c02014c12400300, Shift: 59},
	{Magic: 0x015204808c012004, Shift: 58},
	{Magic: 0x0c10040810040480, Shift: 59},
	{Magic: 0x0000140808044080, Shift: 59},
	{Magic: 0x0088110922120008, Shift: 59},
	{Magic: 0x88a0c82180200008, Shift: 59},
	{Magic: 0x2064040420200000, Shift: 59},
	{Magic: 0x006000880440018a, Shift: 59},
	{Magic: 0x201204a818023020, Shift: 59},
	{Magic: 0x2140384402011040, Shift: 59},
	{Magic: 0x00600449200c4080, Shift: 59},
	{Magic: 0x0004200204880200, Shift: 59},
	{Magic: 0x0010282800811010, Shift: 57},
	{Magic: 0x8508000282004100, Shift: 57},
	{Magic: 0x2004041180a00002, Shift: 57},
	{Magic: 0x0800400208200402, Shift: 57},
	{Magic: 0x8000a00508151001, Shift: 59},
	{Magic: 0x0184804024230800, Shift: 59},
	{Magic: 0x4804044040100401, Shift: 59},
	{Magic: 0x0124040802101400, Shift: 59},
	{Magic: 0x403c900008084030, Shift: 57},
	{Magic: 0x80400400054100a0, Shift: 55},
	{Magic: 0x00a1020144008400, Shift: 55},
	{Magic: 0x08d10100a2100082, Shift: 57},
	{Magic: 0x0064004604220200, Shift: 59},
	{Magic: 0x00052201892a0100, Shift: 59},
	{Magic: 0x08220820288ca000, Shift: 59},
	{Magic: 0x100210020890c200, Shift: 59},
	{Magic: 0x0004002400480040, Shift: 57},
	{Magic: 0x00420042409c0100, Shift: 55},
	{Magic: 0x0018c10040040040, Shift: 55},
	{Magic: 0x0086058200290800, Shift: 57},
	{Magic: 0x4072080200010ac1, Shift: 59},
	{Magic: 0x0004048200002100, Shift: 59},
	{Magic: 0x4006088c0441c052, Shift: 59},
	{Magic: 0x400a080434100220, Shift: 59},
	{Magic: 0x000600202410080a, Shift: 57},
	{Magic: 0x0200c86254000800, Shift: 57},
	{Magic: 0x01000208a4000a00, Shift: 57},
	{Magic: 0x000a424042000d00, Shift: 57},
	{Magic: 0x08a8100101403200, Shift: 59},
	{Magic: 0x180400840bc00312, Shift: 59},
	{Magic: 0x0084d42404408200, Shift: 59},
	{Magic: 0x801200c14c100000, Shift: 59},
	{Magic: 0x5000811441100000, Shift: 59},
	{Magic: 0x0444000284041115, Shift: 59},
	{Magic: 0x0020812c208a0400, Shift: 59},
	{Magic: 0x000c081025020200, Shift: 59},
	{Magic: 0x0009200914050015, Shift: 59},
	{Magic: 0x00120a040d060000, Shift: 59},
	{Magic: 0x0000402098201013, Shift: 58},
	{Magic: 0x2800250141100848, Shift: 59},
	{Magic: 0xc144000206010428, Shift: 59},
	{Magic: 0x2108030024420880, Shift: 59},
	{Magic: 0x8408082c10021208, Shift: 59},
	{Magic: 0x0904001202500500, Shift: 59},
	{Magic: 0x0010100430308200, Shift: 59},
	{Magic: 0x2088200402893100, Shift: 58},
}
