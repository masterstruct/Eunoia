// Taken from Rodent V - thank you!
// https://github.com/nescitus/Rodent-V/blob/main/prefetch_amd64.go

package tt

import "unsafe"

func ttPrefetch(address unsafe.Pointer)

func (tt *Table) Prefetch(key uint64) {
	ttPrefetch(unsafe.Pointer(&tt.entries[key&tt.mask]))
}
