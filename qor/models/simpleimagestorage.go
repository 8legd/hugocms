package models

import (
	"io"
	"os"

	"github.com/qor/qor/media_library"
)

var fs media_library.FileSystem
var opt media_library.Option

func init() {
	opt = make(map[string]string)
	opt["path"] = "./static"
}

type SimpleImageStorage struct{ media_library.Base }

func (s SimpleImageStorage) GetFullPath(url string, option *media_library.Option) (path string, err error) {
	return fs.GetFullPath(url, &opt)
}

func (s SimpleImageStorage) Store(url string, option *media_library.Option, reader io.Reader) error {
	return fs.Store(url, &opt, reader)
}

func (s SimpleImageStorage) Retrieve(url string) (*os.File, error) {
	return fs.Retrieve(url)
}
