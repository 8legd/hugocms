package qor

import (
	"strings"

	//"github.com/astaxie/beego/session" TODO handle auth / session management
	"github.com/jinzhu/gorm"
	"github.com/qor/media_library"
	"github.com/qor/qor"
	"github.com/qor/qor/admin"
	"github.com/qor/qor/resource"
	"github.com/qor/qor/utils"
	"github.com/qor/qor/validations"
	"github.com/qor/sorting"

	"github.com/8legd/hugocms/config"
	"github.com/8legd/hugocms/qor/models"
)

var (
	Tables []interface{}

	users    *admin.Resource
	pages    *admin.Resource
	videos   *admin.Resource
	header   *admin.Resource
	settings *admin.Resource
)

func init() {

	// Define database tables used by CMS
	Tables = []interface{}{

		&models.User{},
		&models.UserImage{},

		&models.Settings{},
		&models.SettingsLogo{},
		&models.SettingsContactDetails{},
		&models.SettingsHeader{},
		&models.SettingsSidebar{},
		&models.SettingsIntroVideo{},
		&models.SettingsIntroVideoMeta{},
		&models.SettingsSidebarImage{},

		&models.Video{},

		&models.Page{},
		&models.PageMeta{},
		&models.PageContentRow{},
		&models.PageContentColumn{},
		&models.PageSlideshowImage{},
	}

}

func SetupAdmin() *admin.Admin {

	// Setup Database for QOR Admin
	sorting.RegisterCallbacks(config.QOR.DB)
	validations.RegisterCallbacks(config.QOR.DB)
	media_library.RegisterCallbacks(config.QOR.DB)

	result := admin.New(&qor.Config{DB: config.QOR.DB})

	result.SetSiteName(config.QOR.SiteName)
	result.SetAuth(Auth{})

	// TODO Add Dashboard
	// result.AddMenu(&admin.Menu{Name: "Dashboard", Link: "/admin"})

	// Add Asset Manager, for rich editor
	assetManager := result.AddResource(&media_library.AssetManager{}, &admin.Config{Invisible: true})

	users = result.AddResource(&models.User{}, &admin.Config{Name: "Users"})
	users.IndexAttrs("ID", "Name")
	users.EditAttrs("Name", "Image")

	columns := result.NewResource(&models.PageContentColumn{}, &admin.Config{Invisible: true})
	columns.Meta(&admin.Meta{Name: "TextContent", Type: "rich_editor", Resource: assetManager})
	columns.Meta(&admin.Meta{Name: "VideoOptions", Type: "select_one", Collection: []string{"Auto Play"}})
	columns.NewAttrs("-ContentRow")
	columns.EditAttrs("-ContentRow", "Heading", "Video")

	rows := result.NewResource(&models.PageContentRow{}, &admin.Config{Invisible: true})
	rows.Meta(&admin.Meta{Name: "ContentColumns", Resource: columns})
	rows.NewAttrs("-Page")
	rows.EditAttrs("-Page")

	pages = result.AddResource(&models.Page{}, &admin.Config{Name: "Pages"})
	pages.IndexAttrs("Path", "Name")
	pages.Meta(&admin.Meta{
		Name: "MenuWeight",
		Type: "select_one",
		Collection: func(o interface{}, context *qor.Context) [][]string {
			// Build menu weight drop down on the fly...
			var result [][]string
			// Check we have a path (if not set menu weight to 0)
			if p, ok := o.(*models.Page); ok && p.Path != "" {
				// TODO find out the current max menu weight for this path
				//var pages []models.Page
				//db.DB.Find(&pages)
				result = append(result, []string{"0", "0"})
				result = append(result, []string{"1", "1"})
			} else {
				result = append(result, []string{"0", "0"})
			}
			return result

		},
	})
	pages.Meta(&admin.Meta{Name: "ContentRows", Resource: rows})

	pages.Meta(&admin.Meta{Name: "Path", Type: "select_one", Collection: config.QOR.Paths})

	// define scopes for pages
	for _, path := range config.QOR.Paths {
		path := path // The anonymous function below captures the variable `path` not its value
		// So because the range variable is re-assigned a value on each iteration, if we just used it,
		// the actual value being used would just end up being the same (last value of iteration).
		// By redeclaring `path` within the range block's scope a new variable is in effect created for each iteration
		// and that specific variable is used in the anonymous function instead
		// Another solution would be to pass the range variable into a function as a parameter which then returns the
		// original function you wanted creating a `closure` around the passed in parameter (you often  come accross this in JavaScript)
		pages.Scope(&admin.Scope{
			Name:  path,
			Group: "Path",
			Handle: func(db *gorm.DB, context *qor.Context) *gorm.DB {
				return db.Where(models.Page{Path: path})
			},
		})
	}
	pages.AddValidator(func(record interface{}, metaValues *resource.MetaValues, context *qor.Context) error {
		if meta := metaValues.Get("Name"); meta != nil {
			if name := utils.ToString(meta.Value); strings.TrimSpace(name) == "" {
				return validations.NewError(record, "Name", "Name can not be blank")
			}
		}
		if meta := metaValues.Get("Path"); meta != nil {
			if name := utils.ToString(meta.Value); strings.TrimSpace(name) == "" {
				return validations.NewError(record, "Path", "Path can not be blank")
			}
		}
		return nil
	})

	videos = result.AddResource(&models.Video{}, &admin.Config{Name: "Videos"})
	videos.IndexAttrs("Name")

	header := result.NewResource(&models.SettingsHeader{}, &admin.Config{Invisible: true})
	header.Meta(&admin.Meta{Name: "TextContent", Type: "rich_editor", Resource: assetManager})

	settings = result.AddResource(&models.Settings{}, &admin.Config{Singleton: true})
	settings.Meta(&admin.Meta{Name: "Header", Resource: header})
	settings.Meta(&admin.Meta{Name: "Footer", Type: "rich_editor", Resource: assetManager})

	return result
}

type Auth struct{}

func (Auth) LoginURL(c *admin.Context) string {
	return "/admin"
}

func (Auth) LogoutURL(c *admin.Context) string {
	return "/admin"
}

func (Auth) GetCurrentUser(c *admin.Context) qor.CurrentUser {
	return &models.User{Name: "Admin"}
}
