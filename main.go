package main

import (
	"fmt"
	"os"
	"runtime"
	"time"
)

const (
	BufferSize = 4096 * 32
)

func main() {
	start := time.Now()
	filename := os.Args[1]

	var m1, m2 runtime.MemStats
	runtime.GC()
	runtime.ReadMemStats(&m1)

	counter := NewIPCounter()
	var handled uint64 = 0
	var total uint32 = 0

	ch := NewChunkHandler(filename, BufferSize)
	handled, total = ch.Handle(counter)

	runtime.ReadMemStats(&m2)
	elapsed := time.Since(start)

	fmt.Printf("Handled address count: %d", handled)
	fmt.Printf("\nUnique address count: %d", total)

	velocity := float64(handled) / float64(elapsed.Milliseconds())
	fmt.Printf("\nAverage velocity: %f ip / ms", velocity)
	fmt.Printf("\nSpent time: %s", elapsed)
	fmt.Printf("\nMemory usage: %fMB", float64(m2.HeapAlloc-m1.HeapAlloc)/1024.0/1024.0)
	fmt.Println("\nMallocs:", m2.Mallocs-m1.Mallocs)
}
