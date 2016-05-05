package widget

var registeredScopes []*Scope

// Scope widget scope
type Scope struct {
	Name    string
	Visible func(*Context) bool
}

// RegisterScope register scope for widget
func RegisterScope(scope *Scope) {
	registeredScopes = append(registeredScopes, scope)
}
