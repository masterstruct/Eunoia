// Taken from Rodent V - thank you!
// https://github.com/nescitus/Rodent-V/blob/main/prefetch_other.go

//go:build !amd64

package tt

func (tt *Table) Prefetch(key uint64) {}
