package reqctx

// AdminContext doesn't store any significant information but help to distinguish Admin from User endpoints
type AdminContext struct {
	isAdmin bool // should be always be true for all AdminHandler(s)
}

// NewAdminContext creates a new AdminContext
func NewAdminContext() *AdminContext {
	return &AdminContext{
		isAdmin: true,
	}
}
