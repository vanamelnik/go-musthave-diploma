package middleware

import (
	"net/http"

	appContext "github.com/vanamelnik/go-musthave-diploma/pkg/ctx"

	"github.com/rs/zerolog"
)

// WithLogger returns a middleware func that attaches the logger provided to client's request context.
// This middleware should be the first in the middlewres chain.
func WithLogger(logger zerolog.Logger) mwFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := appContext.WithLogger(r.Context(), logger)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
