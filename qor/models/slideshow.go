package models

import (
	"github.com/jinzhu/gorm"
)

type Slideshow struct {
	gorm.Model

	Name     string
	Interval int
	Slides   []SlideshowSlide
}

type SlideshowSlide struct {
	gorm.Model

	SlideshowID uint
	Image       SimpleImageStorage `sql:"type:varchar(4096)"`
	Alt         string
}
