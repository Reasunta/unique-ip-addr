package main

// IPCounter represents a huge number with a length of 2^32 bits. It is equal to max amount of unique ip addresses.
type IPCounter struct {
	counters [1 << 26]uint64
}

func NewIPCounter() *IPCounter {
	return &IPCounter{}
}

// has checks a bit with number equal integer representation of given ip
func (i *IPCounter) has(ip uint32) bool {
	index := ip >> 6
	bit := ip % 64

	item := i.counters[index]
	return item&(1<<bit) != 0
}

// add sets a bit with number equal integer representation of given ip
func (i *IPCounter) add(ip uint32) {
	index := ip >> 6
	bit := ip % 64

	i.counters[index] += 1 << bit
}
