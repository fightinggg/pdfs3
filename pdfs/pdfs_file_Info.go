package pdfs

import (
	"io/fs"
	"log"
	"time"
)

type PdfsFileInfo struct {
	info SimplePdfsFileInfo
}

func (p *PdfsFileInfo) Name() string {
	return p.info.Filename
}

func (p *PdfsFileInfo) Size() int64 {
	//log.Printf("WARN: do not support filesize")
	//return int64(p.info.Size)
	return 0
}

func (p *PdfsFileInfo) Mode() fs.FileMode {
	log.Printf("WARN : mode=777")
	return 777
}

func (p *PdfsFileInfo) ModTime() time.Time {
	log.Printf("WARN : do not support modtime")
	return time.Now()
}

func (p *PdfsFileInfo) IsDir() bool {
	return p.info.IsDir
}

func (p *PdfsFileInfo) Sys() interface{} {
	//TODO implement me
	panic("implement me")
}
