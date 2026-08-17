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
	{Mask: 0x000101010101017e, Magic: 0x1280008230400020, Shift: 52},
	{Mask: 0x000202020202027c, Magic: 0x0840002008411002, Shift: 53},
	{Mask: 0x000404040404047a, Magic: 0x8a001181400a0020, Shift: 53},
	{Mask: 0x0008080808080876, Magic: 0x2080100080360800, Shift: 53},
	{Mask: 0x001010101010106e, Magic: 0x0200141020880200, Shift: 53},
	{Mask: 0x002020202020205e, Magic: 0x83000400184a0300, Shift: 53},
	{Mask: 0x004040404040403e, Magic: 0x0080008022002500, Shift: 53},
	{Mask: 0x008080808080807e, Magic: 0x8e00010400802452, Shift: 52},
	{Mask: 0x0001010101017e00, Magic: 0x424080004008a080, Shift: 53},
	{Mask: 0x0002020202027c00, Magic: 0x0000402000401000, Shift: 54},
	{Mask: 0x0004040404047a00, Magic: 0x0481802000100081, Shift: 54},
	{Mask: 0x0008080808087600, Magic: 0x1000803002800802, Shift: 54},
	{Mask: 0x0010101010106e00, Magic: 0x000a00100600a008, Shift: 54},
	{Mask: 0x0020202020205e00, Magic: 0x2c0200480a005490, Shift: 54},
	{Mask: 0x0040404040403e00, Magic: 0xa321002200430004, Shift: 54},
	{Mask: 0x0080808080807e00, Magic: 0x0102800900006080, Shift: 53},
	{Mask: 0x00010101017e0100, Magic: 0x0021808000400220, Shift: 53},
	{Mask: 0x00020202027c0200, Magic: 0x0010094003200042, Shift: 54},
	{Mask: 0x00040404047a0400, Magic: 0x0801c20010802200, Shift: 54},
	{Mask: 0x0008080808760800, Magic: 0x404103001000a288, Shift: 54},
	{Mask: 0x00101010106e1000, Magic: 0x0009010014100800, Shift: 54},
	{Mask: 0x00202020205e2000, Magic: 0x0020480134402010, Shift: 54},
	{Mask: 0x00404040403e4000, Magic: 0x0030040005020830, Shift: 54},
	{Mask: 0x00808080807e8000, Magic: 0x0012060000408401, Shift: 53},
	{Mask: 0x000101017e010100, Magic: 0x0080810200204200, Shift: 53},
	{Mask: 0x000202027c020200, Magic: 0x24a0400440201001, Shift: 54},
	{Mask: 0x000404047a040400, Magic: 0x2001200280100882, Shift: 54},
	{Mask: 0x0008080876080800, Magic: 0x4011002100081000, Shift: 54},
	{Mask: 0x001010106e101000, Magic: 0x0086000a00100c20, Shift: 54},
	{Mask: 0x002020205e202000, Magic: 0x0002040080800200, Shift: 54},
	{Mask: 0x004040403e404000, Magic: 0x000490040001a802, Shift: 54},
	{Mask: 0x008080807e808000, Magic: 0x0001000100048042, Shift: 53},
	{Mask: 0x0001017e01010100, Magic: 0x0020400281800028, Shift: 53},
	{Mask: 0x0002027c02020200, Magic: 0x0009028022004205, Shift: 54},
	{Mask: 0x0004047a04040400, Magic: 0x0700802008801000, Shift: 54},
	{Mask: 0x0008087608080800, Magic: 0x0400180080801000, Shift: 54},
	{Mask: 0x0010106e10101000, Magic: 0x0c12800802801400, Shift: 54},
	{Mask: 0x0020205e20202000, Magic: 0x0002001022000804, Shift: 54},
	{Mask: 0x0040403e40404000, Magic: 0x2000081004004601, Shift: 54},
	{Mask: 0x0080807e80808000, Magic: 0x2008488042001401, Shift: 53},
	{Mask: 0x00017e0101010100, Magic: 0x2900400080228000, Shift: 53},
	{Mask: 0x00027c0202020200, Magic: 0x0883810200220040, Shift: 54},
	{Mask: 0x00047a0404040400, Magic: 0x2020080010004040, Shift: 54},
	{Mask: 0x0008760808080800, Magic: 0x8401001000210028, Shift: 54},
	{Mask: 0x00106e1010101000, Magic: 0x0105000800110004, Shift: 54},
	{Mask: 0x00205e2020202000, Magic: 0x0002005104160008, Shift: 54},
	{Mask: 0x00403e4040404000, Magic: 0x1000820001008080, Shift: 54},
	{Mask: 0x00807e8080808000, Magic: 0x40000902a0420004, Shift: 53},
	{Mask: 0x007e010101010100, Magic: 0x0002801041082500, Shift: 53},
	{Mask: 0x007c020202020200, Magic: 0x0081018040002100, Shift: 54},
	{Mask: 0x007a040404040400, Magic: 0x000a084022819200, Shift: 54},
	{Mask: 0x0076080808080800, Magic: 0x1006021042082200, Shift: 54},
	{Mask: 0x006e101010101000, Magic: 0x0004008044080080, Shift: 54},
	{Mask: 0x005e202020202000, Magic: 0x00008c0080120080, Shift: 54},
	{Mask: 0x003e404040404000, Magic: 0x0c08500801160400, Shift: 54},
	{Mask: 0x007e808080808000, Magic: 0x00210020c2820100, Shift: 53},
	{Mask: 0x7e01010101010100, Magic: 0x0010126080004901, Shift: 52},
	{Mask: 0x7c02020202020200, Magic: 0x0091084200802112, Shift: 53},
	{Mask: 0x7a04040404040400, Magic: 0x028200c02084104a, Shift: 53},
	{Mask: 0x7608080808080800, Magic: 0x0804082010010035, Shift: 53},
	{Mask: 0x6e10101010101000, Magic: 0x00060024d0202802, Shift: 53},
	{Mask: 0x5e20202020202000, Magic: 0xa0020030610c0812, Shift: 53},
	{Mask: 0x3e40404040404000, Magic: 0x000010020801188c, Shift: 53},
	{Mask: 0x7e80808080808000, Magic: 0x000401062400408a, Shift: 52},
}

