package storage

import (
	"bytes"
	"errors"
	"io"
	"log"
)

type memoryStorage struct {
	memory    [][]byte
	blockSize int
	blockNums int
}

func (m memoryStorage) BlockByteSize() int {
	return m.blockSize
}

func (m memoryStorage) TotalSize() int {
	return m.blockSize * m.blockNums
}

func NewMemoryStorage(blockSize, blockNums int) Storage {
	storage := memoryStorage{}

	storage.blockSize = blockSize
	storage.blockNums = blockNums
	for i := 0; i < storage.blockSize; i++ {
		storage.memory = append(storage.memory, make([]byte, blockNums))
	}

	return storage
}

func (m memoryStorage) Read(blockIndex int) io.Reader {
	log.Printf("INFO : read block %d", blockIndex)
	return bytes.NewReader(m.memory[blockIndex])
}

func (m memoryStorage) Write(blockIndex int, data io.Reader) error {
	log.Printf("INFO : write block %d", blockIndex)
	all, err := io.ReadAll(data)
	if err != nil {
		return err
	}
	if len(all) > m.blockSize {
		return errors.New("size err")
	}
	m.memory[blockIndex] = all
	return nil
}
