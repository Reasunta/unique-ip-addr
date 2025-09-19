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

	return fh.counter.handle(result)
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

func (fh *FileHandler) countAddresses(h chan uint64, u chan uint32, done chan bool, limit uint64, size int64) {
	buffer := make([]byte, fh.bufferSize)
	var handled, sum uint64 = 0, 0
	var unique uint32 = 0
	var readSize int64 = 0

	for n, err := fh.reader.Read(buffer); n > 0 && readSize < size; n, err = fh.reader.Read(buffer) {
		if err != nil && err != io.EOF {
			fmt.Println(err)
			break
		}
		if readSize+int64(n) > size {
			n = int(size - readSize)
		}
		hd, td := fh.handleBuffer(buffer, n)
		readSize += int64(n)

		handled += hd
		sum += hd
		unique += td

		if handled > SendLimit {
			h <- handled
			handled = 0
		}

		if limit > 0 && sum > limit {
			break
		}
	}
	h <- handled
	u <- unique
	done <- true
}
