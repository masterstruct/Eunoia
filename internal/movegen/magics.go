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

var RookMagics = [64]MagicEntry{
	{Mask: 0x000101010101017e, Magic: 0x9080001482204000, IndexBits: 12},
	{Mask: 0x000202020202027c, Magic: 0x4880200280104004, IndexBits: 11},
	{Mask: 0x000404040404047a, Magic: 0x2080100028822000, IndexBits: 11},
	{Mask: 0x0008080808080876, Magic: 0x0100080500100020, IndexBits: 11},
	{Mask: 0x001010101010106e, Magic: 0x0200080490206200, IndexBits: 11},
	{Mask: 0x002020202020205e, Magic: 0x0200108124082200, IndexBits: 11},
	{Mask: 0x004040404040403e, Magic: 0x0200012a82000408, IndexBits: 11},
	{Mask: 0x008080808080807e, Magic: 0x3200004600228504, IndexBits: 12},
	{Mask: 0x0001010101017e00, Magic: 0x0400800240102081, IndexBits: 11},
	{Mask: 0x0002020202027c00, Magic: 0x8800402010004000, IndexBits: 10},
	{Mask: 0x0004040404047a00, Magic: 0x0412002201814010, IndexBits: 10},
	{Mask: 0x0008080808087600, Magic: 0x0020801800100080, IndexBits: 10},
	{Mask: 0x0010101010106e00, Magic: 0x0021001004080100, IndexBits: 10},
	{Mask: 0x0020202020205e00, Magic: 0x2311001802040100, IndexBits: 10},
	{Mask: 0x0040404040403e00, Magic: 0x000e008102006804, IndexBits: 10},
	{Mask: 0x0080808080807e00, Magic: 0x0000800041000a80, IndexBits: 11},
	{Mask: 0x00010101017e0100, Magic: 0x0008a88000814000, IndexBits: 11},
	{Mask: 0x00020202027c0200, Magic: 0x40100c4000486000, IndexBits: 10},
	{Mask: 0x00040404047a0400, Magic: 0x0b0a020011406280, IndexBits: 10},
	{Mask: 0x0008080808760800, Magic: 0x000101003000a308, IndexBits: 10},
	{Mask: 0x00101010106e1000, Magic: 0x0204808008000400, IndexBits: 10},
	{Mask: 0x00202020205e2000, Magic: 0x2000808004000200, IndexBits: 10},
	{Mask: 0x00404040403e4000, Magic: 0x0802840010110208, IndexBits: 10},
	{Mask: 0x00808080807e8000, Magic: 0xa100c20001008144, IndexBits: 11},
	{Mask: 0x000101017e010100, Magic: 0x004000c080208001, IndexBits: 11},
	{Mask: 0x000202027c020200, Magic: 0x0c400260a0081000, IndexBits: 10},
	{Mask: 0x000404047a040400, Magic: 0xa080200500410110, IndexBits: 10},
	{Mask: 0x0008080876080800, Magic: 0x0008028080100008, IndexBits: 10},
	{Mask: 0x001010106e101000, Magic: 0x2102010600102028, IndexBits: 10},
	{Mask: 0x002020205e202000, Magic: 0x0805000300481400, IndexBits: 10},
	{Mask: 0x004040403e404000, Magic: 0x0100082400021009, IndexBits: 10},
	{Mask: 0x008080807e808000, Magic: 0x00c10d020008488c, IndexBits: 11},
	{Mask: 0x0001017e01010100, Magic: 0x0000854014800020, IndexBits: 11},
	{Mask: 0x0002027c02020200, Magic: 0x0001018202002040, IndexBits: 10},
	{Mask: 0x0004047a04040400, Magic: 0x6008852001801000, IndexBits: 10},
	{Mask: 0x0008087608080800, Magic: 0x0002002012000842, IndexBits: 10},
	{Mask: 0x0010106e10101000, Magic: 0x0182080005001100, IndexBits: 10},
	{Mask: 0x0020205e20202000, Magic: 0x0404204008011004, IndexBits: 10},
	{Mask: 0x0040403e40404000, Magic: 0x280310080c000102, IndexBits: 10},
	{Mask: 0x0080807e80808000, Magic: 0x0055801040802100, IndexBits: 11},
	{Mask: 0x00017e0101010100, Magic: 0x0040204000808002, IndexBits: 11},
	{Mask: 0x00027c0202020200, Magic: 0x0810084020004000, IndexBits: 10},
	{Mask: 0x00047a0404040400, Magic: 0x8050001020008080, IndexBits: 10},
	{Mask: 0x0008760808080800, Magic: 0x0101002490010008, IndexBits: 10},
	{Mask: 0x00106e1010101000, Magic: 0x4400440008008080, IndexBits: 10},
	{Mask: 0x00205e2020202000, Magic: 0x0012000400808042, IndexBits: 10},
	{Mask: 0x00403e4040404000, Magic: 0x8010024801840010, IndexBits: 10},
	{Mask: 0x00807e8080808000, Magic: 0x0480040882420003, IndexBits: 11},
	{Mask: 0x007e010101010100, Magic: 0x610a214201018200, IndexBits: 11},
	{Mask: 0x007c020202020200, Magic: 0x3000912002400080, IndexBits: 10},
	{Mask: 0x007a040404040400, Magic: 0x0001004032200500, IndexBits: 10},
	{Mask: 0x0076080808080800, Magic: 0x5826004020289200, IndexBits: 10},
	{Mask: 0x006e101010101000, Magic: 0x44030c0128008080, IndexBits: 10},
	{Mask: 0x005e202020202000, Magic: 0x084200f008044200, IndexBits: 10},
	{Mask: 0x003e404040404000, Magic: 0xa203008406000700, IndexBits: 10},
	{Mask: 0x007e808080808000, Magic: 0x000100008200c100, IndexBits: 11},
	{Mask: 0x7e01010101010100, Magic: 0x5002008061034032, IndexBits: 12},
	{Mask: 0x7c02020202020200, Magic: 0x14054302008112a2, IndexBits: 11},
	{Mask: 0x7a04040404040400, Magic: 0x901e12004020808a, IndexBits: 11},
	{Mask: 0x7608080808080800, Magic: 0x8000049000190021, IndexBits: 11},
	{Mask: 0x6e10101010101000, Magic: 0x0102000824102002, IndexBits: 11},
	{Mask: 0x5e20202020202000, Magic: 0x0882001028040102, IndexBits: 11},
	{Mask: 0x3e40404040404000, Magic: 0x00802f1000820824, IndexBits: 11},
	{Mask: 0x7e80808080808000, Magic: 0x0100110401c02782, IndexBits: 12},
}

