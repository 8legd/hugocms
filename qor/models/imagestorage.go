package models

import (
	"io"
	"os"

	"github.com/qor/media_library"
)

var fs media_library.FileSystem
var opt media_library.Option

func init() {
	opt = make(map[string]string)
	opt["PATH"] = "./static"
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

type LogoImageStorage struct{ SimpleImageStorage }

func (LogoImageStorage) GetSizes() map[string]media_library.Size {
	return map[string]media_library.Size{
		"mobile":  {Width: 75, Height: 50},
		"desktop": {Width: 254, Height: 170},
	}
}

type VideoImageStorage struct{ SimpleImageStorage }

type HeaderImageStorage struct{ SimpleImageStorage }

type SidebarImageStorage struct{ SimpleImageStorage }

type ContentImageStorage struct{ SimpleImageStorage }
