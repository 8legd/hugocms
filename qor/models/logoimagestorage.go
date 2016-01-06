package models

import (
	"github.com/qor/qor/media_library"
)

type LogoImageStorage struct{ SimpleImageStorage }

func (LogoImageStorage) GetSizes() map[string]media_library.Size {
	return map[string]media_library.Size{
		"mobile":  {Width: 75, Height: 50},
		"desktop": {Width: 254, Height: 170},
	}
}
