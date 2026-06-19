package http

import (
	"net/http"

	"github.com/forcexdd/portfoliomanager/internal/logger"
	"github.com/google/uuid"
)

func ContextLoggerMiddleware(baseLogger logger.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			reqID := r.Header.Get("X-Request-ID")
			if reqID == "" {
				reqID = uuid.New().String()
			}

			enrichedLogger := baseLogger.With("req_id", reqID, "method", r.Method, "path", r.URL.Path)
			ctx := logger.IntoContext(r.Context(), enrichedLogger)

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
