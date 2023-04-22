package pdfs_lower_api

import (
	"io"
	"io/fs"
	"log"
	"time"
)

type Menu struct {
	FileName   string    `json:"fileName"`   // < 255B
	FileSize   int64     `json:"fileSize"`   // = 8B
	FileUpdate time.Time `json:"fileUpdate"` // =10B
	IsDirV     bool      `json:"isDir"`      // =1B
}

func (p *Menu) Name() string {
	return p.FileName
}

func (p *Menu) Size() int64 {
	return p.FileSize
}

func (p *Menu) Mode() fs.FileMode {
	log.Printf("WARN : mode=777")
	return 777
}

func (p *Menu) ModTime() time.Time {
	return p.FileUpdate
}

func (p *Menu) IsDir() bool {
	return p.IsDirV
}

func (p *Menu) Sys() interface{} {
	log.Printf("WARN : not Sys")
	return nil
}

type LowerApi interface {
	Ls(path string) ([]Menu, error)
	Read(path string) (io.Reader, error)
	Mkdir(path string) error
	Write(path string, offset int64, len int64, reader io.Reader) error
	Delete(path string) error
	Stat(name string) (*Menu, error)
	Rename(name string, name2 string) error
}
