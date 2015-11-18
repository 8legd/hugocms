package models

import (
	"github.com/jinzhu/gorm"

	"github.com/qor/slug"
)

type Page struct {
	gorm.Model

	Path string

	Name         string
	NameWithSlug slug.Slug

	SEO PageMeta

	ContentRows []PageContentRow
}

type PageMeta struct {
	gorm.Model

	PageID uint

	PageTitle   string
	Description string
}

type PageContentRow struct {
	gorm.Model

	PageID         uint
	Page           Page
	ContentColumns []PageContentColumn
}

type PageContentColumn struct {
	gorm.Model

	PageContentRowID  uint
	ContentRow        PageContentRow
	Heading           string
	TextContent       string             `sql:"size:2000"`
	Image             SimpleImageStorage `sql:"type:varchar(4096)"`
	ImageOptions      string
	Link              string
	VideoID           uint
	Video             Video
	VideoOptions      string
	Slideshow         []PageSlideshowImage
	SlideshowInterval int
}

type PageSlideshowImage struct {
	gorm.Model

	PageContentColumnID uint
	Image               SimpleImageStorage `sql:"type:varchar(4096)"`
}
