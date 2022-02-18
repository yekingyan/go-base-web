package sharetrace

import (
	"fmt"
	"time"
)

// XRequestIDKey is the request id key.
const XRequestIDKey = "X-Request-ID"

// NewRequestID create a request id.
func NewRequestID() string {
	return fmt.Sprintf("req-id-%d", time.Now().UnixNano())
}
