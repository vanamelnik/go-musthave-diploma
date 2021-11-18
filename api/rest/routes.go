package rest

import (
	"github.com/go-chi/chi"
	"github.com/rs/zerolog"
	"github.com/vanamelnik/go-musthave-diploma/api/handlers"
	"github.com/vanamelnik/go-musthave-diploma/pkg/middleware"
	"github.com/vanamelnik/go-musthave-diploma/service/gophermart"
	"github.com/vanamelnik/go-musthave-diploma/storage"
)

// SetupRoutes configures mux.
func SetupRoutes(service gophermart.Service, db storage.Storage, log zerolog.Logger) *chi.Mux {
	h := handlers.New(service, db)

	// Setup routes
	r := chi.NewRouter()
	r.Use(middleware.WithLogger(log))
	r.Use(middleware.GzipMdlw)

	r.Post("/api/user/register", h.Register)
	r.Post("/api/user/login", h.Login)

	// // These endpoints required the user authenticated:
	// r.Route("/api/user/orders", func(r chi.Router) {
	// 	r.Use(middleware.RequireUser(db))
	// 	r.Post("/", h.PostOrder)
	// 	r.Get("/", h.GetOrders)
	// })
	// r.Route("/api/user/balance", func(r chi.Router) {
	// 	r.Use(middleware.RequireUser(db))
	// 	r.Get("/", h.GetBalance)
	// 	r.Post("/withdraw", h.Withdraw)
	// 	r.Get("/withdrawals", h.GetWithdrawals)
	// })
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
