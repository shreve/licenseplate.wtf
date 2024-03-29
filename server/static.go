package server

import (
	"embed"
	"io/fs"
	"net/http"
	"os"
	"path"
	"runtime"
)

//go:embed static/*.css static/*.js static/*.json static/*.png
var files embed.FS
var dev = os.Getenv("ENV") == "development"

// https://blog.carlmjohnson.net/post/2021/how-to-use-go-embed/#website-files
var StaticFiles = func() http.Handler {
	var filesys http.FileSystem

	if dev {
		// Get the current file and find static dir relative to it.
		_, file, _, _ := runtime.Caller(1)
		f := path.Join(path.Dir(file), "static")
		filesys = http.FS(os.DirFS(f))
	} else {
		f, _ := fs.Sub(files, "static")
		filesys = http.FS(f)
	}

	return http.FileServer(filesys)
}()
