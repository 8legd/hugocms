package config

import (
	"github.com/adrianduke/configr"
	_ "github.com/adrianduke/configr/sources/file/toml"
	"github.com/jinzhu/gorm"
)

var (
	QOR  QORConfig
	Hugo HugoConfig
)

type QORConfig struct {
	Port     int
	SiteName string
	DB       *gorm.DB
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

func Parse() error {

	qorConf := configr.New()
	qorConf.RegisterKey("port", "QOR admin port", 8000)
	qorConf.RegisterKey("sitename", "QOR admin site name", "QOR Admin")
	qorConf.AddSource(configr.NewFile("qor.toml"))
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

	Hugo.MetaDataFormat = "toml"

	hugoConf := configr.New()
	hugoConf.RegisterKey("baseurl", "Hugo site baseurl", "/")
	hugoConf.RegisterKey("staticdir", "Hugo site static dir", "static")
	hugoConf.RegisterKey("publishdir", "Hugo site publish dir", "public")
	hugoConf.RegisterKey("languageCode", "Hugo site languageCode", "en")
	hugoConf.RegisterKey("disableRSS", "Hugo site disableRSS", true)
	hugoConf.RegisterKey("menu", "Hugo site menus", make(map[string]interface{}))

	hugoConf.AddSource(configr.NewFile("hugo.toml"))
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

	return nil

}