var BishopMagics = [64]MagicEntry{
	{Mask: 0x0040201008040200, Magic: 0x4003200808010200, IndexBits: 12},
	{Mask: 0x0000402010080400, Magic: 0xa400400810441121, IndexBits: 11},
	{Mask: 0x0000004020100a00, Magic: 0x0002204100091400, IndexBits: 11},
	{Mask: 0x0000000040221400, Magic: 0x0001412004020000, IndexBits: 11},
	{Mask: 0x0000000002442800, Magic: 0x2006440080801182, IndexBits: 11},
	{Mask: 0x0000000204085000, Magic: 0x2000020840040020, IndexBits: 11},
	{Mask: 0x0000020408102000, Magic: 0x00080b0282802001, IndexBits: 11},
	{Mask: 0x0002040810204000, Magic: 0x8000208020114800, IndexBits: 12},
	{Mask: 0x0020100804020000, Magic: 0x102020204006002e, IndexBits: 11},
	{Mask: 0x0040201008040000, Magic: 0x2802004404040210, IndexBits: 10},
	{Mask: 0x00004020100a0000, Magic: 0x40400021001a0000, IndexBits: 10},
	{Mask: 0x0000004022140000, Magic: 0x1a00004222602002, IndexBits: 10},
	{Mask: 0x0000000244280000, Magic: 0x8b02001004040000, IndexBits: 10},
	{Mask: 0x0000020408500000, Magic: 0x2000101408110008, IndexBits: 10},
	{Mask: 0x0002040810200000, Magic: 0x1421828080080098, IndexBits: 10},
	{Mask: 0x0004081020400000, Magic: 0x0000000c04008805, IndexBits: 11},
	{Mask: 0x0010080402000200, Magic: 0x8180316000108020, IndexBits: 11},
	{Mask: 0x0020100804000400, Magic: 0x0001000400c0060a, IndexBits: 10},
	{Mask: 0x004020100a000a00, Magic: 0x4010502c02000821, IndexBits: 10},
	{Mask: 0x0000402214001400, Magic: 0x2052040412000100, IndexBits: 10},
	{Mask: 0x0000024428002800, Magic: 0x0804000014446068, IndexBits: 10},
	{Mask: 0x0002040850005000, Magic: 0x0080081004100204, IndexBits: 10},
	{Mask: 0x0004081020002000, Magic: 0x3004008100401000, IndexBits: 10},
	{Mask: 0x0008102040004000, Magic: 0x8002002008400802, IndexBits: 11},
	{Mask: 0x0008040200020400, Magic: 0x9101400801220280, IndexBits: 11},
	{Mask: 0x0010080400040800, Magic: 0x2100112401100040, IndexBits: 10},
	{Mask: 0x0020100a000a1000, Magic: 0x2120220060840501, IndexBits: 10},
	{Mask: 0x0040221400142200, Magic: 0x0060020063005004, IndexBits: 10},
	{Mask: 0x0002442800284400, Magic: 0x0104040000401040, IndexBits: 10},
	{Mask: 0x0004085000500800, Magic: 0x10400849101c1000, IndexBits: 10},
	{Mask: 0x0008102000201000, Magic: 0x05004002008a010c, IndexBits: 10},
	{Mask: 0x0010204000402000, Magic: 0x0040402080806258, IndexBits: 11},
	{Mask: 0x0004020002040800, Magic: 0x00010010000080a0, IndexBits: 11},
	{Mask: 0x0008040004081000, Magic: 0x4520622000008608, IndexBits: 10},
	{Mask: 0x00100a000a102000, Magic: 0x0522004020110208, IndexBits: 10},
	{Mask: 0x0022140014224000, Magic: 0x0402020021001042, IndexBits: 10},
	{Mask: 0x0044280028440200, Magic: 0x809902000c049040, IndexBits: 10},
	{Mask: 0x0008500050080400, Magic: 0x004150001805000a, IndexBits: 10},
	{Mask: 0x0010200020100800, Magic: 0x1801001a1020b804, IndexBits: 10},
	{Mask: 0x0020400040201000, Magic: 0x1020400428108204, IndexBits: 11},
	{Mask: 0x0002000204081000, Magic: 0x2380022110002404, IndexBits: 11},
	{Mask: 0x0004000408102000, Magic: 0x2008802222010200, IndexBits: 10},
	{Mask: 0x000a000a10204000, Magic: 0x0010b04230100808, IndexBits: 10},
	{Mask: 0x0014001422400000, Magic: 0x0520004404300040, IndexBits: 10},
	{Mask: 0x0028002844020000, Magic: 0xcd20842028200c00, IndexBits: 10},
	{Mask: 0x0050005008040200, Magic: 0xa002061100480034, IndexBits: 10},
	{Mask: 0x0020002010080400, Magic: 0x1481002008100144, IndexBits: 10},
	{Mask: 0x0040004020100800, Magic: 0x2103201002000900, IndexBits: 11},
	{Mask: 0x0000020408102000, Magic: 0x020ac02801024000, IndexBits: 11},
	{Mask: 0x0000040810204000, Magic: 0x0002202104084001, IndexBits: 10},
	{Mask: 0x00000a1020400000, Magic: 0x4000614002822400, IndexBits: 10},
	{Mask: 0x0000142240000000, Magic: 0x0021000040f21800, IndexBits: 10},
	{Mask: 0x0000284402000000, Magic: 0x0010304800205c25, IndexBits: 10},
	{Mask: 0x0000500804020000, Magic: 0x7000412004010400, IndexBits: 10},
	{Mask: 0x0000201008040200, Magic: 0x0000441000308000, IndexBits: 10},
	{Mask: 0x0000402010080400, Magic: 0x4020420008020008, IndexBits: 11},
	{Mask: 0x0002040810204000, Magic: 0x0088180404010008, IndexBits: 12},
	{Mask: 0x0004081020400000, Magic: 0x8000000e4a08c200, IndexBits: 11},
	{Mask: 0x000a102040000000, Magic: 0x0022020804000230, IndexBits: 11},
	{Mask: 0x0014224000000000, Magic: 0x0004210220408080, IndexBits: 11},
	{Mask: 0x0028440200000000, Magic: 0x5001008800212020, IndexBits: 11},
	{Mask: 0x0050080402000000, Magic: 0x2200080800300020, IndexBits: 11},
	{Mask: 0x0020100804020000, Magic: 0x0004020424018001, IndexBits: 11},
	{Mask: 0x0040201008040200, Magic: 0x1000801041220002, IndexBits: 12},
}

var RookMoves [64][]board.Bitboard
var BishopMoves [64][]board.Bitboard

type MagicEntry struct {
	Mask      board.Bitboard
	Magic     uint64
	IndexBits uint8
}

func MagicIndex(entry *MagicEntry, occupied board.Bitboard) int {
	occupied &= entry.Mask
	hash := uint64(occupied) * entry.Magic
	return int(hash >> (64 - entry.IndexBits))
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
	initMagics(true, &RookMagics, &RookMoves, "rook")
	initMagics(false, &BishopMagics, &BishopMoves, "bishop")
}

func initMagics(isRook bool, magics *[64]MagicEntry, moves *[64][]board.Bitboard, name string) {
	for sq := board.A1; sq <= board.H8; sq++ {
		entry := &magics[sq]

		table, ok := TryMakeTable(isRook, sq, entry)
		if !ok {
			panic(fmt.Sprintf("Bad %s magics! Square %v", name, sq))
		}
		moves[sq] = table
	}
}

func TryMakeTable(isRook bool, sq board.Square, entry *MagicEntry) ([]board.Bitboard, bool) {
	tableSize := 1 << entry.IndexBits
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
