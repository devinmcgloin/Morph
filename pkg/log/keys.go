package log

type ContextKey string

const (
	IDKey     = ContextKey("fokal-id-key")
	IPKey     = ContextKey("fokal-ip-key")
	JWTID     = ContextKey("fokal-jwt-key")
	UserIDKey = ContextKey("fokal-user-key")
)
