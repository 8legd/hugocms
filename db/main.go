// +build ignore

package main

import (
	"github.com/8legd/hugocms/admin/models"
	"github.com/8legd/hugocms/db"
)

func main() {
	tables := []interface{}{

		&models.User{},
		&models.UserImage{},

		&models.Settings{},
		&models.SettingsLogo{},
		&models.SettingsContactDetails{},
		&models.SettingsStrapline{},
		&models.SettingsSidebar{},
		&models.SidebarImage{},

		&models.Video{},

		&models.Page{},
		&models.PageMeta{},
		&models.PageContentRow{},
		&models.PageContentColumn{},
		&models.PageSlideshowImage{},
	}

	for _, table := range tables {
		if err := db.DB.DropTableIfExists(table).Error; err != nil {
			panic(err)
		}
		if err := db.DB.AutoMigrate(table).Error; err != nil {
			panic(err)
		}
	}

}
