package embed

import (
	"embed"
	"io/fs"
	"net/http"
)

//go:embed all:dist
var dist embed.FS

var DistFS http.FileSystem

func init() {
	sub, err := fs.Sub(dist, "dist")
	if err != nil {
		panic("embed: failed to create sub-filesystem: " + err.Error())
	}
	DistFS = &fallbackFS{http.FS(sub)}
}

type fallbackFS struct {
	fs http.FileSystem
}

func (f *fallbackFS) Open(name string) (http.File, error) {
	file, err := f.fs.Open(name)
	if err != nil {
		return f.fs.Open("index.html")
	}
	return file, nil
}
