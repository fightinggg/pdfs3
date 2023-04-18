package httpwebdav

import (
	"context"
	"fmt"
	"golang.org/x/net/webdav"
	"log"
	"net/http"
	"os"
)

type HttpWebDav struct {
	webdav.Handler
}

func (u *HttpWebDav) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	if r.Method == "GET" && handleDirList(u.Handler.FileSystem, w, r) {
		return
	}

	u.Handler.ServeHTTP(w, r)
}

func handleDirList(fs webdav.FileSystem, w http.ResponseWriter, req *http.Request) bool {
	ctx := context.Background()
	f, err := fs.OpenFile(ctx, req.URL.Path, os.O_RDONLY, 0)
	if err != nil {
		return false
	}
	defer f.Close()
	if fi, _ := f.Stat(); fi != nil && !fi.IsDir() {
		return false
	}
	dirs, err := f.Readdir(-1)
	if err != nil {
		log.Print(w, "Error reading directory", http.StatusInternalServerError)
		return false
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if len(dirs) == 0 {
		fmt.Fprintf(w, "<pre>Nothing Here</pre>")
	} else {
		fmt.Fprintf(w, "<pre>\n")
		for _, d := range dirs {
			name := d.Name()
			if d.IsDir() {
				name += "/"
			}
			fmt.Fprintf(w, "<a href=\"%s\">%s</a>\n", name, name)
		}
		fmt.Fprintf(w, "</pre>\n")
	}
	return true
}
