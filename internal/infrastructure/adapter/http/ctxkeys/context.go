package ctxkeys

// ContextKey is the type for context keys to avoid collisions.
type ContextKey string

// RequestIDKey is the context key for the request ID.
const RequestIDKey ContextKey = "request_id"
