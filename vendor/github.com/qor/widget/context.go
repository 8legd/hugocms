package widget

import "errors"

// NewContext new widget context
func NewContext(options map[string]interface{}, availableWidgets ...string) *Context {
	return &Context{
		Widgets: availableWidgets,
		Options: options,
	}
}

// Context widget context
type Context struct {
	Widgets []string
	Options map[string]interface{}
}

// Get get option with name
func (context Context) Get(name string) (interface{}, error) {
	if value, ok := context.Options[name]; ok {
		return value, nil
	}

	return nil, errors.New("not found")
}

// Set set option by name
func (context *Context) Set(name string, value interface{}) {
	if context.Options == nil {
		context.Options = map[string]interface{}{}
	}
	context.Options[name] = value
}
