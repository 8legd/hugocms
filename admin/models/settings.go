package models

import (
	"github.com/jinzhu/gorm"
)

type Settings struct {
	gorm.Model

	Logo           SettingsLogo
	ContactDetails SettingsContactDetails
	Strapline      SettingsStrapline
	Sidebar        SettingsSidebar
	Copyright      string
	Footer         string `sql:"size:2000"`
}

type SettingsLogo struct {
	gorm.Model

	SettingsID uint

	Image LogoImageStorage `sql:"type:varchar(4096)"`
}

type SettingsContactDetails struct {
	gorm.Model

	SettingsID uint

	Heading string
	Tel     string
	Email   string
}

type SettingsStrapline struct {
	gorm.Model

	SettingsID uint

	TextContent string `sql:"size:2000"`
	TextMobile  string
	Image       SimpleImageStorage `sql:"type:varchar(4096)"`
	Link        string
}

type SettingsSidebar struct {
	gorm.Model

	SettingsID uint
	VideoID    uint
	Video      Video
	Images     []SidebarImage
}

type SidebarImage struct {
	gorm.Model

	SettingsSidebarID uint
	Image             SimpleImageStorage `sql:"type:varchar(4096)"`
}
