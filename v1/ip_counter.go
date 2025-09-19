package main

import "sync/atomic"

// IPCounter represents a huge number with a length of 2^32 bits. It is equal to max amount of unique ip addresses.
type IPCounter struct {
	counters [1 << 26]uint64
}

func NewIPCounter() *IPCounter {
	return &IPCounter{}
}

func (i *IPCounter) handle(ip uint32) uint32 {
	index := ip >> 6
	bit := ip % 64
	mask := uint64(1 << bit)

	old := atomic.OrUint64(&i.counters[index], mask)
	if old&mask == 0 {
		return 1
	}

	return 0
}
