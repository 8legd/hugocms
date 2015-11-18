package admin

import (
	"fmt"
	"os"
	//"reflect"

	"github.com/adrianduke/configr"

	"github.com/jinzhu/gorm"

	"github.com/qor/qor"
	"github.com/qor/qor/admin"
	//	"github.com/qor/qor/resource"
	//	"github.com/qor/qor/utils"
	//	"github.com/qor/qor/validations"

	"github.com/8legd/hugocms/admin/models"
	"github.com/8legd/hugocms/db"
	//	"strings"
)

var (
	Admin *admin.Admin

	users *admin.Resource

	paths []string

	pages *admin.Resource

	videos *admin.Resource

	strapline *admin.Resource
	settings  *admin.Resource
)

func init() {

	paths = append(paths, "/")

	configr.AddSource(configr.NewFileSource("config/admin.toml"))
	if err := configr.Parse(); err != nil {
		fmt.Println(err)
		os.Exit(1)
		//TODO more graceful exit!
	}

	cfg, err := configr.Get("paths")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
		//TODO more graceful exit!
	}

	if arr, ok := cfg.([]interface{}); ok {
		for _, v := range arr {
			if s, ok := v.(string); ok {
				paths = append(paths, s)
			}
		}
	}

	Admin = admin.New(&qor.Config{DB: db.DB})

	Admin.SetAuth(Auth{})

	// For QOR Admin's rich text editor
	assetManager := Admin.AddResource(&admin.AssetManager{}, &admin.Config{Invisible: true})

	users = Admin.AddResource(&models.User{}, &admin.Config{Name: "Users"})
	users.IndexAttrs("ID", "Name")
	users.EditAttrs("Name", "Image")

	columns := Admin.NewResource(&models.PageContentColumn{}, &admin.Config{Invisible: true})
	columns.Meta(&admin.Meta{Name: "TextContent", Type: "rich_editor", Resource: assetManager})
	columns.Meta(&admin.Meta{Name: "VideoOptions", Type: "select_one", Collection: []string{"Auto Play"}})
	columns.NewAttrs("-ContentRow")
	columns.EditAttrs("-ContentRow", "Heading", "Video")

	rows := Admin.NewResource(&models.PageContentRow{}, &admin.Config{Invisible: true})
	rows.Meta(&admin.Meta{Name: "ContentColumns", Resource: columns})
	rows.NewAttrs("-Page")
	rows.EditAttrs("-Page")

	// TODO validation, example below
	//page.AddValidator(func(record interface{}, metaValues *resource.MetaValues, context *qor.Context) error {
	//	if meta := metaValues.Get("Name"); meta != nil {
	//		if name := utils.ToString(meta.Value); strings.TrimSpace(name) == "" {
	//			return validations.NewError(record, "Name", "Name can't be blank")
	//		}
	//	}
	//	return nil
	//})

	pages = Admin.AddResource(&models.Page{}, &admin.Config{Name: "Pages"})
	pages.IndexAttrs("Path", "Name")
	pages.Meta(&admin.Meta{Name: "ContentRows", Resource: rows})
	pages.Meta(&admin.Meta{Name: "Path", Type: "select_one", Collection: paths})

	// define scopes for pages
	for _, path := range paths {
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

	videos = Admin.AddResource(&models.Video{}, &admin.Config{Name: "Videos"})
	videos.IndexAttrs("Name")

	strapline := Admin.NewResource(&models.SettingsStrapline{}, &admin.Config{Invisible: true})
	strapline.Meta(&admin.Meta{Name: "TextContent", Type: "rich_editor", Resource: assetManager})
	//strapline.NewAttrs("-Settings")
	//strapline.EditAttrs("-Settings")

	settings = Admin.AddResource(&models.Settings{}, &admin.Config{Singleton: true})
	settings.Meta(&admin.Meta{Name: "Strapline", Resource: strapline})
	settings.Meta(&admin.Meta{Name: "Footer", Type: "rich_editor", Resource: assetManager})

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
