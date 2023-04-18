package pdfs

import (
	"bytes"
	"errors"
	"io"
	"io/fs"
	"log"
	"os"
	"pdfs3/pdfs/pdfs_err"
	"strings"
)

type PdfsFile struct {
	simplePdfs *SimplePdfs
	filename   string
	offset     int64
	reader     io.Reader
	flag       int
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
		file, err := pdfsFile.simplePdfs.ReadFile(pdfsFile.filename, 0, 0)
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
		pdfsFile.reader, err = pdfsFile.simplePdfs.ReadFile(pdfsFile.filename, 0, 0)
		if err != nil {
			return 0, err
		}
	} else {
		return 0, errors.New("ER")
	}
	return pdfsFile.offset, nil

}

func (pdfsFile *PdfsFile) Readdir(count int) ([]fs.FileInfo, error) {
	dir := pdfsFile.simplePdfs.listDir(pdfsFile.filename)
	var res []fs.FileInfo

	for _, info := range dir {
		fileInfo := PdfsFileInfo{
			info: info,
		}
		res = append(res, &fileInfo)
	}
	return res, nil
}

func (pdfsFile *PdfsFile) Stat() (fs.FileInfo, error) {
	if strings.HasSuffix(pdfsFile.filename, "/") {
		return &PdfsFileInfo{
			info: SimplePdfsFileInfo{
				Filename: pdfsFile.filename,
				IsDir:    true,
			},
		}, nil
	} else {
		// not lazy dir
		for _, info := range pdfsFile.simplePdfs.metaCache.AllFileList {
			if info.Filename == pdfsFile.filename || info.Filename == pdfsFile.filename+"/" {
				if !info.IsDir {
					file, err := pdfsFile.simplePdfs.ReadFile(pdfsFile.filename, 0, 0)
					if err != nil {
						return nil, err
					}

					all, err := io.ReadAll(file)
					if err != nil {
						return nil, err
					}

					return &PdfsFileInfo{
						info: SimplePdfsFileInfo{
							Filename: pdfsFile.filename,
							IsDir:    info.IsDir,
							Size:     len(all),
						},
					}, nil
				}
				return &PdfsFileInfo{
					info: SimplePdfsFileInfo{
						Filename: pdfsFile.filename,
						IsDir:    info.IsDir,
					},
				}, nil
			}
		}

		if (pdfsFile.flag & os.O_CREATE) != 0 {
			_, err := pdfsFile.Write([]byte{})
			if err != nil {
				return nil, err
			}
			return pdfsFile.Stat()
		} else {
			return nil, pdfs_err.PdfsNotFoundError()
		}

	}

}

func (pdfsFile *PdfsFile) Write(p []byte) (int, error) {
	err := pdfsFile.simplePdfs.WriteFile(pdfsFile.filename, 0, 0, bytes.NewReader(p))
	if err != nil {
		return 0, err
	}
	return len(p), nil
}
