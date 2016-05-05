package widget

import (
	"net/http"

	"github.com/qor/admin"
	"github.com/qor/responder"
)

type widgetController struct {
	Widgets *Widgets
}

func (wc widgetController) Index(context *admin.Context) {
	context = context.NewResourceContext(wc.Widgets.WidgetSettingResource)
	result, err := context.FindMany()
	context.AddError(err)

	if context.HasError() {
		http.NotFound(context.Writer, context.Request)
	} else {
		responder.With("html", func() {
			context.Execute("index", result)
		}).With("json", func() {
			context.JSON("index", result)
		}).Respond(context.Request)
	}
}

func (wc widgetController) Edit(context *admin.Context) {
	qorSetting := &QorWidgetSetting{}
	context.Resource = wc.Widgets.WidgetSettingResource
	err := wc.Widgets.WidgetSettingResource.FindOneHandler(qorSetting, nil, context.Context)
	context.AddError(err)
	context.Execute("edit", qorSetting)
}

func (wc widgetController) Update(context *admin.Context) {
	qorSetting := &QorWidgetSetting{}
	context.Resource = wc.Widgets.WidgetSettingResource
	err := wc.Widgets.WidgetSettingResource.FindOneHandler(qorSetting, nil, context.Context)
	context.AddError(err)
	if context.AddError(context.Resource.Decode(context.Context, qorSetting)); !context.HasError() {
		context.AddError(context.Resource.CallSave(qorSetting, context.Context))
		context.Execute("edit", qorSetting)
		return
	}

	http.Redirect(context.Writer, context.Request, context.Request.URL.Path, http.StatusFound)
}

func (wc widgetController) FronendEdit(context *admin.Context) {
	context.Writer.Write([]byte(context.Render("front_edit")))
}
