package widget

import (
	"github.com/jinzhu/gorm"
	"github.com/qor/admin"
	"github.com/qor/serializable_meta"
)

// QorWidgetSetting default qor widget setting struct
type QorWidgetSetting struct {
	gorm.Model
	Scope string
	Name  string
	serializable_meta.SerializableMeta
}

func findSettingByNameAndKinds(db *gorm.DB, name string, kinds []string) *QorWidgetSetting {
	setting := QorWidgetSetting{}
	if db.Where("name = ? AND kind IN (?)", name, kinds).First(&setting).RecordNotFound() {
		setting.Name = name
		setting.Kind = kinds[0]
		db.Save(&setting)
	}
	return &setting
}

// GetSerializableArgumentResource get setting's argument's resource
func (setting *QorWidgetSetting) GetSerializableArgumentResource() *admin.Resource {
	return GetWidget(setting.Kind).Setting
}
