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
	{Mask: 0x000101010101017e, Magic: 0xd080015220884000, Shift: 52},
	{Mask: 0x000202020202027c, Magic: 0x0440009000c02000, Shift: 53},
	{Mask: 0x000404040404047a, Magic: 0x0100182000524100, Shift: 53},
	{Mask: 0x0008080808080876, Magic: 0x010010000c090020, Shift: 53},
	{Mask: 0x001010101010106e, Magic: 0x0080080042800400, Shift: 53},
	{Mask: 0x002020202020205e, Magic: 0x0700020100040018, Shift: 53},
	{Mask: 0x004040404040403e, Magic: 0x0100240200008100, Shift: 53},
	{Mask: 0x008080808080807e, Magic: 0x848000c10000a080, Shift: 52},
	{Mask: 0x0001010101017e00, Magic: 0x4400800020804001, Shift: 53},
	{Mask: 0x0002020202027c00, Magic: 0x0002400040201002, Shift: 54},
	{Mask: 0x0004040404047a00, Magic: 0x0c87001041006001, Shift: 54},
	{Mask: 0x0008080808087600, Magic: 0x000200401a002110, Shift: 54},
	{Mask: 0x0010101010106e00, Magic: 0x0042000420d00a00, Shift: 54},
	{Mask: 0x0020202020205e00, Magic: 0x0001000400980300, Shift: 54},
	{Mask: 0x0040404040403e00, Magic: 0x0401002200044500, Shift: 54},
	{Mask: 0x0080808080807e00, Magic: 0x0089000300004082, Shift: 53},
	{Mask: 0x00010101017e0100, Magic: 0x0003010020c28000, Shift: 53},
	{Mask: 0x00020202027c0200, Magic: 0x1040030024804108, Shift: 54},
	{Mask: 0x00040404047a0400, Magic: 0x0000808010002000, Shift: 54},
	{Mask: 0x0008080808760800, Magic: 0x0004808010000800, Shift: 54},
	{Mask: 0x00101010106e1000, Magic: 0x1008808008000400, Shift: 54},
	{Mask: 0x00202020205e2000, Magic: 0x200a008004008280, Shift: 54},
	{Mask: 0x00404040403e4000, Magic: 0x00040c0012110810, Shift: 54},
	{Mask: 0x00808080807e8000, Magic: 0x14114a0000408104, Shift: 53},
	{Mask: 0x000101017e010100, Magic: 0x0000608180004004, Shift: 53},
	{Mask: 0x000202027c020200, Magic: 0x6002090200204280, Shift: 54},
	{Mask: 0x000404047a040400, Magic: 0x2005600100410012, Shift: 54},
	{Mask: 0x0008080876080800, Magic: 0x02000a0200204010, Shift: 54},
	{Mask: 0x001010106e101000, Magic: 0x0022001200140820, Shift: 54},
	{Mask: 0x002020205e202000, Magic: 0x4001040080060080, Shift: 54},
	{Mask: 0x004040403e404000, Magic: 0x0110010400481250, Shift: 54},
	{Mask: 0x008080807e808000, Magic: 0x1080800080024100, Shift: 53},
	{Mask: 0x0001017e01010100, Magic: 0x0082400880800022, Shift: 53},
	{Mask: 0x0002027c02020200, Magic: 0x0400201000404000, Shift: 54},
	{Mask: 0x0004047a04040400, Magic: 0x4101002003004030, Shift: 54},
	{Mask: 0x0008087608080800, Magic: 0x04100400c2400800, Shift: 54},
	{Mask: 0x0010106e10101000, Magic: 0x0202800800801c02, Shift: 54},
	{Mask: 0x0020205e20202000, Magic: 0x0842000c0a002811, Shift: 54},
	{Mask: 0x0040403e40404000, Magic: 0x0010020104000830, Shift: 54},
	{Mask: 0x0080807e80808000, Magic: 0x0000006902000084, Shift: 53},
	{Mask: 0x00017e0101010100, Magic: 0x2040800041010020, Shift: 53},
	{Mask: 0x00027c0202020200, Magic: 0x0000500a20024002, Shift: 54},
	{Mask: 0x00047a0404040400, Magic: 0x0013006000410018, Shift: 54},
	{Mask: 0x0008760808080800, Magic: 0x02002012000a0040, Shift: 54},
	{Mask: 0x00106e1010101000, Magic: 0x1100040008008080, Shift: 54},
	{Mask: 0x00205e2020202000, Magic: 0x01110c0002008080, Shift: 54},
	{Mask: 0x00403e4040404000, Magic: 0x0001080102040010, Shift: 54},
	{Mask: 0x00807e8080808000, Magic: 0x08040a40a9020004, Shift: 53},
	{Mask: 0x007e010101010100, Magic: 0x000280010040a100, Shift: 53},
	{Mask: 0x007c020202020200, Magic: 0x0400c00120100ac0, Shift: 54},
	{Mask: 0x007a040404040400, Magic: 0x000082a001100180, Shift: 54},
	{Mask: 0x0076080808080800, Magic: 0x0110a10210006900, Shift: 54},
	{Mask: 0x006e101010101000, Magic: 0x80001800800c0080, Shift: 54},
	{Mask: 0x005e202020202000, Magic: 0x1044802200040080, Shift: 54},
	{Mask: 0x003e404040404000, Magic: 0x0021011008960400, Shift: 54},
	{Mask: 0x007e808080808000, Magic: 0x0804802045000480, Shift: 53},
	{Mask: 0x7e01010101010100, Magic: 0x00010080a2001442, Shift: 52},
	{Mask: 0x7c02020202020200, Magic: 0x010100a081144202, Shift: 53},
	{Mask: 0x7a04040404040400, Magic: 0x1088403200808a62, Shift: 53},
	{Mask: 0x7608080808080800, Magic: 0x2001001000220865, Shift: 53},
	{Mask: 0x6e10101010101000, Magic: 0x0029003800043023, Shift: 53},
	{Mask: 0x5e20202020202000, Magic: 0x2002002804100102, Shift: 53},
	{Mask: 0x3e40404040404000, Magic: 0x080200080c150082, Shift: 53},
	{Mask: 0x7e80808080808000, Magic: 0xc600040464830442, Shift: 52},
}

