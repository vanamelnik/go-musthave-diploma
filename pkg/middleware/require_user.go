package middleware

import (
	"net/http"

	appContext "github.com/vanamelnik/go-musthave-diploma/pkg/ctx"
	"github.com/vanamelnik/go-musthave-diploma/storage"
)

// RequreUser returns a middleware function that checks if there's a user's remember token
// in client's cookies. If OK and the user found in the storage, its object is attached
// to the requst context.
func RequireUser(db storage.Storage) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			log := appContext.Logger(r.Context())
			cookie, err := r.Cookie("gophermart_remember")
			if err != nil {
				log.Error().Err(err).Msg("RequireUser: cookie with remember token not found")
				http.Error(w, "Unauthorized", http.StatusUnauthorized)

				return
			}
			remember := cookie.Value
			user, err := db.UserByRemember(r.Context(), remember)
			if err != nil {
				log.Error().Err(err).Msgf("RequireUser: user with remember token %s not found", remember)
				http.Error(w, "Unauthorized", http.StatusUnauthorized)

				return
			}

			ctx := appContext.WithUser(r.Context(), user)
			log.Info().Str("user", user.Login).Msg("RequireUser: successfully authorized")
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
