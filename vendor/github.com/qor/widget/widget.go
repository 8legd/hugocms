package widget

import (
	"fmt"
	"os"
	"path"

	"github.com/jinzhu/gorm"
	"github.com/qor/admin"
	"github.com/qor/qor/resource"
	"github.com/qor/roles"
)

var (
	root, _           = os.Getwd()
	viewPaths         []string
	registeredWidgets []*Widget
)

// Config widget config
type Config struct {
	DB    *gorm.DB
	Admin *admin.Admin
}

func init() {
	if path := os.Getenv("WEB_ROOT"); path != "" {
		root = path
	}
	registerViewPath(path.Join(root, "app/views/widgets"))
}

// New new widgets container
func New(config *Config) *Widgets {
	instance := &Widgets{Config: config}
	instance.RegisterViewPath("app/views/widgets")
	return instance
}

// Widgets widgets container
type Widgets struct {
	Config                *Config
	Resource              *admin.Resource
	WidgetSettingResource *admin.Resource
}

// ConfigureQorResource a method used to config Widget for qor admin
func (widgets *Widgets) ConfigureQorResource(res resource.Resourcer) {
	if res, ok := res.(*admin.Resource); ok {
		// register view paths
		admin.RegisterViewPath("github.com/qor/widget/views")

		// set resources
		res.Name = "Widget"
		widgets.Resource = res

		// set setting resource
		if widgets.WidgetSettingResource == nil {
			widgets.WidgetSettingResource = res.GetAdmin().NewResource(&QorWidgetSetting{}, &admin.Config{Permission: roles.Deny(roles.Create, roles.Anyone)})
			widgets.WidgetSettingResource.IndexAttrs("ID", "Kind", "Name")
			widgets.WidgetSettingResource.Name = res.Name
		}

		// configure routes
		controller := widgetController{Widgets: widgets}
		router := res.GetAdmin().GetRouter()
		router.Get(res.ToParam(), controller.Index)
		router.Get(fmt.Sprintf("%v/frontend-edit", res.ToParam()), controller.FronendEdit)
		router.Get(fmt.Sprintf("%v/%v", res.ToParam(), res.ParamIDName()), controller.Edit)
		router.Get(fmt.Sprintf("%v/%v/edit", res.ToParam(), res.ParamIDName()), controller.Edit)
		router.Put(fmt.Sprintf("%v/%v", res.ToParam(), res.ParamIDName()), controller.Update)
	}
}

// RegisterWidget register a new widget
func (widgets *Widgets) RegisterWidget(w *Widget) {
	registeredWidgets = append(registeredWidgets, w)
}

// Widget widget struct
type Widget struct {
	Name     string
	Template string
	Setting  *admin.Resource
	Context  func(context *Context, setting interface{}) *Context
}

// GetWidget get widget by name
func GetWidget(name string) *Widget {
	for _, w := range registeredWidgets {
		if w.Name == name {
			return w
		}
	}
	return nil
}
