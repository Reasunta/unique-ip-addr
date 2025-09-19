package main

import (
	"bytes"
	"errors"
	"fmt"
	"math"
	"os"
	"runtime"
	"sync"
	"time"
)

const (
	BufferSize               = 1 << 20 //Size of file reader buffer in bytes
	ChannelBufferSize        = 1 << 8  //Size of channels buffer
	Limit                    = 0       // Debug feature, allows to stop every working thread after handling at least Limit addresses
	SendLimit         uint64 = 1000000 // Working thread sends number of handled ips to main thread when at least SendLimit addresses are handled
)

func calcChunks(filename string, parts int) ([]int64, error) {
	file, err := os.Open(filename)
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			panic(err)
		}
	}(file)

	if err != nil {
		return nil, err
	}
	if parts < 1 {
		return nil, errors.New("parts must be greater than 0")
	}

	buffer := make([]byte, 20)
	result := make([]int64, 2*parts)

	stat, _ := file.Stat()
	approxLength := stat.Size() / int64(parts)

	result[0] = 0
	for i := 0; i < parts-1; i++ {
		lastIndex := approxLength * int64(i+1)
		_, err = file.ReadAt(buffer, lastIndex)
		if err != nil {
			return nil, err
		}
		lastIndex += int64(bytes.IndexByte(buffer, '\n'))

		result[2*i+1] = lastIndex - result[2*i] + 1 //size = end - start + 1
		result[2*i+2] = lastIndex + 1               // start of next chunk
	}
	result[2*parts-1] = stat.Size() - result[2*parts-2]
	return result, nil
}

func main() {
	start := time.Now()
	filename := os.Args[1]

	var workers = runtime.NumCPU() - 2
	if workers < 1 {
		workers = 1
	}
	fmt.Printf("Working thread count: %d\n", workers)
	var m1, m2 runtime.MemStats
	runtime.GC()
	runtime.ReadMemStats(&m1)

	var handled uint64 = 0
	var unique uint32 = 0

	var h = make(chan uint64, ChannelBufferSize)
	var u = make(chan uint32, ChannelBufferSize)
	counter := NewIPCounter()

	chunks, err := calcChunks(filename, workers)

	if err != nil {
		panic(err)
	}

	var finishedChannels = 0

	wg := &sync.WaitGroup{}
	wg.Add(workers + 1)
	go func() {
		defer wg.Done()
		for {
			if finishedChannels == 2*workers {
				break
			}
			select {
			case hd := <-h:
				if hd == math.MaxUint64 {
					finishedChannels++
				} else {
					handled += hd
					fmt.Printf("Handled address count: %d\r", handled)
				}
			case ud := <-u:
				if ud == math.MaxUint32 {
					finishedChannels++
				} else {
					unique += ud
				}
			}
		}
	}()

	for i := 0; i < workers; i++ {
		ch := NewChunkHandler(filename, BufferSize, chunks[2*i], chunks[2*i+1])
		go ch.Handle(counter, h, u, Limit, wg)
	}
	wg.Wait()

	runtime.ReadMemStats(&m2)
	elapsed := time.Since(start)

	fmt.Printf("Handled address count: %d", handled)
	fmt.Printf("\nUnique address count: %d", unique)

	velocity := float64(handled) / float64(elapsed.Milliseconds())
	fmt.Printf("\nAverage velocity: %f ip / ms", velocity)
	fmt.Printf("\nSpent time: %s", elapsed)
	fmt.Printf("\nMemory usage: %fMB", float64(m2.HeapAlloc-m1.HeapAlloc)/1024.0/1024.0)
	fmt.Println("\nMallocs:", m2.Mallocs-m1.Mallocs)
}
