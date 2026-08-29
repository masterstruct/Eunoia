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
	mask = bishopAttacksSlow(sq, board.EmptyBB)
	return mask &^ board.EdgesBB
}

type MagicEntry struct {
	Mask   board.Bitboard
	Magic  board.Bitboard
	Shift  uint8
	Offset int
}

func MagicIndex(entry *MagicEntry, occupied board.Bitboard) int {
	return int(((occupied&entry.Mask)*entry.Magic)>>entry.Shift) + entry.Offset
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
		rookEntry.Shift = uint8(64 - rookEntry.Mask.CountBits())
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
		bishopEntry.Shift = uint8(64 - bishopEntry.Mask.CountBits())
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
			moves = rookAttacksSlow(sq, blockers)
		} else {
			moves = bishopAttacksSlow(sq, blockers)
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

var RookMoves [RookTableSize]board.Bitboard
var BishopMoves [BishopTableSize]board.Bitboard

const RookTableSize = 102400
const BishopTableSize = 5248

var RookMagics = [64]MagicEntry{
	{Magic: 0x1180001a40008460}, {Magic: 0x0100104002810021},
	{Magic: 0x6b00082001051040}, {Magic: 0x8880041000800800},
	{Magic: 0x0200100408204200}, {Magic: 0x1200045801020010},
	{Magic: 0x0080088001000200}, {Magic: 0x82001040850201a4},
	{Magic: 0x0140800080314000}, {Magic: 0x1021004000208b04},
	{Magic: 0x0001801000200081}, {Magic: 0x0101002100081001},
	{Magic: 0x0612001008042200}, {Magic: 0x0100800200802c00},
	{Magic: 0x108100010052006c}, {Magic: 0x0000801080004100},
	{Magic: 0x0000208000804000}, {Magic: 0x0021010020904001},
	{Magic: 0x0010008014200080}, {Magic: 0xa068818008041000},
	{Magic: 0x4030110008000501}, {Magic: 0x2001010008420c00},
	{Magic: 0x00080c0061081210}, {Magic: 0x980082000100844c},
	{Magic: 0x0200400480002080}, {Magic: 0xc801200980400180},
	{Magic: 0x0061811200420020}, {Magic: 0x0010080080100080},
	{Magic: 0x8000050100080050}, {Magic: 0x0004008080060004},
	{Magic: 0x1016110400300802}, {Magic: 0x0801001100026082},
	{Magic: 0x4000400421800081}, {Magic: 0x002080c000802000},
	{Magic: 0x0800104901002000}, {Magic: 0x0010000821001100},
	{Magic: 0x1800080180800400}, {Magic: 0x0202002802001004},
	{Magic: 0x0000081004009201}, {Magic: 0x0401000041000086},
	{Magic: 0x2408800105470020}, {Magic: 0x2410012002404000},
	{Magic: 0x001004080020a000}, {Magic: 0x0102001058420020},
	{Magic: 0x0805080004008080}, {Magic: 0x2002000430020008},
	{Magic: 0x0480040200010100}, {Magic: 0x0080a04484020011},
	{Magic: 0x0220400020800480}, {Magic: 0x2410812102004200},
	{Magic: 0x8008410018200100}, {Magic: 0x00001001b8008080},
	{Magic: 0x4008000904008080}, {Magic: 0x4232800400420080},
	{Magic: 0x00010004ce000100}, {Magic: 0x410600830c004600},
	{Magic: 0x9080118001004021}, {Magic: 0x0898104201012086},
	{Magic: 0x420101406000a811}, {Magic: 0x8024082420100101},
	{Magic: 0x4111002210080085}, {Magic: 0x0006001001442806},
	{Magic: 0x0403004084020011}, {Magic: 0x4041016100804406},
}

var BishopMagics = [64]MagicEntry{
	{Magic: 0x4016041002004100}, {Magic: 0x0002240106020024},
	{Magic: 0x0204040082108400}, {Magic: 0x2124051200000204},
	{Magic: 0x010c046000080000}, {Magic: 0x6209104210408002},
	{Magic: 0x0042009044100040}, {Magic: 0x9000570808040400},
	{Magic: 0x2401488604040400}, {Magic: 0x40084830194200a0},
	{Magic: 0x0000220204002a11}, {Magic: 0x0060044100200000},
	{Magic: 0x0022462110000000}, {Magic: 0x6001010460140020},
	{Magic: 0x0010448084a02004}, {Magic: 0x080200a2010c2080},
	{Magic: 0x0040019154070420}, {Magic: 0x0182602424040400},
	{Magic: 0x0006020404015200}, {Magic: 0x0090240804204201},
	{Magic: 0x840420e202010280}, {Magic: 0x800100120100a600},
	{Magic: 0x0001c60201042080}, {Magic: 0x8000800024014802},
	{Magic: 0x002a501040040800}, {Magic: 0x08440200105a2810},
	{Magic: 0x0096208210018284}, {Magic: 0x0002080024004218},
	{Magic: 0x6183001001004001}, {Magic: 0x0200a20000880400},
	{Magic: 0x0203820815051000}, {Magic: 0x0c82004a0a0101a8},
	{Magic: 0x2010102430080849}, {Magic: 0x1214b00900040800},
	{Magic: 0x01a1141005020080}, {Magic: 0x9062040400080210},
	{Magic: 0x4005010c00020021}, {Magic: 0x0008100442008800},
	{Magic: 0x0e48008081111800}, {Magic: 0xc003120082120040},
	{Magic: 0x0804241240000828}, {Magic: 0x80008a8410a06001},
	{Magic: 0x14000a0090000600}, {Magic: 0x00d842842020a400},
	{Magic: 0x0101904200821810}, {Magic: 0x0820480108080040},
	{Magic: 0x002224140c006090}, {Magic: 0x00480a0c00400420},
	{Magic: 0x040a082d04104000}, {Magic: 0xa007024210441008},
	{Magic: 0x000002240608000d}, {Magic: 0x0300801284040600},
	{Magic: 0x18a8000c25040000}, {Magic: 0x0838502001090400},
	{Magic: 0x0820219a02004000}, {Magic: 0x0403100102008924},
	{Magic: 0x0004808401200600}, {Magic: 0x220210b209100800},
	{Magic: 0x00200ca02a011000}, {Magic: 0x1000008102420880},
	{Magic: 0x0000880840104100}, {Magic: 0x1102a00604280200},
	{Magic: 0x0800200404008400}, {Magic: 0x0048810844840080},
}
