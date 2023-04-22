package pdfs_lower_api

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"log"
	"pdfs3/pdfs/pdfs_err"
	"pdfs3/pdfs/storage"
	"strings"
	"time"
)

type MenuBlock struct {
	PathPrefix    string              `json:"pathPrefix"`
	NextMenuIndex int                 `json:"nextMenuIndex"`
	Menus         []*MenuWithDataMenu `json:"menus"`
}

type MenuWithDataMenu struct {
	PathIndex     int       `json:"pathIndex"`
	FileName      string    `json:"fileName"`      // < 255B
	FileSize      int64     `json:"fileSize"`      // = 8B
	FileUpdate    time.Time `json:"fileUpdate"`    // =10B
	IsDir         bool      `json:"isDir"`         // =1B
	DirBlockIndex int       `json:"dirBlockIndex"` // =2B
	FileData      []byte    `json:"fileData"`      // <= 1KB
}

type PdfsLowerApi struct {
	Storage storage.Storage
}

func (api *PdfsLowerApi) Rename(name string, name2 string) error {
	//TODO implement me
	return errors.New("NO")
}

func (api *PdfsLowerApi) Stat(path string) (*Menu, error) {
	log.Printf("Stat %s", path)

	if strings.HasSuffix(path, "/") {
		path = path[:len(path)-1]
	}
	var nextIndex = 0
	for true {
		menus, err := api.readMenus(nextIndex)
		if err != nil {
			return nil, err
		}

		nextIndex = menus.NextMenuIndex

		for _, menu := range menus.Menus {
			fileName := menus.PathPrefix + menu.FileName
			if fileName == path {
				return &Menu{
					FileName:   fileName,
					FileSize:   menu.FileSize,
					FileUpdate: menu.FileUpdate,
					IsDirV:     menu.IsDir,
				}, nil
			}
		}

		if nextIndex <= 0 {
			// TODO
			return nil, pdfs_err.PdfsNotFoundError()
		}
	}
	// never here
	return nil, nil
}

func (api *PdfsLowerApi) writeMenus(blockIndex int, menus *MenuBlock) error {
	marshal, err := json.Marshal(menus)
	if err != nil {
		return err
	}
	return api.Storage.Write(blockIndex, bytes.NewReader(marshal))
}

func (api *PdfsLowerApi) readMenus(blockIndex int) (*MenuBlock, error) {

	read := api.Storage.Read(blockIndex)
	all, err := io.ReadAll(read)
	if err != nil {
		return nil, err
	}

	var menus MenuBlock
	err = json.Unmarshal(all, &menus)
	if err != nil {
		return &MenuBlock{
			PathPrefix:    "",
			NextMenuIndex: 0,
			Menus: []*MenuWithDataMenu{
				{
					FileName:      "",
					FileSize:      0,
					FileUpdate:    time.Now(),
					IsDir:         true,
					DirBlockIndex: 0,
					FileData:      nil,
				},
			},
		}, nil
	}

	return &menus, nil
}

