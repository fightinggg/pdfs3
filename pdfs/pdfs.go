package pdfs

import (
	"context"
	"golang.org/x/net/webdav"
	"os"
	"pdfs3/pdfs/storage"
)

type Pdfs struct {
	simplePdfs SimplePdfs
}

func NewPdfsFromStorage(storage storage.Storage) (webdav.FileSystem, error) {
	fromStorage, err := NewSimplePdfsFromStorage(storage)
	return &Pdfs{
		simplePdfs: fromStorage,
	}, err
}

func NewPdfsFromStorageAndIgnoreAllError(storage storage.Storage) webdav.FileSystem {
	fromStorage := NewSimplePdfsFromStorageAndIgnoreAllError(storage)
	return &Pdfs{
		simplePdfs: fromStorage,
	}
}

func (pdfs *Pdfs) Mkdir(ctx context.Context, name string, perm os.FileMode) error {
	pdfs.simplePdfs.mkdir(name)
	return nil
}

func (pdfs *Pdfs) OpenFile(ctx context.Context, name string, flag int, perm os.FileMode) (webdav.File, error) {
	return &PdfsFile{
		simplePdfs: &pdfs.simplePdfs,
		filename:   name,
		flag:       flag,
	}, nil

}

func (pdfs *Pdfs) RemoveAll(ctx context.Context, name string) error {
	//TODO implement me
	panic("implement me")
}

func (pdfs *Pdfs) Rename(ctx context.Context, oldName, newName string) error {
	list := pdfs.simplePdfs.metaCache.AllFileList
	for i := range list {
		if list[i].Filename == oldName {
			list[i].Filename = newName
		}
	}
	return nil
	// TODO to disk
}

func (pdfs *Pdfs) Stat(ctx context.Context, name string) (os.FileInfo, error) {
	file, err := pdfs.OpenFile(ctx, name, 0, 0)
	if err != nil {
		return nil, err
	}
	stat, err := file.Stat()
	if err != nil {
		return nil, err
	}
	err = file.Close()
	if err != nil {
		return nil, err
	}

	return stat, err

}
