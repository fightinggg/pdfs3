package pdfs

import (
	"errors"
	"golang.org/x/net/webdav"
	"io/fs"
)

type CachePdfsFile struct {
	proxy     webdav.File
	cache     []byte
	cacheSize int
}

func NewCachePdfsFile(proxy webdav.File, cacheWriteCap int) *CachePdfsFile {
	return &CachePdfsFile{
		proxy:     proxy,
		cache:     make([]byte, cacheWriteCap),
		cacheSize: 0,
	}
}

func (f *CachePdfsFile) Flush() error {
	if f.cacheSize != 0 {
		write, err := f.proxy.Write(f.cache[0:f.cacheSize])
		if err != nil {
			return err
		}
		if f.cacheSize != write {
			return errors.New("f.cacheSize != write")
		}
		f.cacheSize = 0
	}
	return nil
}

func (f *CachePdfsFile) Close() error {
	err := f.Flush()
	if err != nil {
		return err
	}
	return f.proxy.Close()
}

func (f *CachePdfsFile) Read(p []byte) (n int, err error) {
	return f.proxy.Read(p)
}

func (f *CachePdfsFile) Seek(offset int64, whence int) (int64, error) {
	err := f.Flush()
	if err != nil {
		return 0, err
	}
	return f.proxy.Seek(offset, whence)
}

func (f *CachePdfsFile) Readdir(count int) (res []fs.FileInfo, err error) {
	return f.proxy.Readdir(count)
}

func (f *CachePdfsFile) Stat() (fs.FileInfo, error) {
	return f.proxy.Stat()
}

func (f *CachePdfsFile) Write(p []byte) (int, error) {
	write := len(f.cache) - f.cacheSize
	if len(p) < write {
		write = len(p)
	}
	for i := 0; i < write; i++ {
		f.cache[i+f.cacheSize] = p[i]
	}
	f.cacheSize += write

	if f.cacheSize == len(f.cache) {
		err := f.Flush()
		return write, err
	} else {
		return write, nil
	}
}
