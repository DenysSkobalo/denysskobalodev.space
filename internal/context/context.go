package context

type ContextKey string

const (
	AdminHashKey ContextKey = "ADMIN_HASH"
	SaltKey      ContextKey = "SALT"
)
