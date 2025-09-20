package main

import (
	"bufio"
	"fmt"
	"io"
	"math"
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

func (fh *FileHandler) handleIP(line []byte) (uint32, int) {
	var result uint32 = 0
	var octet, shift uint32 = 0, 24

	for i := 0; i < len(line); i++ {
		if line[i] == 13 {
			continue
		}
		if line[i] == 10 {
			result |= octet
			return fh.counter.handle(result), i
		}
		if line[i] == 46 {
			if octet > 255 {
				panic(fmt.Sprintf("Too big octet: %d in string %s", octet, string(line)))
			}
			result |= octet << shift
			octet = 0
			shift -= 8
			continue
		}
		octet = octet*10 + uint32(line[i]-'0')
	}
	return 0, -1
}

func (fh *FileHandler) handleBuffer(buffer []byte, size int) (uint64, uint32) {
	startIndex := 0
	var handled uint64 = 0
	var unique uint32 = 0

	u, endIndex := fh.handleIP(append(fh.tail, buffer[:MaxIpSize]...))
	unique += u
	handled++
	startIndex = endIndex - len(fh.tail) + 1

	for i := startIndex; i < size; i++ {
		ipSize := MaxIpSize
		if size-i < ipSize {
			ipSize = size - i
		}

		u, endIndex := fh.handleIP(buffer[i : i+ipSize])
		if endIndex > 0 {
			unique += u
			handled++

			i += endIndex
			startIndex = i + 1
		} else {
			startIndex = i
			break
		}
	}
	fh.tail = slices.Clone(buffer[startIndex:size])
	return handled, unique
}

func (fh *FileHandler) countAddresses(h chan uint64, u chan uint32, limit uint64, size int64) {
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
	h <- math.MaxUint64
	u <- math.MaxUint32
}
