package pdfs_lower_api

import (
	"encoding/json"
	"io"
	"pdfs3/pdfs/storage"
	"strings"
	"time"
)

type MenuBlock struct {
	PathPrefix    string           `json:"pathPrefix"`
	NextMenuIndex int              `json:"nextMenuIndex"`
	Menus         []PdfsLowerAPiLS `json:"menus"`
}

type PdfsLowerAPiLS struct {
	FileName      string    `json:"fileName,omitempty"` // < 255B
	FileSize      int64     `json:"fileSize,omitempty"` // = 8B
	FileUpdate    time.Time `json:"fileUpdate"`         // =10B
	IsDir         bool      `json:"isDir,omitempty"`    // =1B
	DirBlockIndex int       `json:"dirBlockIndex"`      // =2B
	FileData      []byte    `json:"fileData"`           // <= 1KB
}

type PdfsLowerApi struct {
	storage storage.Storage
}

func (api *PdfsLowerApi) readMenus(blockIndex int) (*MenuBlock, error) {
	read := api.storage.Read(blockIndex)
	all, err := io.ReadAll(read)
	if err != nil {
		return nil, err
	}

	var menus MenuBlock
	err = json.Unmarshal(all, &menus)
	if err != nil {
		return nil, err
	}

	return &menus, nil
}

func (api *PdfsLowerApi) Ls(path string) ([]PdfsLowerAPiLS, error) {
	if !strings.HasSuffix(path, "/") {
		path = path + "/"
	}
	var res []PdfsLowerAPiLS
	var nextIndex = 0
	for true {
		menus, err := api.readMenus(nextIndex)
		if err != nil {
			return nil, err
		}

		nextIndex = menus.NextMenuIndex

		for _, menu := range menus.Menus {
			menu.FileName = menus.PathPrefix + menu.FileName
			if menu.FileName == path && menu.DirBlockIndex > 0 {
				nextIndex = menu.DirBlockIndex
				break
			} else if strings.HasPrefix(menu.FileName, path) {
				res = append(res, menu)
			}
		}

		if nextIndex <= 0 {
			// TODO
			return res, nil
		}
	}
	// never here
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
