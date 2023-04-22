package pdfs_lower_api

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"pdfs3/pdfs/pdfs_blocks"
	"pdfs3/pdfs/pdfs_err"
	"pdfs3/pdfs/storage"
	"strings"
)

type SimplePdfsFileInfo struct {
	Filename       string `json:"filename"`
	FileBlockIndex int    `json:"fileBlockIndex"` // file in which block ?
	FileBlockAlias int    `json:"fileBlockAlias"` // file in this block with which name ?
	IsDir          bool   `json:"isDir"`
	Size           int    `json:"size"`
}

type BlockInfo struct {
	EmptyCap int `json:"emptyCap"`
}

type BlockFsMeta struct {
	// all blockInfo
	AllBlockInfo []BlockInfo `json:"allBlockInfo"`

	// files Tree, PDFS v3.0 only want to support a little file system with no more than 10000 files
	// so were put this tree as a file list
	AllFileList []SimplePdfsFileInfo `json:"allFileList"`
}

type SimplePdfs struct {
	storage   storage.Storage
	metaCache BlockFsMeta
}

func NewSimplePdfsFromStorage(storage storage.Storage) (SimplePdfs, error) {

	res := SimplePdfs{
		storage:   storage,
		metaCache: BlockFsMeta{},
	}

	// pdfs v3.0 only use the first block as head block
	head, err := io.ReadAll(storage.Read(0))
	if err != nil {
		return SimplePdfs{}, err
	}

	err = json.Unmarshal(head, &res.metaCache)
	if err != nil {
		return SimplePdfs{}, err
	}

	return res, nil
}

func NewSimplePdfsFromStorageAndIgnoreAllError(storage storage.Storage) SimplePdfs {
	log.Printf("WARN: ignore all filesystem error , some files may loss")

	infos := make([]BlockInfo, storage.TotalSize()/storage.BlockByteSize())
	for i := range infos {
		if i == 0 {
			infos[i].EmptyCap = 0
		} else {
			infos[i].EmptyCap = storage.BlockByteSize() - 1 - 8 - 1
		}
	}
	return SimplePdfs{
		storage: storage,
		metaCache: BlockFsMeta{
			AllBlockInfo: infos,
			AllFileList:  []SimplePdfsFileInfo{},
		},
	}
}

func (s *SimplePdfs) ReadFile(filename string, start, len int64) (io.Reader, error) {
	for _, file := range s.metaCache.AllFileList {
		if file.Filename == filename {
			return parseBlock(s.storage, file.FileBlockIndex, file.FileBlockAlias)
		}
	}
	return nil, pdfs_err.PdfsNotFoundError()
}

func (s *SimplePdfs) WriteFile(filename string, start, end int64, reader io.Reader) error {
	data, err := io.ReadAll(reader)
	if err != nil {
		return err
	}

	// update file
	for _, file := range s.metaCache.AllFileList {
		if file.Filename == filename {
			_, cap, err := updateBlock(s.storage, file.FileBlockIndex, file.FileBlockAlias, data)
			s.metaCache.AllBlockInfo[file.FileBlockIndex].EmptyCap = cap
			// TODO may write or not
			return err
		}
	}

	// create file
	for i, info := range s.metaCache.AllBlockInfo {
		if info.EmptyCap > len(data) {

			alias, cap, err := updateBlock(s.storage, i, -1, data)
			if err != nil {
				return err
			}
			s.metaCache.AllBlockInfo[i].EmptyCap = cap

			var simplePdfsFileInfo SimplePdfsFileInfo
			simplePdfsFileInfo.FileBlockIndex = i
			simplePdfsFileInfo.Filename = filename
			simplePdfsFileInfo.FileBlockAlias = alias
			simplePdfsFileInfo.IsDir = false
			//simplePdfsFileInfo.Size = len(data)

			s.metaCache.AllFileList = append(s.metaCache.AllFileList, simplePdfsFileInfo)
			// TODO may write or not
			return nil
		}
	}

	return pdfs_err.PdfsOutOfDiskError()

}

func (s *SimplePdfs) RemoveFile(filename string) {
	return
}

func (s *SimplePdfs) listDir(path string) (res []SimplePdfsFileInfo) {
	if !strings.HasSuffix(path, "/") {
		path = path + "/"
	}
	mp := map[string]SimplePdfsFileInfo{}
	for _, file := range s.metaCache.AllFileList {
		if strings.HasPrefix(file.Filename, path) {
			suffix := file.Filename[len(path):]
			if strings.Contains(suffix, "/") {
				name := suffix[:strings.Index(suffix, "/")]
				mp[name] = SimplePdfsFileInfo{
					Filename: name,
					IsDir:    true,
				}
			} else if suffix != "" {
				mp[suffix] = SimplePdfsFileInfo{
					Filename:       suffix,
					FileBlockIndex: file.FileBlockIndex,
					FileBlockAlias: file.FileBlockAlias,
					IsDir:          false,
				}
			}
		}
	}
	for _, info := range mp {
		res = append(res, info)
	}
	return
}

func (s *SimplePdfs) mkdir(name string) {
	if !strings.HasSuffix(name, "/") {
		name = name + "/"
	}
	s.metaCache.AllFileList = append(s.metaCache.AllFileList, SimplePdfsFileInfo{
		Filename: name,
		IsDir:    true,
	})
	// TODO write
}

func parseBlock(blockfs storage.Storage, index, alias int) (io.Reader, error) {
	// all file data here, do not need to read others block
	allData, err := io.ReadAll(blockfs.Read(index))
	if err != nil {
		return nil, err
	}

	if allData[0] == 0 {
		all, readErr := pdfs_blocks.ReadAll(allData)
		if readErr != nil {
			return nil, readErr
		}
		res, readErr := all.Read(alias)
		if readErr != nil {
			return nil, readErr
		}
		return bytes.NewReader(res), nil
	}

	return nil, nil
}

func updateBlock(storage storage.Storage, index, alias int, data []byte) (newAlias, cap int, err error) {

	// all file data here, do not need to read others block
	allData, err := io.ReadAll(storage.Read(index))
	if err != nil {
		return 0, 0, err
	}

	if allData[0] == 0 {
		all, readErr := pdfs_blocks.ReadAll(allData)
		if readErr != nil {
			return 0, 0, readErr
		}

		newAlias, cap = all.Write(alias, data)

		writeAll, readErr := all.WriteAll()
		if readErr != nil {
			return 0, 0, readErr
		}

		err := storage.Write(index, bytes.NewReader(writeAll))
		if err != nil {
			return 0, 0, err
		}
	}

	return newAlias, storage.BlockByteSize() - 1 - 2*4 - 1 - cap, nil
}
