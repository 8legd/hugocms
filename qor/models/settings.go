package models

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/8legd/hugocms/config"

	"github.com/jinzhu/gorm"
)

type Settings struct {
	gorm.Model

	Logo           SettingsLogo
	ContactDetails SettingsContactDetails
	Header         SettingsHeader
	CallToAction   SettingsCallToAction
	IntroVideo     SettingsIntroVideo
	Sidebar        []SettingsSidebarContent
	Copyright      string
	Footer         string `sql:"size:2000"`
}

type SettingsLogo struct {
	gorm.Model

	SettingsID uint

	Image LogoImageStorage `sql:"type:varchar(4096)"`
	Alt   string
}

type SettingsContactDetails struct {
	gorm.Model

	SettingsID uint

	Title               string
	Tel                 string
	Email               string
	OpeningHoursDesktop string `sql:"size:2000"`
	OpeningHoursMobile  string
}

type SettingsHeader struct {
	gorm.Model

	SettingsID uint
	Image      HeaderImageStorage `sql:"type:varchar(4096)"`
	Alt        string
	Link       string
}

type SettingsCallToAction struct {
	gorm.Model

	SettingsID uint

	ActionText string `sql:"size:2000"`
	Link       string
}

type SettingsIntroVideo struct {
	gorm.Model
	SettingsID uint
	VideoID    uint
	Video      Video
	SEO        SettingsIntroVideoMeta
}

type SettingsIntroVideoMeta struct {
	gorm.Model

	SettingsIntroVideoID uint

	PageTitle   string
	Description string
}

type SettingsSidebarContent struct {
	gorm.Model

	SettingsID uint

	Image SidebarImageStorage `sql:"type:varchar(4096)"`
	Alt   string
	Link  string
}

func (s *Settings) AfterSave() error {

	// If we have one, fetch the associated IntroVideo's Video model
	// (We need to do this because of the way the relationship is for SettingsIntroVideo > Video)
	if s.IntroVideo.VideoID > 0 {
		var video Video
		config.DB.First(&video, s.IntroVideo.VideoID)
		s.IntroVideo.Video = video

		// Save to Intro Video content page
		content, err := json.MarshalIndent(
			struct {
				Title       string `json:"Title"`
				Description string `json:"Description"`
				Date        string `json:"Date"`
			}{
				s.IntroVideo.SEO.PageTitle,
				s.IntroVideo.SEO.Description,
				s.IntroVideo.CreatedAt.Format("2006-01-02T15:04:05Z"),
			},
			"",
			"  ",
		)
		if err != nil {
			return err
		}

		// TODO use hugo config to get content dir
		contentFile := "content/intro-video.html"
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

	}

	// Update Hugo config - this is where our settings are stored (as custom Params)
	config := config.Hugo
	config.Params = make(map[string]interface{})
	config.Params["Settings"] = s
	rs := make(map[string]map[string]interface{})
	config.Params["ResponsiveSettings"] = rs
	if s.Logo.Image.Url != "" {
		rsl := make(map[string]interface{})
		rs["Logo"] = rsl
		ext := filepath.Ext(s.Logo.Image.Url)
		for k, _ := range s.Logo.Image.CropOptions {
			rsl[k] = strings.Replace(s.Logo.Image.Url, ext, "."+k+ext, 1)
		}
	}

	output, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return err
	}
	err = ioutil.WriteFile("config.json", output, 0644)
	if err != nil {
		return err
	}

	return nil
}