var BishopMagics = [64]MagicEntry{
	{Mask: 0x0040201008040200, Magic: 0x0010550808008860, Shift: 58},
	{Mask: 0x0000402010080400, Magic: 0x0a020a2841010048, Shift: 59},
	{Mask: 0x0000004020100a00, Magic: 0x6068081060800000, Shift: 59},
	{Mask: 0x0000000040221400, Magic: 0x00209081008a00a8, Shift: 59},
	{Mask: 0x0000000002442800, Magic: 0x2c04042004004082, Shift: 59},
	{Mask: 0x0000000204085000, Magic: 0x0001901420001000, Shift: 59},
	{Mask: 0x0000020408102000, Magic: 0x0064020110281488, Shift: 59},
	{Mask: 0x0002040810204000, Magic: 0x0002002088041010, Shift: 58},
	{Mask: 0x0020100804020000, Magic: 0x8000881004089400, Shift: 59},
	{Mask: 0x0040201008040000, Magic: 0x2108080148022045, Shift: 59},
	{Mask: 0x00004020100a0000, Magic: 0x14004804c4008241, Shift: 59},
	{Mask: 0x0000004022140000, Magic: 0x0208080a00201428, Shift: 59},
	{Mask: 0x0000000244280000, Magic: 0x18181110411a4108, Shift: 59},
	{Mask: 0x0000020408500000, Magic: 0x40002c2208400000, Shift: 59},
	{Mask: 0x0002040810200000, Magic: 0x0080030802022004, Shift: 59},
	{Mask: 0x0004081020400000, Magic: 0x4184120046280482, Shift: 59},
	{Mask: 0x0010080402000200, Magic: 0x81400e4888019400, Shift: 59},
	{Mask: 0x0020100804000400, Magic: 0x0026042008061080, Shift: 59},
	{Mask: 0x004020100a000a00, Magic: 0x100200010a040900, Shift: 57},
	{Mask: 0x0000402214001400, Magic: 0xa40084480a024040, Shift: 57},
	{Mask: 0x0000024428002800, Magic: 0x0001010820080212, Shift: 57},
	{Mask: 0x0002040850005000, Magic: 0x1402010048020800, Shift: 57},
	{Mask: 0x0004081020002000, Magic: 0x80820000a2212000, Shift: 59},
	{Mask: 0x0008102040004000, Magic: 0x1100c002004a0850, Shift: 59},
	{Mask: 0x0008040200020400, Magic: 0x00840c04d1101013, Shift: 59},
	{Mask: 0x0010080400040800, Magic: 0x0009040420080602, Shift: 59},
	{Mask: 0x0020100a000a1000, Magic: 0x03024800b00024c0, Shift: 57},
	{Mask: 0x0040221400142200, Magic: 0x000a012002008200, Shift: 55},
	{Mask: 0x0002442800284400, Magic: 0x0081010008304000, Shift: 55},
	{Mask: 0x0004085000500800, Magic: 0x2430008089081b14, Shift: 57},
	{Mask: 0x0008102000201000, Magic: 0x08010404430c01a2, Shift: 59},
	{Mask: 0x0010204000402000, Magic: 0x6004010000804300, Shift: 59},
	{Mask: 0x0004020002040800, Magic: 0x031010041010a400, Shift: 59},
	{Mask: 0x0008040004081000, Magic: 0x2052080400021000, Shift: 59},
	{Mask: 0x00100a000a102000, Magic: 0x8000105000180088, Shift: 57},
	{Mask: 0x0022140014224000, Magic: 0x800a240108040100, Shift: 55},
	{Mask: 0x0044280028440200, Magic: 0x01020a0400020028, Shift: 55},
	{Mask: 0x0008500050080400, Magic: 0x4010040120111002, Shift: 57},
	{Mask: 0x0010200020100800, Magic: 0xa31101010002281b, Shift: 59},
	{Mask: 0x0020400040201000, Magic: 0x0001893201084202, Shift: 59},
	{Mask: 0x0002000204081000, Magic: 0x0201092020805000, Shift: 59},
	{Mask: 0x0004000408102000, Magic: 0x030104016c842004, Shift: 59},
	{Mask: 0x000a000a10204000, Magic: 0x8042001044110800, Shift: 57},
	{Mask: 0x0014001422400000, Magic: 0x00040d4204801800, Shift: 57},
	{Mask: 0x0028002844020000, Magic: 0x4000040102100400, Shift: 57},
	{Mask: 0x0050005008040200, Magic: 0x02466810030004e0, Shift: 57},
	{Mask: 0x0020002010080400, Magic: 0x10a0480b01001440, Shift: 59},
	{Mask: 0x0040004020100800, Magic: 0xa806081841004180, Shift: 59},
	{Mask: 0x0000020408102000, Magic: 0x0108840402401400, Shift: 59},
	{Mask: 0x0000040810204000, Magic: 0x1800444404200091, Shift: 59},
	{Mask: 0x00000a1020400000, Magic: 0x0404812402480010, Shift: 59},
	{Mask: 0x0000142240000000, Magic: 0x0800010284040100, Shift: 59},
	{Mask: 0x0000284402000000, Magic: 0x000210a00a048003, Shift: 59},
	{Mask: 0x0000500804020000, Magic: 0x2002088208020440, Shift: 59},
	{Mask: 0x0000201008040200, Magic: 0x0878101408004010, Shift: 59},
	{Mask: 0x0000402010080400, Magic: 0x0812080124009048, Shift: 59},
	{Mask: 0x0002040810204000, Magic: 0x0001040841041002, Shift: 58},
	{Mask: 0x0004081020400000, Magic: 0x0504020100821080, Shift: 59},
	{Mask: 0x000a102040000000, Magic: 0x0800010114010400, Shift: 59},
	{Mask: 0x0014224000000000, Magic: 0x200501080d208800, Shift: 59},
	{Mask: 0x0028440200000000, Magic: 0x4a00015020604100, Shift: 59},
	{Mask: 0x0050080402000000, Magic: 0x006801c2022c0101, Shift: 59},
	{Mask: 0x0020100804020000, Magic: 0x8301400801440080, Shift: 59},
	{Mask: 0x0040201008040200, Magic: 0x8808281004004019, Shift: 58},
}
