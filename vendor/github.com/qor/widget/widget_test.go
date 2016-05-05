package widget_test

import (
	"fmt"
	"testing"

	"net/http"
	"strings"

	"github.com/fatih/color"
	"github.com/jinzhu/gorm"
	"github.com/qor/admin"
	"github.com/qor/qor"
	"github.com/qor/qor/test/utils"
	"github.com/qor/widget"
)

var db *gorm.DB

func init() {
	db = utils.TestDB()
}

// Runner
func TestRender(t *testing.T) {
	db.DropTable(&widget.QorWidgetSetting{})
	db.AutoMigrate(&widget.QorWidgetSetting{})

	Widgets := widget.New(&widget.Config{
		DB: db,
	})
	Widgets.RegisterViewPath("github.com/qor/widget/test")

	Admin := admin.New(&qor.Config{DB: db})
	Admin.AddResource(Widgets)
	Admin.MountTo("/admin", http.NewServeMux())

	type bannerArgument struct {
		Title    string
		SubTitle string
	}

	Widgets.RegisterWidget(&widget.Widget{
		Name:     "Banner",
		Template: "banner",
		Setting:  Admin.NewResource(&bannerArgument{}),
		Context: func(context *widget.Context, setting interface{}) *widget.Context {
			if setting != nil {
				argument := setting.(*bannerArgument)
				context.Options["Title"] = argument.Title
				context.Options["SubTitle"] = argument.SubTitle
			}
			return context
		},
	})

	html := Widgets.Render("HomeBanner", nil, "Banner")
	if !strings.Contains(string(html), "Hello, \n<h1></h1>\n<h2></h2>\n") {
		t.Errorf(color.RedString(fmt.Sprintf("\nWidget Render TestCase #%d: Failure Result:\n %s\n", 1, html)))
	}

	widgetContext := widget.NewContext(map[string]interface{}{
		"CurrentUser": "Qortex",
	})
	html = Widgets.Render("HomeBanner", widgetContext, "Banner")
	if !strings.Contains(string(html), "Hello, Qortex\n<h1></h1>\n<h2></h2>\n") {
		t.Errorf(color.RedString(fmt.Sprintf("\nWidget Render TestCase #%d: Failure Result:\n %s\n", 2, html)))
	}

	setting := widget.QorWidgetSetting{}
	db.Where("name = ? AND kind = ?", "HomeBanner", "Banner").First(&setting)
	setting.SetSerializableArgumentValue(&bannerArgument{Title: "Title", SubTitle: "SubTitle"})
	db.Save(&setting)

	html = Widgets.Render("HomeBanner", widgetContext, "Banner")
	if !strings.Contains(string(html), "Hello, Qortex\n<h1>Title</h1>\n<h2>SubTitle</h2>\n") {
		t.Errorf(color.RedString(fmt.Sprintf("\nWidget Render TestCase #%d: Failure Result:\n %s\n", 3, html)))
	}
}
