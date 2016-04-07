package qor

import (
	"strings"

	"github.com/jinzhu/gorm"
	"github.com/qor/admin"
	"github.com/qor/media_library"
	"github.com/qor/qor"
	"github.com/qor/qor/resource"
	"github.com/qor/qor/utils"
	"github.com/qor/sorting"
	"github.com/qor/validations"

	"github.com/8legd/hugocms/config"
	"github.com/8legd/hugocms/qor/models"
)

var (
	Tables []interface{}

	settings   *admin.Resource
	videos     *admin.Resource
	slideshows *admin.Resource
	pages      *admin.Resource
)

func init() {

	// Define database tables used by CMS
	Tables = []interface{}{

		&models.Settings{},
		&models.SettingsLogo{},
		&models.SettingsContactDetails{},
		&models.SettingsHeader{},
		&models.SettingsCallToAction{},
		&models.SettingsIntroVideo{},
		&models.SettingsIntroVideoMeta{},
		&models.SettingsSidebarContent{},

		&models.Video{},

		&models.Slideshow{},
		&models.SlideshowSlide{},

		&models.Page{},
		&models.PageMeta{},
		&models.PageContentColumn{},
		&models.PageContentColumnVideo{},
		&models.PageContentColumnSlideshow{},
		&models.PageLink{},
	}

}

func SetupAdmin() *admin.Admin {

	// Setup Database for QOR Admin
	sorting.RegisterCallbacks(config.DB)
	validations.RegisterCallbacks(config.DB)
	media_library.RegisterCallbacks(config.DB)

	result := admin.New(&qor.Config{DB: config.DB})

	result.SetSiteName(config.QOR.SiteName)
	result.SetAuth(config.Auth)

	// Add Asset Manager, for rich editor
	assetManager := result.AddResource(&media_library.AssetManager{}, &admin.Config{Invisible: true})

	image := result.NewResource(&models.PageContentColumnImage{}, &admin.Config{Invisible: true})
	image.Meta(&admin.Meta{
		Name: "Alignment",
		Type: "select_one",
		Collection: func(o interface{}, context *qor.Context) [][]string {
			var result [][]string
			result = append(result, []string{"media-left media-top", "left top"})
			result = append(result, []string{"media-left media-middle", "left middle"})
			result = append(result, []string{"media-left media-bottom", "left bottom"})
			result = append(result, []string{"media-right media-top", "right top"})
			result = append(result, []string{"media-right media-middle", "right middle"})
			result = append(result, []string{"media-right media-bottom", "right bottom"})
			return result
		},
	})
	image.NewAttrs("-PageContentColumn")
	image.EditAttrs("-PageContentColumn")

	columns := result.NewResource(&models.PageContentColumn{}, &admin.Config{Invisible: true})
	columns.Meta(&admin.Meta{
		Name: "Width",
		Type: "select_one",
		Collection: func(o interface{}, context *qor.Context) [][]string {
			var result [][]string
			result = append(result, []string{"col-md-6", "50% on desktop, 100% on mobile"})
			result = append(result, []string{"col-md-12", "100% on desktop, 100% on mobile"})
			return result
		},
	})
	columns.Meta(&admin.Meta{Name: "TextContent", Type: "rich_editor", Resource: assetManager})
	columns.Meta(&admin.Meta{Name: "ImageContent", Resource: image})
	columns.NewAttrs("-Page")
	columns.EditAttrs("-Page")

	links := result.NewResource(&models.PageLink{}, &admin.Config{Invisible: true})
	links.Meta(&admin.Meta{Name: "LinkText", Type: "rich_editor", Resource: assetManager})
	links.NewAttrs("-Page")
	links.EditAttrs("-Page")

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

	pages.Meta(&admin.Meta{Name: "ContentColumns", Resource: columns})

	pages.Meta(&admin.Meta{Name: "Links", Resource: links})

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

		// TODO make SEO required

		// if we have content check it is valid
		if meta := metaValues.Get("ContentRows"); meta != nil {
			if metas := meta.MetaValues.Values; len(metas) > 0 {
				for _, v := range metas {
					if v.Name == "ContentColumns" {
						// first collect information on the fields we need to check
						img := ""
						imgPos := ""

						if fields := v.MetaValues.Values; len(fields) > 0 {
							for _, f := range fields {

								if f.Name == "Image" && f.Value != nil {
									for _, s := range f.Value.([]string) {
										if s != "" {
											img = s
										}
									}
								}
								if f.Name == "Alignment" && f.Value != nil {
									for _, s := range f.Value.([]string) {
										if s != "" {
											imgPos = s
										}
									}
								}
							}
						}

						// all images need an image position

						if img != "" && imgPos == "" {
							return validations.NewError(record, "ContentRows", "All Images need an Image Position")
						}
					}
				}
			}

		}
		return nil
	})

	videos = result.AddResource(&models.Video{}, &admin.Config{Name: "Videos"})
	videos.IndexAttrs("Name")

	slideshows = result.AddResource(&models.Slideshow{}, &admin.Config{Name: "Slideshow"})
	slideshows.IndexAttrs("Name")

	contact := result.NewResource(&models.SettingsContactDetails{}, &admin.Config{Invisible: true})
	contact.Meta(&admin.Meta{Name: "OpeningHoursDesktop", Type: "rich_editor", Resource: assetManager})

	callToAction := result.NewResource(&models.SettingsCallToAction{}, &admin.Config{Invisible: true})
	callToAction.Meta(&admin.Meta{Name: "ActionText", Type: "rich_editor", Resource: assetManager})

	settings = result.AddResource(&models.Settings{}, &admin.Config{Singleton: true})
	settings.Meta(&admin.Meta{Name: "ContactDetails", Resource: contact})
	settings.Meta(&admin.Meta{Name: "CallToAction", Resource: callToAction})
	settings.Meta(&admin.Meta{Name: "Footer", Type: "rich_editor", Resource: assetManager})

	return result
}
