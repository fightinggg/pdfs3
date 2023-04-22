package main

import (
	"context"
	"golang.org/x/net/webdav"
	"io"
	"log"
	"net/http"
	"pdfs3/httpwebdav"
	"pdfs3/pdfs"
	"pdfs3/pdfs/pdfs_lower_api"
	"pdfs3/pdfs/storage"
	"strconv"
)

func mkdir(newPdfs webdav.FileSystem, dir string) {
	err := newPdfs.Mkdir(context.Background(), dir, 0)
	if err != nil {
		panic(err)
	}
}

func createFile(newPdfs webdav.FileSystem, path string, v string) {
	file, err2 := newPdfs.OpenFile(context.Background(), path, 0, 0)
	if err2 != nil {
		panic(err2.Error())
	}

	_, err2 = io.WriteString(file, v)
	if err2 != nil {
		panic(err2.Error())
	}
}

func addSomeFiles(newPdfs webdav.FileSystem) {

	createFile(newPdfs, "/a.txt", "fuck a.txt")
	createFile(newPdfs, "/b.txt", "fuck b.txt")
	//
	mkdir(newPdfs, "/abc")
	mkdir(newPdfs, "/小黑子")

	createFile(newPdfs, "/小黑子/小黑子.txt", "你干嘛，哎呦")
	createFile(newPdfs, "/abc/a.txt", "fuck abc/a.txt")
	createFile(newPdfs, "/abc/b.txt", "fuck abc/b.txt")

}

func main() {

	port := 8080
	memoryStorage := storage.NewMemoryStorage(1<<20, 1)
	newPdfs := &pdfs.Pdfs{
		LowerApi: &pdfs_lower_api.PdfsLowerApi{
			Storage: memoryStorage,
		},
	}

	//
	//newPdfs = webdav.NewMemFS()
	////
	//file, err2 := newPdfs.OpenFile(context.Background(), "a.txt", 578, 438)
	//print(file, err2)

	//addSomeFiles(newPdfs)

	server := httpwebdav.HttpWebDav{
		Handler: webdav.Handler{
			FileSystem: newPdfs,
			LockSystem: webdav.NewMemLS(),
		},
	}
	log.Printf("server start at port %d\n", port)
	err := http.ListenAndServe(":"+strconv.Itoa(port), &server)
	if err != nil {
		panic(err)
	}
}
