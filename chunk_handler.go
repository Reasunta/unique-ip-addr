package main

import "os"

type ChunkHandler struct {
	filename   string
	bufferSize int
}

func NewChunkHandler(filename string, bufferSize int) *ChunkHandler {
	return &ChunkHandler{filename: filename, bufferSize: bufferSize}
}

func (ch *ChunkHandler) Handle(counter *IPCounter) (uint64, uint32) {
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

	fileHandler := NewFileHandler(file, ch.bufferSize, counter)
	return fileHandler.countAddresses()
}
