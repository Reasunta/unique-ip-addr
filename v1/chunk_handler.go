package main

import "os"

type ChunkHandler struct {
	filename   string
	bufferSize int
	from       int64
	size       int64
}

func NewChunkHandler(filename string, bufferSize int, from int64, size int64) *ChunkHandler {
	return &ChunkHandler{filename: filename, bufferSize: bufferSize, from: from, size: size}
}

func (ch *ChunkHandler) Handle(counter *IPCounter, h chan uint64, u chan uint32, done chan bool, limit uint64) {
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
	fileHandler.countAddresses(h, u, done, limit, ch.size)
}
