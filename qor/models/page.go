package models

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/jinzhu/gorm"
)

var slugger = regexp.MustCompile("[^a-z0-9]+")

type Page struct {
	gorm.Model

	Path       string
	prevPath   string
	MenuWeight uint

	Name     string
	prevName string

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

func (p *Page) Slug() string {
	if p.Name == "" {
		return ""
	}
	return strings.Trim(slugger.ReplaceAllString(strings.ToLower(p.Name), "-"), "-")
}

func (p *Page) AfterFind() (err error) {
	// handle renames
	p.prevPath = p.Path
	p.prevName = p.Name
	return
}

func (p *Page) AfterSave() (err error) {
	// handle renames
	if p.prevPath != p.Path || p.prevName != p.Name {
		p.syncRemoveRef()
	}
	if false {
		err = errors.New("Page AfterSave Error")
	} else {
		p.syncWrite()
	}
	return
}

func (p *Page) AfterDelete() (err error) {
	if false {
		err = errors.New("Page AfterDelete Error")
	} else {
		p.syncRemoveRef()
	}
	return
}

// Syncs creation and update events for a page with Hugo
func (p *Page) syncWrite() (err error) {
	var path = p.Path + p.Slug()
	output, err := json.MarshalIndent(p, "", "  ")
	if err != nil {
		return err
	}
	// Write the data file for Hugo
	// TODO use hugo config to get data dir
	dataFile := "data" + path + ".json"
	// If required, create content dir first
	err = os.MkdirAll(filepath.Dir(dataFile), 0777)
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(dataFile, output, 0644)
	if err != nil {
		return err
	}
	// Write the content file for Hugo
	// TODO if p.MenuWeight < 1 hidden?
	menuWeight := make(map[string]uint)
	menuWeight["weight"] = p.MenuWeight
	menu := make(map[string]map[string]uint)
	menu["about_us"] = menuWeight
	content, err := json.MarshalIndent(
		struct {
			Title       string                     `json:"Title"`
			Description string                     `json:"Description"`
			Date        string                     `json:"Date"`
			Menu        map[string]map[string]uint `json:"Menu"`
		}{
			p.SEO.PageTitle,
			p.SEO.Description,
			p.CreatedAt.Format("2006-01-02T15:04:05Z"),
			menu,
		},
		"",
		"  ",
	)
	// TODO use hugo config to get content dir
	contentFile := "content" + path + ".json"
	// If required, create content dir first
	err = os.MkdirAll(filepath.Dir(contentFile), 0777)
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(contentFile, content, 0644)
	if err != nil {
		return err
	}

	return
}

// Syncs rename or deletion of a page with Hugo
func (p *Page) syncRemoveRef() (err error) {
	// TODO just remove content file from Hugo (data files can remain)
	// this way the page will in effect be un-published
	// TODO if after removing the content file the section directory is empty then
	// also remove it to clear up (again the sections will remain in the data dir)
	// TODO use hugo config to get content dir
	var filename = "content" + p.Path + p.Slug() + ".json"
	fmt.Printf("\nTODO remove %s \n", filename)
	fmt.Printf("TODO and if required remove empty section dir %s \n", p.Path)
	return
}
