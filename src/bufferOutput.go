package bufferOutput

import "os"

type BufferOutput struct {
	data []byte
	pos int
}

func (buffer *BufferOutput) Create(size int) {
	buffer.data = make([]byte, size)
	buffer.pos = 0
}

func (buffer *BufferOutput) Reset() {
	buffer.pos = 0
}

func (buffer *BufferOutput) Write(file *os.File) (n int, err error) {
	n, err = file.Write(buffer.data[:buffer.pos])
	buffer.pos = 0
	return
}

func (buffer *BufferOutput) AddUint32(value uint32) {
	pos := buffer.pos
	memory := buffer.data
	memory[pos] = byte((value >>  0) & 0xFF)
	memory[pos+1] = byte((value >>  8) & 0xFF)
	memory[pos+2] = byte((value >> 16) & 0xFF)
	memory[pos+3] = byte((value >> 24) & 0xFF)
	buffer.pos += 4
}

func (buffer *BufferOutput) AddUint16(value uint16) {
	pos := buffer.pos
	memory := buffer.data
	memory[pos] = byte((value >> 0) & 0xFF)
	memory[pos+1] = byte((value >> 8) & 0xFF)
	buffer.pos += 2
}

