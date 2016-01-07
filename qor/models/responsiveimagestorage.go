package models

import (
	"github.com/qor/media_library"
)

type ResponsiveImageStorage struct{ media_library.FileSystem }

func (ResponsiveImageStorage) GetSizes() map[string]media_library.Size {
	return map[string]media_library.Size{
		"mobile":  {Width: 320, Height: 320},
		"tablet":  {Width: 640, Height: 640},
		"desktop": {Width: 1280, Height: 1280},
	}
}
