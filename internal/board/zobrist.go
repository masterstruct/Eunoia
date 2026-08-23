package board

type PRNG struct {
	state uint64
}

// Splitmix64 pseudorandom number generator
func (p *PRNG) Next() uint64 {
	p.state += 0x9e3779b97f4a7c15
	z := p.state
	z = (z ^ (z >> 30)) * 0xbf58476d1ce4e5b9
	z = (z ^ (z >> 27)) * 0x94d049bb133111eb
	return z ^ (z >> 31)
}
