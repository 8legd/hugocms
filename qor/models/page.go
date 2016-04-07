package models

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/8legd/hugocms/config"

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

	ContentColumns []PageContentColumn

	Links []PageLink
}

type PageMeta struct {
	gorm.Model

	PageID uint

	PageTitle   string
	Description string
}

type PageContentColumn struct {
	gorm.Model

	PageID uint

	Width            string
	Heading          string
	TextContent      string `sql:"size:2000"`
	ImageContent     []PageContentColumnImage
	Link             string
	VideoContent     PageContentColumnVideo
	SlideshowContent PageContentColumnSlideshow
}

type PageContentColumnImage struct {
	gorm.Model

	PageContentColumnID uint
	Image               ContentImageStorage `sql:"type:varchar(4096)"`
	Alt                 string
	Alignment           string
}

type PageContentColumnVideo struct {
	gorm.Model

	PageContentColumnID uint
	VideoID             uint
	Video               Video
}

type PageContentColumnSlideshow struct {
	gorm.Model

	PageContentColumnID uint
	SlideshowID         uint
	Slideshow           Slideshow
}

type PageLink struct {
	gorm.Model

	PageID uint

	LinkText string `sql:"size:2000"`
	Link     string
}

func slug(s string) string {
	if s == "" {
		return ""
	}
	return strings.Trim(slugger.ReplaceAllString(strings.ToLower(s), "-"), "-")
}

func (p *Page) Slug() string {
	return slug(p.Name)
}

func (p *Page) PrevSlug() string {
	return slug(p.prevName)
}

func (p *Page) AfterFind() error {
	// handle renames
	p.prevPath = p.Path
	p.prevName = p.Name
	return nil
}

func (p *Page) AfterSave() error {
	// handle renames
	if p.prevPath != "" && (p.prevPath != p.Path || p.prevName != p.Name) {
		// Remove content file from Hugo but rename data files in case we ever need to restore :)
		// TODO use hugo config to get content dir
		filename := "content" + p.prevPath + p.PrevSlug() + ".json"
		if err := os.Remove(filename); err != nil {
			return err
		}
		// TODO use hugo config to get data dir
		filename = "data" + p.prevPath + p.PrevSlug() + ".json"
		if err := os.Rename(filename, filename+".deleted_at_"+time.Now().Format("20060102150405")); err != nil {
			return err
		}
	}
	return p.syncWrite()
}

func (p *Page) AfterDelete() error {
	// Remove content file from Hugo but rename data files in case we ever need to restore :)
	// TODO use hugo config to get content dir
	filename := "content" + p.Path + p.Slug() + ".json"
	if err := os.Remove(filename); err != nil {
		return err
	}
	// TODO use hugo config to get data dir
	filename = "data" + p.Path + p.Slug() + ".json"
	if err := os.Rename(filename, filename+".deleted_at_"+time.Now().Format("20060102150405")); err != nil {
		return err
	}
	return nil
}

// Syncs creation and update events for a page with Hugo
func (p *Page) syncWrite() error {

	// If we have one, fetch the associated Video's model
	// (We need to do this because of the way the relationship is for SettingsIntroVideo > Video)
	for _, c := range p.ContentColumns {
		if c.VideoContent.VideoID > 0 {
			var video Video
			config.DB.First(&video, c.VideoContent.VideoID)
			c.VideoContent.Video = video
		}
	}

	var path = p.Path + p.Slug()
	output, err := json.MarshalIndent(p, "", "  ")
	if err != nil {
		return err
	}
	// Write the data file for Hugo
	// TODO use hugo config to get data dir
	dataFile := "data" + path + ".json"
	// If required, create data dir first
	if _, err := os.Stat("./data"); os.IsNotExist(err) {
		err = os.MkdirAll("./data", os.ModePerm)
		if err != nil {
			return err
		}
	}
	err = ioutil.WriteFile(dataFile, output, 0644)
	if err != nil {
		return err
	}
	// Write the content file for Hugo
	menu := make(map[string]map[string]uint)
	if p.MenuWeight > 0 {
		menuWeight := make(map[string]uint)
		menuWeight["weight"] = p.MenuWeight
		if p.Path != "" && p.Path != "/" {
			menuName := slug(p.Path)
			menu[menuName] = menuWeight
		}
	}
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
	if _, err := os.Stat("./content"); os.IsNotExist(err) {
		err = os.MkdirAll("./content", os.ModePerm)
		if err != nil {
			return err
		}
	}
	err = ioutil.WriteFile(contentFile, content, 0644)
	if err != nil {
		return err
	}

	return nil
}
