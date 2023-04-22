package pdfs

import (
	"context"
	"golang.org/x/net/webdav"
	"os"
	"pdfs3/pdfs/pdfs_lower_api"
)

type Pdfs struct {
	LowerApi pdfs_lower_api.LowerApi
}

func (pdfs *Pdfs) Mkdir(ctx context.Context, name string, perm os.FileMode) error {
	return pdfs.LowerApi.Mkdir(name)
}

func (pdfs *Pdfs) OpenFile(ctx context.Context, name string, flag int, perm os.FileMode) (webdav.File, error) {
	return &PdfsFile{
		pdfs:     pdfs,
		filename: name,
		offset:   0,
		reader:   nil,
		flag:     flag,
	}, nil
}

func (pdfs *Pdfs) RemoveAll(ctx context.Context, name string) error {
	return pdfs.LowerApi.Delete(name)
}

func (pdfs *Pdfs) Rename(ctx context.Context, oldName, newName string) error {
	return pdfs.LowerApi.Rename(oldName, newName)
}

func (pdfs *Pdfs) Stat(ctx context.Context, name string) (os.FileInfo, error) {
	stat, err := pdfs.LowerApi.Stat(name)
	if err != nil {
		return nil, err
	}
	return stat, err
}
