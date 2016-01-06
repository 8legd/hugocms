package models

import (
	"github.com/jinzhu/gorm"
)

type Video struct {
	gorm.Model

	Name      string
	YouTubeID string
}

func (v *Video) DisplayName() string {
	return v.Name
}
