// Taken from Rodent V - thank you!
// https://github.com/nescitus/Rodent-V/blob/main/prefetch_amd64.s

#include "textflag.h" // for NOSPLIT

TEXT ·ttPrefetch(SB), NOSPLIT, $0-8
	MOVQ address+0(FP), AX  // loads TT entry address into AX
	PREFETCHT0 (AX)         // tries to cache this entry
	RET
