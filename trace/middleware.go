package trace

import (
	"net/http"
)

// CorrelationIDHeader is the header key for the correlation ID.
const CorrelationIDHeader = "X-Correlation-ID"

// WithCorrelationID is a middleware to add a correlation ID to the request context.
func WithCorrelationID(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		correlationID := r.Header.Get(CorrelationIDHeader)
		if correlationID == "" {
			correlationID = GenerateCorrelationID()
		}

		// Add the correlation ID to the response headers.
		w.Header().Set(CorrelationIDHeader, correlationID)

		// Add the correlation ID to the context.
		r = r.WithContext(SetCorrelationID(r.Context(), correlationID))

		// Call the next handler in the chain.
		next.ServeHTTP(w, r)
	})
}
