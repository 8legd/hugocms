package models

import (
	"github.com/jinzhu/gorm"
)

type Video struct {
	gorm.Model

	Name           string
	YouTubeID      string
	Title          string
	Length         string
	ThumbnailImage VideoImageStorage `sql:"type:varchar(4096)"`
	Alt            string
}

func (v *Video) DisplayName() string {
	return v.Name
}
