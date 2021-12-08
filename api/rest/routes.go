package rest

import (
	"github.com/go-chi/chi"
	"github.com/rs/zerolog"
	"github.com/vanamelnik/gophermart/api/handlers"
	"github.com/vanamelnik/gophermart/pkg/middleware"
	"github.com/vanamelnik/gophermart/service/gophermart"
	"github.com/vanamelnik/gophermart/storage"
)

// SetupRoutes configures mux.
func SetupRoutes(service gophermart.Service, db storage.Storage, log zerolog.Logger) *chi.Mux {
	h := handlers.New(service, db)

	// Setup routes
	r := chi.NewRouter()
	r.Use(middleware.WithLogger(log))
	r.Use(middleware.GzipMdlw)

	r.Route("/api/user", func(r chi.Router) {
		r.Post("/register", h.Register)
		r.Post("/login", h.Login)

		r.Route("/", func(r chi.Router) {
			r.Use(middleware.UserCtx(db))

			r.Post("/orders", h.PostOrder)
			r.Get("/orders", h.GetOrders)
			r.Get("/balance", h.GetBalance)
			r.Post("/balance/withdraw", h.Withdraw)
			r.Get("/balance/withdrawals", h.GetWithdrawals)
		})
	})

	return r
}
