package models

import (
	"errors"
	sync "github.com/8legd/hugocms/admin/sync"
	"github.com/jinzhu/gorm"
	"github.com/qor/slug"
)

type Page struct {
	gorm.Model

	Path     string
	prevPath string

	Name         string
	NameWithSlug slug.Slug
	prevSlug     string

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

	PageID uint

	ContentColumns []PageContentColumn
}

type PageContentColumn struct {
	gorm.Model

	PageContentRowID uint

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

func (p *Page) AfterFind() (err error) {
	// handle renames
	p.prevPath = p.Path
	p.prevSlug = p.NameWithSlug.Slug
	return
}

func (p *Page) AfterSave() (err error) {
	// handle renames
	if p.prevPath != p.Path || p.prevSlug != p.NameWithSlug.Slug {
		sync.RemovePageRef(p.prevPath, p.prevSlug)
	}
	// TODO check if any images have been created

	if false {
		err = errors.New("Page AfterSave Error")
	} else {
		sync.WritePage(p, p.Path, p.NameWithSlug.Slug)
	}
	return
}

func (p *Page) AfterDelete() (err error) {
	if false {
		err = errors.New("Page AfterDelete Error")
	} else {
		sync.RemovePageRef(p.Path, p.NameWithSlug.Slug)
	}
	return
}