var BishopMagics = [64]MagicEntry{
	{Mask: 0x0040201008040200, Magic: 0x0000910228210020, Shift: 52},
	{Mask: 0x0000402010080400, Magic: 0x0102020050211920, Shift: 53},
	{Mask: 0x0000004020100a00, Magic: 0x0080448010400400, Shift: 53},
	{Mask: 0x0000000040221400, Magic: 0x000a824422824042, Shift: 53},
	{Mask: 0x0000000002442800, Magic: 0x2000092100205000, Shift: 53},
	{Mask: 0x0000000204085000, Magic: 0x1804020008200100, Shift: 53},
	{Mask: 0x0000020408102000, Magic: 0x0001002050012008, Shift: 53},
	{Mask: 0x0002040810204000, Magic: 0x0041320090072084, Shift: 52},
	{Mask: 0x0020100804020000, Magic: 0x0000048510834001, Shift: 53},
	{Mask: 0x0040201008040000, Magic: 0x0200011002001240, Shift: 54},
	{Mask: 0x00004020100a0000, Magic: 0x2a10040004082212, Shift: 54},
	{Mask: 0x0000004022140000, Magic: 0x0082603044019202, Shift: 54},
	{Mask: 0x0000000244280000, Magic: 0x1110003882810084, Shift: 54},
	{Mask: 0x0000020408500000, Magic: 0x0044041224048040, Shift: 54},
	{Mask: 0x0002040810200000, Magic: 0x24080148170080a6, Shift: 54},
	{Mask: 0x0004081020400000, Magic: 0x1000000920054280, Shift: 53},
	{Mask: 0x0010080402000200, Magic: 0x0430100010004240, Shift: 53},
	{Mask: 0x0020100804000400, Magic: 0x81001100a0e10008, Shift: 54},
	{Mask: 0x004020100a000a00, Magic: 0x0002200100064200, Shift: 54},
	{Mask: 0x0000402214001400, Magic: 0x21c1400541c20020, Shift: 54},
	{Mask: 0x0000024428002800, Magic: 0x2095409438a00800, Shift: 54},
	{Mask: 0x0002040850005000, Magic: 0x0000904801011500, Shift: 54},
	{Mask: 0x0004081020002000, Magic: 0x004004090e402001, Shift: 54},
	{Mask: 0x0008102040004000, Magic: 0x0400080040180100, Shift: 53},
	{Mask: 0x0008040200020400, Magic: 0x0041050040140180, Shift: 53},
	{Mask: 0x0010080400040800, Magic: 0x000422088a002064, Shift: 54},
	{Mask: 0x0020100a000a1000, Magic: 0x60023814100c8010, Shift: 54},
	{Mask: 0x0040221400142200, Magic: 0x8240208040200850, Shift: 54},
	{Mask: 0x0002442800284400, Magic: 0x202081001480404c, Shift: 54},
	{Mask: 0x0004085000500800, Magic: 0x8282000998080104, Shift: 54},
	{Mask: 0x0008102000201000, Magic: 0x0044445005000c20, Shift: 54},
	{Mask: 0x0010204000402000, Magic: 0x4041000182820020, Shift: 53},
	{Mask: 0x0004020002040800, Magic: 0x820142005011400c, Shift: 53},
	{Mask: 0x0008040004081000, Magic: 0x82210120080014c4, Shift: 54},
	{Mask: 0x00100a000a102000, Magic: 0x8401000800809804, Shift: 54},
	{Mask: 0x0022140014224000, Magic: 0x000444010000a002, Shift: 54},
	{Mask: 0x0044280028440200, Magic: 0x1000a06014020801, Shift: 54},
	{Mask: 0x0008500050080400, Magic: 0x4420900044090100, Shift: 54},
	{Mask: 0x0010200020100800, Magic: 0x40220900a012401a, Shift: 54},
	{Mask: 0x0020400040201000, Magic: 0x0008800108018008, Shift: 53},
	{Mask: 0x0002000204081000, Magic: 0x0060081010000100, Shift: 53},
	{Mask: 0x0004000408102000, Magic: 0x0020888580103000, Shift: 54},
	{Mask: 0x000a000a10204000, Magic: 0x010d009070910110, Shift: 54},
	{Mask: 0x0014001422400000, Magic: 0x0414008348010401, Shift: 54},
	{Mask: 0x0028002844020000, Magic: 0x020405a902000070, Shift: 54},
	{Mask: 0x0050005008040200, Magic: 0x41202000c1401010, Shift: 54},
	{Mask: 0x0020002010080400, Magic: 0x0000200881201020, Shift: 54},
	{Mask: 0x0040004020100800, Magic: 0x8011a02002800080, Shift: 53},
	{Mask: 0x0000020408102000, Magic: 0x000020c408400380, Shift: 53},
	{Mask: 0x0000040810204000, Magic: 0x0420010804500001, Shift: 54},
	{Mask: 0x00000a1020400000, Magic: 0x0122042541300000, Shift: 54},
	{Mask: 0x0000142240000000, Magic: 0x4020402004101411, Shift: 54},
	{Mask: 0x0000284402000000, Magic: 0x0200204004098248, Shift: 54},
	{Mask: 0x0000500804020000, Magic: 0x0200b43200c08000, Shift: 54},
	{Mask: 0x0000201008040200, Magic: 0x1004a00462210100, Shift: 54},
	{Mask: 0x0000402010080400, Magic: 0x20110081414a0200, Shift: 53},
	{Mask: 0x0002040810204000, Magic: 0x0000040080288200, Shift: 52},
	{Mask: 0x0004081020400000, Magic: 0x1602020022c04040, Shift: 53},
	{Mask: 0x000a102040000000, Magic: 0x0008009230140020, Shift: 53},
	{Mask: 0x0014224000000000, Magic: 0x0003046482400188, Shift: 53},
	{Mask: 0x0028440200000000, Magic: 0x4411088009a00480, Shift: 53},
	{Mask: 0x0050080402000000, Magic: 0x00100001100008c4, Shift: 53},
	{Mask: 0x0020100804020000, Magic: 0x028a001210518880, Shift: 53},
	{Mask: 0x0040201008040200, Magic: 0x0219000c02005002, Shift: 52},
}

var RookMoves [64][]board.Bitboard
var BishopMoves [64][]board.Bitboard

type MagicEntry struct {
	Mask  board.Bitboard
	Magic uint64
	Shift uint8
}

func MagicIndex(entry *MagicEntry, occupied board.Bitboard) int {
	hash := uint64(occupied&entry.Mask) * entry.Magic
	return int(hash >> entry.Shift)
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
