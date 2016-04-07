package config

import (
	"github.com/adrianduke/configr"
	_ "github.com/adrianduke/configr/sources/file/toml"
	"github.com/jinzhu/gorm"
	"github.com/qor/admin"
)

var (
	QOR  QORConfig
	Hugo HugoConfig
	DB   *gorm.DB
	Auth admin.Auth
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

func Setup(qorConfigFile string, hugoConfigFile string, db *gorm.DB, auth admin.Auth) error {

	qorConf := configr.New()
	qorConf.RegisterKey("port", "QOR admin port", 8000)
	qorConf.RegisterKey("sitename", "QOR admin site name", "QOR Admin")
	qorConf.AddSource(configr.NewFile(qorConfigFile))
	if err := qorConf.Parse(); err != nil {
		return err
	}
	port, err := qorConf.Int("port")
	if err != nil {
		return err
	}
	QOR.Port = port

	sitename, err := qorConf.String("sitename")
	if err != nil {
		return err
	}
	QOR.SiteName = sitename

	// As a minumum add the root path for our site
	QOR.Paths = append(QOR.Paths, "/")

	Hugo.MetaDataFormat = "json"

	hugoConf := configr.New()
	hugoConf.RegisterKey("baseurl", "Hugo site baseurl", "/")
	hugoConf.RegisterKey("staticdir", "Hugo site static dir", "static")
	hugoConf.RegisterKey("publishdir", "Hugo site publish dir", "public")
	hugoConf.RegisterKey("languageCode", "Hugo site languageCode", "en")
	hugoConf.RegisterKey("disableRSS", "Hugo site disableRSS", true)
	hugoConf.RegisterKey("menu", "Hugo site menus", make(map[string]interface{}))

	hugoConf.AddSource(configr.NewFile(hugoConfigFile))
	if err := hugoConf.Parse(); err != nil {
		return err
	}

	baseurl, err := hugoConf.String("baseurl")
	if err != nil {
		return err
	}
	Hugo.BaseURL = baseurl

	staticdir, err := hugoConf.String("staticdir")
	if err != nil {
		return err
	}
	Hugo.StaticDir = staticdir

	publishdir, err := hugoConf.String("publishdir")
	if err != nil {
		return err
	}
	Hugo.PublishDir = publishdir

	languageCode, err := hugoConf.String("languageCode")
	if err != nil {
		return err
	}
	Hugo.LanguageCode = languageCode

	disableRSS, err := hugoConf.Bool("disableRSS")
	if err != nil {
		return err
	}
	Hugo.DisableRSS = disableRSS

	rawMenu, err := hugoConf.Get("menu")
	if err != nil {
		return err
	}
	if menu, ok := rawMenu.(map[string]interface{}); ok {
		Hugo.Menu = menu
		// Add additional site paths from main menu items
		if rawMainMenu, ok := menu["main"]; ok {
			if mainMenu, ok := rawMainMenu.([]map[string]interface{}); ok {
				for _, item := range mainMenu {
					if url, ok := item["url"].(string); ok {
						if url != "" && url != "/" {
							QOR.Paths = append(QOR.Paths, url)
						}
					}
				}
			}
		}
	}

	DB = db

	Auth = auth

	return nil

}
