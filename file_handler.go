package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"slices"
)

type FileHandler struct {
	reader     *bufio.Reader
	counter    *IPCounter
	bufferSize int
	tail       []byte
}

func NewFileHandler(file *os.File, bufferSize int, counter *IPCounter) *FileHandler {
	return &FileHandler{
		reader:     bufio.NewReaderSize(file, bufferSize),
		counter:    counter,
		bufferSize: bufferSize,
		tail:       make([]byte, 0, 20),
	}
}

func (fh *FileHandler) handleIP(line []byte) uint32 {
	var result uint32 = 0
	var powIndex, itemIndex uint32 = 1, 1

	for i := len(line) - 1; i >= 0; i-- {
		if line[i] == 10 || line[i] == 13 {
			continue
		}
		if line[i] == 46 {
			powIndex *= 256
			itemIndex = 1
			continue
		}
		result += uint32(line[i]-48) * itemIndex * powIndex
		itemIndex *= 10
	}

	if !fh.counter.has(result) {
		fh.counter.add(result)
		return 1
	}
	return 0
}

func (fh *FileHandler) handleBuffer(buffer []byte, size int) (uint64, uint32) {
	startIndex := 0
	var handled uint64 = 0
	var unique uint32 = 0
	var useTail = true

	for i := 0; i < size; i++ {
		if buffer[i] == '\n' {
			if useTail {
				unique += fh.handleIP(append(fh.tail, buffer[:i]...))
				useTail = false
			} else {
				unique += fh.handleIP(buffer[startIndex:i])
			}
			handled++
			startIndex = i
		}
	}
	fh.tail = slices.Clone(buffer[startIndex:size])
	return handled, unique
}

func (fh *FileHandler) countAddresses() (uint64, uint32) {
	buffer := make([]byte, fh.bufferSize)
	var handled uint64 = 0
	var unique uint32 = 0

	for n, err := fh.reader.Read(buffer); n > 0; n, err = fh.reader.Read(buffer) {
		if err != nil && err != io.EOF {
			fmt.Println(err)
			break
		}
		hd, td := fh.handleBuffer(buffer, n)
		handled += hd
		unique += td

		fmt.Printf("Handled address count: %d\r", handled)
	}

	return handled, unique
}
