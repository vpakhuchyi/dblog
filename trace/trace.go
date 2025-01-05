package trace

import (
	"context"

	"github.com/google/uuid"
)

// CorrelationID is the type used to store the correlation ID in the context.
type CorrelationID struct{}

// contextKey is the key used to store the correlation ID in the context.
var contextKey = CorrelationID{}

// GenerateCorrelationID generates a new UUID as correlation ID.
func GenerateCorrelationID() string {
	return uuid.NewString()
}

// GetCorrelationID extracts the correlation ID from the context.
func GetCorrelationID(ctx context.Context) string {
	if v := ctx.Value(contextKey); v != nil {
		if id, ok := v.(string); ok {
			return id
		}
	}

	return ""
}

// SetCorrelationID adds the correlation ID to the context.
func SetCorrelationID(ctx context.Context, id string) context.Context {
	return context.WithValue(ctx, contextKey, id)
}
