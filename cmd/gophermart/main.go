package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"

	"github.com/vanamelnik/go-musthave-diploma/api/handlers"
	"github.com/vanamelnik/go-musthave-diploma/cmd/gophermart/config"
	appContext "github.com/vanamelnik/go-musthave-diploma/pkg/ctx"
	"github.com/vanamelnik/go-musthave-diploma/pkg/logging"
	"github.com/vanamelnik/go-musthave-diploma/pkg/middleware"
	"github.com/vanamelnik/go-musthave-diploma/provider/accrual"
	"github.com/vanamelnik/go-musthave-diploma/service/gophermart"
	"github.com/vanamelnik/go-musthave-diploma/storage/psql"

	"github.com/go-chi/chi"
)

func main() {
	// Load config.
	cfg := config.LoadConfig("./config.toml")
	must(cfg.Validate())

	// Create the logger.
	log := logging.NewLogger(logging.WithConsoleOutput(cfg.Logger.Console), logging.WithLevel(cfg.Logger.Level))
	ctx := appContext.WithLogger(context.Background(), log)
	log.Trace().Msgf("config loaded: %+v", cfg)

	// Connect to the database.
	db, err := psql.New(
		psql.WithConfig(cfg.Database),
		psql.WithAutoMigrate(log, "file://storage/psql/migration"),
	)
	must(err)
	defer db.Close()

	// Start Gophermart Service
	service, err := gophermart.New(
		ctx, db,
		gophermart.WithConfig(cfg.Service),
		gophermart.WithAccrualClient(accrual.New(cfg.AccrualSystemAddr)),
	)
	must(err)
	defer service.Close()

	// Setup handlers.
	h := handlers.New(service, db)

	// Setup routes
	router := chi.NewRouter()
	router.Use(middleware.WithLogger(log))
	router.Use(middleware.GzipMdlw)

	router.Post("/api/user/register", h.Register)
	router.Post("/api/user/login", h.Login)

	router.With(middleware.RequireUser(db)).Post("/api/user/orders", h.PostOrder)
	router.With(middleware.RequireUser(db)).Get("/api/user/orders", h.GetOrders)
	router.With(middleware.RequireUser(db)).Get("/api/user/balance", h.GetBalance)
	router.With(middleware.RequireUser(db)).Post("/api/user/balance/withdraw", h.Withdraw)
	router.With(middleware.RequireUser(db)).Get("/api/user/balance/withdrawals", h.GetWithdrawals)

	server := http.Server{
		Addr:    cfg.RunAddr,
		Handler: router,
	}
	sigint := make(chan os.Signal, 1)
	signal.Notify(sigint, os.Interrupt)

	go func() {
		err := server.ListenAndServe()
		log.Info().Err(err).Msg("server stopped")
	}()
	log.Info().Msgf("main: the server is listening at %s", cfg.RunAddr)

	<-sigint
	log.Info().Msg("main: shutting down... ")
	if err := server.Shutdown(ctx); err != nil {
		log.Error().Err(err).Msg("server shutdown")
	}
}

func must(err error) {
	if err != nil {
		panic(err)
	}
}
