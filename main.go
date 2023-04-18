package main

import (
	"context"
	"golang.org/x/net/webdav"
	"io"
	"log"
	"net/http"
	"pdfs3/httpwebdav"
	"pdfs3/pdfs"
	"pdfs3/pdfs/storage"
	"strconv"
)

func addSomeFiles(newPdfs webdav.FileSystem) {
	file, err2 := newPdfs.OpenFile(context.Background(), "/a.txt", 0, 0)
	if err2 != nil {
		panic(err2.Error())
	}

	_, err2 = io.WriteString(file, "你好")
	if err2 != nil {
		panic(err2.Error())
	}

	file, err2 = newPdfs.OpenFile(context.Background(), "/b.txt", 0, 0)
	if err2 != nil {
		panic(err2.Error())
	}

	_, err2 = io.WriteString(file, "你好b")
	if err2 != nil {
		panic(err2.Error())
	}

	file, err2 = newPdfs.OpenFile(context.Background(), "/小黑子/b.txt", 0, 0)
	if err2 != nil {
		panic(err2.Error())
	}

	_, err2 = io.WriteString(file, "小黑子就是你")
	if err2 != nil {
		panic(err2.Error())
	}

	err := newPdfs.Mkdir(context.Background(), "/abc", 0)
	if err != nil {
		return
	}

}

func main() {

	port := 8080
	memoryStorage := storage.NewMemoryStorage(1<<20, 8)
	newPdfs, err3 := pdfs.NewPdfsFromStorage(memoryStorage)

	if err3 != nil {
		newPdfs = pdfs.NewPdfsFromStorageAndIgnoreAllError(memoryStorage)
	}

	newPdfs.Mkdir(context.Background(), "/abc", 0)

	//addSomeFiles(newPdfs)
	//
	//newPdfs = webdav.NewMemFS()
	////
	//file, err2 := newPdfs.OpenFile(context.Background(), "a.txt", 578, 438)
	//print(file, err2)

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
