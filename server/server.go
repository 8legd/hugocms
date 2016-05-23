package server

import (
	"fmt"
	"net/http"
	"os"
	
	"github.com/adrianduke/configr"
	_ "github.com/adrianduke/configr/sources/file/toml"
	"github.com/astaxie/beego/session"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	_ "github.com/mattn/go-sqlite3"
	"github.com/qor/admin"
	"github.com/qor/i18n"
	"github.com/qor/i18n/backends/database"
	"github.com/qor/qor"

	"github.com/8legd/hugocms/config"
	hugocms_qor "github.com/8legd/hugocms/qor"
	"github.com/8legd/hugocms/qor/models"
)

var SessionManager *session.Manager

type Auth struct {
	UserName string
	Password string
}

func (a Auth) LoginURL(c *admin.Context) string {
	return "/login"
}

func (a Auth) LogoutURL(c *admin.Context) string {
	return "/admin/logout"
}

func (a Auth) GetCurrentUser(c *admin.Context) qor.CurrentUser {
	w := c.Writer
	r := c.Request
	sess, err := SessionManager.SessionStart(w, r)
	if err != nil {
		handleError(err)
	}
	defer sess.SessionRelease(w)

	if r.URL.String() == "/admin/auth" &&
		r.FormValue("inputAccount") != "" &&
		(r.FormValue("inputAccount") == a.UserName) &&
		r.FormValue("inputPassword") != "" &&
		(r.FormValue("inputPassword") == a.Password) {
		sess.Set("User", User{a.UserName})
	}
	if u, ok := sess.Get("User").(User); ok && u.Name != "" {
		return u
	}
	return nil
}

type User struct {
	Name string
}

func (u User) DisplayName() string {
	return u.Name
}

type DatabaseType int

const (
	_ DatabaseType = iota
	DB_SQLite
	DB_MySQL
)

func ListenAndServe(addr string, auth Auth, dbType DatabaseType) {
	var db *gorm.DB
	var err error

	if dbType == DB_MySQL {
		dbConn := fmt.Sprintf("%s:%s@tcp(127.0.0.1:3306)", auth.UserName, auth.Password)
		db, err = gorm.Open("mysql", dbConn+"/hugocms_"+auth.UserName+"?charset=utf8&parseTime=True&loc=Local")
	} else {
		db, err = gorm.Open("sqlite3", "hugocms_"+auth.UserName+".db")
	}

	if err != nil {
		handleError(err)
	}
	db.LogMode(true)

	// to reset to an empty database drop the settings table
	if !db.HasTable(&models.Settings{}) {
		for _, table := range hugocms_qor.Tables {
			if err := db.DropTableIfExists(table).Error; err != nil {
				handleError(err)
			}
		}
	}

	for _, table := range hugocms_qor.Tables {
		if err := db.AutoMigrate(table).Error; err != nil {
			handleError(err)
		}
	}

	// Because this is a singleton a single record must already exist in the database (there is no create through QOR admin just update)
	db.FirstOrCreate(&models.Settings{})

	siteName := fmt.Sprintf("%s - Hugo CMS", auth.UserName)
	if err := setupConfig(addr, siteName, db, auth); err != nil {
		handleError(err)
	}

	// Add session support - used by Auth
	sessionLifetime := 3600 // session lifetime in seconds
	SessionManager, err = session.NewManager("memory", fmt.Sprintf(`{"cookieName":"gosessionid","gclifetime":%d}`, sessionLifetime))
	if err != nil {
		handleError(err)
	}
	go SessionManager.GC()

	// Create Hugo's content directory if it doesnt exist
	// TODO read content dir from config
	if _, err := os.Stat("./content"); os.IsNotExist(err) {
		err = os.MkdirAll("./content", os.ModePerm)
	}

	mux := http.NewServeMux()

	mux.Handle("/", http.FileServer(http.Dir("public")))

	adm := hugocms_qor.SetupAdmin()

	adm.MountTo("/admin", mux)
	adm.GetRouter().Post("/auth", func(ctx *admin.Context) {
		// we will only hit this on succesful login - redirect to admin dashboard
		w := ctx.Writer
		r := ctx.Request
		http.Redirect(w, r, "/admin", http.StatusFound)
	})
	adm.GetRouter().Get("/logout", func(ctx *admin.Context) {
		w := ctx.Writer
		r := ctx.Request
		sess, err := SessionManager.SessionStart(w, r)
		if err != nil {
			handleError(err)
		}
		defer sess.SessionRelease(w)
		sess.Delete("User")
		http.Redirect(w, r, "/login", http.StatusFound)
	})

	// NOTE: `system` is where QOR admin will upload files e.g. images - we map this to Hugo's static dir along with our other static assets
	// TODO read static dir from config
	// TODO read static assets list from config
	for _, path := range []string{"system", "css", "fonts", "images", "js", "login"} {
		mux.Handle(fmt.Sprintf("/%s/", path), http.FileServer(http.Dir("static")))
	}

	if err := http.ListenAndServe(config.QOR.Addr, mux); err != nil {
		handleError(err)
	}

	// to re-generate site delete `config.json`
	if _, err := os.Stat("config.json"); os.IsNotExist(err) {
		hugocms_qor.CallSave(adm)
	}

	fmt.Printf("Listening on: %s\n", config.QOR.Addr)
}

