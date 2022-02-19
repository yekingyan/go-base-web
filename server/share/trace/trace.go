package sharetrace

import "github.com/google/uuid"

// XRequestIDKey is the request id key.
const XRequestIDKey = "X-Request-ID"

// NewRequestID create a request id.
func NewRequestID() string {
	return uuid.New().String()
}
