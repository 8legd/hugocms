package config

import (
	"github.com/jinzhu/gorm"
	"github.com/qor/admin"
	"github.com/qor/i18n"
)

var (
	QOR  QORConfig
	Hugo HugoConfig
	DB   *gorm.DB
	Auth admin.Auth
	I18n *i18n.I18n
)

type QORConfig struct {
	Port     int
	SiteName string
	Paths    []string
}

type HugoConfig struct {
	BaseURL        string `json:"baseurl"`
	StaticDir      string `json:"staticdir"`
	PublishDir     string `json:"publishdir"`
	LanguageCode   string `json:"languageCode"`
	DisableRSS     bool   `json:"disableRSS"`
	MetaDataFormat string
	Menu           map[string]interface{}
	Params         map[string]interface{}
}
