package storage

import "io"

type Storage interface {
	BlockByteSize() int                         // block byte size, I think 1<<20 = 1MB is a good value
	TotalSize() int                             // total file system size
	Read(blockIndex int) io.Reader              // Read a block
	Write(blockIndex int, data io.Reader) error // Write a block
}