func setupConfig(addr string, sitename string, db *gorm.DB, auth admin.Auth) error {

	config.QOR.Addr = addr
	config.QOR.SiteName = sitename

	// As a minumum add the root path for our site
	config.QOR.Paths = append(config.QOR.Paths, "/")

	config.Hugo.MetaDataFormat = "json"

	hugoConf := configr.New()
	hugoConf.RegisterKey("baseurl", "Hugo site baseurl", "/")
	hugoConf.RegisterKey("staticdir", "Hugo site static dir", "static")
	hugoConf.RegisterKey("publishdir", "Hugo site publish dir", "public")
	hugoConf.RegisterKey("languageCode", "Hugo site languageCode", "en")
	hugoConf.RegisterKey("disableRSS", "Hugo site disableRSS", true)
	hugoConf.RegisterKey("menu", "Hugo site menus", make(map[string]interface{}))

	hugoConfigFile := "hugo.toml"
	hugoConf.AddSource(configr.NewFile(hugoConfigFile))
	if err := hugoConf.Parse(); err != nil {
		return err
	}

	baseurl, err := hugoConf.String("baseurl")
	if err != nil {
		return err
	}
	config.Hugo.BaseURL = baseurl

	staticdir, err := hugoConf.String("staticdir")
	if err != nil {
		return err
	}
	config.Hugo.StaticDir = staticdir

	publishdir, err := hugoConf.String("publishdir")
	if err != nil {
		return err
	}
	config.Hugo.PublishDir = publishdir

	languageCode, err := hugoConf.String("languageCode")
	if err != nil {
		return err
	}
	config.Hugo.LanguageCode = languageCode

	disableRSS, err := hugoConf.Bool("disableRSS")
	if err != nil {
		return err
	}
	config.Hugo.DisableRSS = disableRSS

	rawMenu, err := hugoConf.Get("menu")
	if err != nil {
		return err
	}
	if menu, ok := rawMenu.(map[string]interface{}); ok {
		config.Hugo.Menu = menu
		// Add additional site paths from main menu items
		if rawMainMenu, ok := menu["main"]; ok {
			if mainMenu, ok := rawMainMenu.([]map[string]interface{}); ok {
				for _, item := range mainMenu {
					if url, ok := item["url"].(string); ok {
						if url != "" && url != "/" {
							config.QOR.Paths = append(config.QOR.Paths, url)
						}
					}
				}
			}
		}
	}

	config.DB = db
	config.I18n = i18n.New(database.New(db))

	config.Auth = auth

	return nil

}

func handleError(err error) {
	fmt.Println(err)
	os.Exit(1)
	//TODO more graceful exit!
}
