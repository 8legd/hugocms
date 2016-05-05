package widget

import (
	"bytes"
	"errors"
	"fmt"
	"html/template"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/qor/qor/utils"
)

// Render find widget by name, render it based on current context
func (widgets *Widgets) Render(widgetName string, context *Context, availableWidgets ...string) template.HTML {
	if context == nil {
		context = NewContext(map[string]interface{}{})
	}

	if len(availableWidgets) == 0 {
		availableWidgets = append(availableWidgets, widgetName)
	}

	var (
		setting      = findSettingByNameAndKinds(widgets.Config.DB, widgetName, availableWidgets)
		widgetObj    = GetWidget(setting.Kind)
		settingValue = setting.GetSerializableArgument(setting)
		newContext   = widgetObj.Context(context, settingValue)
		url          = widgets.settingEditURL(setting)
		prefix       = widgets.Resource.GetAdmin().GetRouter().Prefix
	)

	return template.HTML(fmt.Sprintf(
		"<script data-prefix=\"%v\" src=\"%v/assets/javascripts/widget_check.js?theme=widget\"></script><div class=\"qor-widget qor-widget-%v\" data-widget-frontend-edit-url=\"%v\" data-url=\"%v\">\n%v\n</div>",
		prefix,
		prefix,
		utils.ToParamString(widgetObj.Name),
		fmt.Sprintf("%v/%v/frontend-edit", prefix, widgets.Resource.ToParam()),
		url,
		widgetObj.Render(newContext, url),
	))
}

func (widgets *Widgets) settingEditURL(setting *QorWidgetSetting) string {
	prefix := widgets.WidgetSettingResource.GetAdmin().GetRouter().Prefix
	return fmt.Sprintf("%v/%v/%v/edit", prefix, widgets.WidgetSettingResource.ToParam(), setting.ID)
}

// FuncMap return view functions map
func (widgets *Widgets) FuncMap() template.FuncMap {
	funcMap := template.FuncMap{}

	funcMap["render_widget"] = func(key string, context *Context, availableWidgets ...string) template.HTML {
		return widgets.Render(key, context, availableWidgets...)
	}

	return funcMap
}

// Render register widget itself content
func (w *Widget) Render(context *Context, url string) template.HTML {
	var err error
	var result = bytes.NewBufferString("")
	file := w.Template

	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("Get error when render file %v: %v", file, r)
			utils.ExitWithMsg(err)
		}
	}()

	if file, err = w.findTemplate(file + ".tmpl"); err == nil {
		if tmpl, err := template.New(filepath.Base(file)).ParseFiles(file); err == nil {
			if err = tmpl.Execute(result, context.Options); err == nil {
				return template.HTML(result.String())
			}
		}
	}

	return template.HTML(err.Error())
}

// RegisterViewPath register views directory
func (widgets *Widgets) RegisterViewPath(p string) {
	for _, gopath := range strings.Split(os.Getenv("GOPATH"), ":") {
		if registerViewPath(path.Join(gopath, "src", p)) == nil {
			return
		}
	}
}

func isExistingDir(pth string) bool {
	fi, err := os.Stat(pth)
	if err != nil {
		return false
	}
	return fi.Mode().IsDir()
}

func registerViewPath(path string) error {
	if isExistingDir(path) {
		var found bool

		for _, viewPath := range viewPaths {
			if path == viewPath {
				found = true
				break
			}
		}

		if !found {
			viewPaths = append(viewPaths, path)
		}
		return nil
	}
	return errors.New("path not found")
}

func (w *Widget) findTemplate(layouts ...string) (string, error) {
	for _, layout := range layouts {
		for _, p := range viewPaths {
			if _, err := os.Stat(filepath.Join(p, layout)); !os.IsNotExist(err) {
				return filepath.Join(p, layout), nil
			}
		}
	}
	return "", fmt.Errorf("template not found: %v", layouts)
}
