package pdfs

import (
	"bytes"
	"context"
	"errors"
	"io"
	"io/fs"
	"log"
	"os"
)

type PdfsFile struct {
	pdfs     *Pdfs
	filename string
	offset   int64
	reader   io.Reader
	flag     int
}

func (pdfsFile *PdfsFile) Close() error {
	//TODO implement me
	log.Printf("implement me PdfsFile.Close()")
	return nil
}

func (pdfsFile *PdfsFile) Read(p []byte) (n int, err error) {
	if pdfsFile.reader == nil {
		_, err := pdfsFile.Seek(0, io.SeekStart)
		if err != nil {
			return 0, err
		}
	}
	return pdfsFile.reader.Read(p)
}

func (pdfsFile *PdfsFile) Seek(offset int64, whence int) (int64, error) {
	if whence == io.SeekEnd { // reading...
		file, err := pdfsFile.pdfs.LowerApi.Read(pdfsFile.filename)
		if err != nil {
			return 0, err
		}

		all, err := io.ReadAll(file)
		if err != nil {
			return 0, err
		}

		pdfsFile.offset = int64(len(all)) + offset
		pdfsFile.reader = bytes.NewReader([]byte{})
	} else if whence == io.SeekStart { // reading...
		pdfsFile.offset = offset
		var err error
		pdfsFile.reader, err = pdfsFile.pdfs.LowerApi.Read(pdfsFile.filename)
		if err != nil {
			return 0, err
		}
	} else {
		return 0, errors.New("ER")
	}
	return pdfsFile.offset, nil

}

func (pdfsFile *PdfsFile) Readdir(count int) (res []fs.FileInfo, err error) {
	ls, err := pdfsFile.pdfs.LowerApi.Ls(pdfsFile.filename)
	if err != nil {
		return nil, err
	}

	for i := range ls {
		res = append(res, &ls[i])
	}

	return
}

func (pdfsFile *PdfsFile) Stat() (fs.FileInfo, error) {
	stat, err := pdfsFile.pdfs.Stat(context.Background(), pdfsFile.filename)
	if pdfsFile.flag&os.O_CREATE != 0 && err == os.ErrNotExist {
		err := pdfsFile.pdfs.LowerApi.Write(pdfsFile.filename, 0, 0, bytes.NewReader([]byte{}))
		if err != nil {
			return nil, err
		}
		return pdfsFile.pdfs.Stat(context.Background(), pdfsFile.filename)
	}
	return stat, err
}

func (pdfsFile *PdfsFile) Write(p []byte) (int, error) {
	err := pdfsFile.pdfs.LowerApi.Write(pdfsFile.filename, pdfsFile.offset, int64(len(p)), bytes.NewReader(p))
	pdfsFile.offset += int64(len(p))
	if err != nil {
		return 0, err
	}
	log.Printf("write %d", len(p))
	return len(p), nil
}