func (api *PdfsLowerApi) Ls(path string) ([]Menu, error) {
	log.Printf("Ls %s", path)

	if !strings.HasSuffix(path, "/") {
		path = path + "/"
	}
	var res []Menu
	var nextIndex = 0
	for true {
		menus, err := api.readMenus(nextIndex)
		if err != nil {
			return nil, err
		}

		nextIndex = menus.NextMenuIndex

		for _, menu := range menus.Menus {
			fileName := menus.PathPrefix + menu.FileName
			if fileName == path && menu.DirBlockIndex > 0 {
				nextIndex = menu.DirBlockIndex
				break
			} else if strings.HasPrefix(fileName, path) && !strings.Contains(menu.FileName[len(path):], "/") {
				res = append(res, Menu{
					FileName:   fileName[len(path):],
					FileSize:   menu.FileSize,
					FileUpdate: menu.FileUpdate,
					IsDirV:     menu.IsDir,
				})
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

func (api *PdfsLowerApi) Read(path string) (io.Reader, error) {
	log.Printf("read %s", path)
	var nextIndex = 0
	for true {
		menus, err := api.readMenus(nextIndex)
		if err != nil {
			return nil, err
		}

		nextIndex = menus.NextMenuIndex

		for _, menu := range menus.Menus {
			fileName := menus.PathPrefix + menu.FileName
			if fileName == path {
				// TODO multiblock
				return bytes.NewReader(menu.FileData), nil
			} else if strings.HasPrefix(path, fileName) && menu.IsDir && menu.DirBlockIndex > 0 {
				nextIndex = menu.DirBlockIndex
				break
			}
		}

		if nextIndex <= 0 {
			return nil, pdfs_err.PdfsNotFoundError()
		}
	}
	return nil, nil
}

func (api *PdfsLowerApi) Mkdir(path string) error {
	log.Printf("mkdir %s", path)
	var nextIndex = 0
	for true {
		thisIndex := nextIndex
		menus, err := api.readMenus(nextIndex)
		if err != nil {
			return err
		}

		nextIndex = menus.NextMenuIndex

		for _, menu := range menus.Menus {
			fileName := menus.PathPrefix + menu.FileName
			if fileName == path {
				return nil
			} else if strings.HasPrefix(path, fileName) && menu.IsDir && menu.DirBlockIndex > 0 {
				nextIndex = menu.DirBlockIndex
				break
			}
		}

		if nextIndex <= 0 {
			menus.Menus = append(menus.Menus, &MenuWithDataMenu{
				FileName:      path[len(menus.PathPrefix):],
				FileSize:      0,
				FileUpdate:    time.Now(),
				IsDir:         true,
				DirBlockIndex: 0,
				FileData:      nil,
			})
			return api.writeMenus(thisIndex, menus)
		}
	}
	return nil

}

func (api *PdfsLowerApi) Write(path string, offset int64, size int64, reader io.Reader) error {
	log.Printf("write %s", path)
	data, err := io.ReadAll(reader)
	if err != nil {
		return err
	}

	var nextIndex = 0
	for true {
		thisIndex := nextIndex
		menus, err := api.readMenus(nextIndex)
		if err != nil {
			return err
		}

		nextIndex = menus.NextMenuIndex

		for _, menu := range menus.Menus {
			fileName := menus.PathPrefix + menu.FileName
			if fileName == path {
				// update
				if offset+size < menu.FileSize {
					menu.FileData = append(menu.FileData[:offset], append(data, menu.FileData[offset+size:]...)...)
				} else {
					menu.FileData = append(menu.FileData[:offset], data...)
				}
				menu.FileSize = int64(len(menu.FileData))
				return api.writeMenus(thisIndex, menus)
			} else if strings.HasPrefix(path, fileName) && menu.IsDir && menu.DirBlockIndex > 0 {
				nextIndex = menu.DirBlockIndex
				break
			}
		}

		if nextIndex <= 0 {
			// create
			if offset != 0 {
				return errors.New("FUck")
			}
			menus.Menus = append(menus.Menus, &MenuWithDataMenu{
				FileName:      path[len(menus.PathPrefix):],
				FileSize:      int64(len(data)),
				FileUpdate:    time.Now(),
				IsDir:         false,
				DirBlockIndex: 0,
				FileData:      data,
			})
			return api.writeMenus(thisIndex, menus)
		}
	}
	return nil

}

func (api *PdfsLowerApi) Delete(path string) error {
	log.Printf("delete %s", path)
	var nextIndex = 0
	for true {
		thisIndex := nextIndex
		menus, err := api.readMenus(nextIndex)
		if err != nil {
			return err
		}

		nextIndex = menus.NextMenuIndex

		for i, menu := range menus.Menus {
			fileName := menus.PathPrefix + menu.FileName
			if fileName == path {
				menus.Menus = append(menus.Menus[:i], menus.Menus[i+1:]...)
				return api.writeMenus(thisIndex, menus)
			} else if strings.HasPrefix(path, fileName) && menu.IsDir && menu.DirBlockIndex > 0 {
				nextIndex = menu.DirBlockIndex
				break
			}
		}

		if nextIndex <= 0 {
			return nil
		}
	}
	return nil

}
