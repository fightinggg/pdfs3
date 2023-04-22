package pdfs_lower_api

import (
	"pdfs3/pdfs/storage"
	"time"
)

type PdfsLowerAPiLS struct {
	FileName   string    `json:"fileName,omitempty"` // < 255B
	FileSize   int64     `json:"fileSize,omitempty"` // = 8B
	FileUpdate time.Time `json:"fileUpdate"`         // =10B
	IsDir      bool      `json:"isDir,omitempty"`    // =1B
	FileData   []byte    `json:"fileData"`           // <= 1KB
}

type PdfsLowerApi struct {
	storage storage.Storage
}

func (api *PdfsLowerApi) ls(path string) ([]PdfsLowerAPiLS, error) {

	return nil, nil
}

func (api *PdfsLowerApi) mkdir(path string) ([]PdfsLowerAPiLS, error) {
	return nil, nil

}

func (api *PdfsLowerApi) write(path string) ([]PdfsLowerAPiLS, error) {

	return nil, nil
}

func (api *PdfsLowerApi) read(path string) ([]PdfsLowerAPiLS, error) {
	return nil, nil

}

func (api *PdfsLowerApi) delete(path string) ([]PdfsLowerAPiLS, error) {
	return nil, nil

}
