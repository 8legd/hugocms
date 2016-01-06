package models

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/8legd/hugocms/config"

	"github.com/jinzhu/gorm"
)

type Settings struct {
	gorm.Model

	Logo           SettingsLogo
	ContactDetails SettingsContactDetails
	Header         SettingsHeader
	IntroVideo     SettingsIntroVideo
	Sidebar        SettingsSidebar
	Copyright      string
	Footer         string `sql:"size:2000"`
}

type SettingsLogo struct {
	gorm.Model

	SettingsID uint

	Text  string
	Image LogoImageStorage `sql:"type:varchar(4096)"`
	Alt   string
	Link  string
}

type SettingsContactDetails struct {
	gorm.Model

	SettingsID uint

	Title string
	Tel   string
	Email string
}

type SettingsHeader struct {
	gorm.Model

	SettingsID uint

	TextContent string `sql:"size:2000"`
	TextMobile  string
	Image       SimpleImageStorage `sql:"type:varchar(4096)"`
	Alt         string
	Link        string
}

type SettingsIntroVideo struct {
	gorm.Model
	SettingsID uint
	VideoID    uint
	Video      Video
	Title      string
	Length     string
	Image      SimpleImageStorage `sql:"type:varchar(4096)"`
	Alt        string
	SEO        SettingsIntroVideoMeta
}

type SettingsIntroVideoMeta struct {
	gorm.Model

	SettingsIntroVideoID uint

	PageTitle   string
	Description string
}

type SettingsSidebar struct {
	gorm.Model

	SettingsID uint
	Images     []SettingsSidebarImage
}

type SettingsSidebarImage struct {
	gorm.Model

	SettingsSidebarID uint
	Image             SimpleImageStorage `sql:"type:varchar(4096)"`
	Alt               string
	Link              string
}

func (s *Settings) AfterSave() (err error) {

	// If we have one, fetch the associated IntroVideo's Video model
	// (We need to do this because of the way the releationship is for SettingsIntroVideo > Video)
	if s.IntroVideo.VideoID > 0 {
		var video Video
		config.QOR.DB.First(&video, s.IntroVideo.VideoID)
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
		err = os.MkdirAll(filepath.Dir(contentFile), 0777)
		if err != nil {
			return err
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

	output, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return err
	}
	err = ioutil.WriteFile("config.json", output, 0644)
	if err != nil {
		return err
	}

	return
}
