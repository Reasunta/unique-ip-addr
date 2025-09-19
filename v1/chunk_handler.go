package main

import (
	"os"
	"sync"
)

type ChunkHandler struct {
	filename   string
	bufferSize int
	from       int64
	size       int64
}

func NewChunkHandler(filename string, bufferSize int, from int64, size int64) *ChunkHandler {
	return &ChunkHandler{filename: filename, bufferSize: bufferSize, from: from, size: size}
}

func (ch *ChunkHandler) Handle(counter *IPCounter, h chan uint64, u chan uint32, limit uint64, wg *sync.WaitGroup) {
	defer wg.Done()

	file, err := os.Open(ch.filename)
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			panic(err)
		}
	}(file)

	if err != nil {
		panic(err)
	}
	_, _ = file.Seek(ch.from, 0)

	fileHandler := NewFileHandler(file, ch.bufferSize, counter)
	fileHandler.countAddresses(h, u, limit, ch.size)
}
