package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"

	"github.com/vanamelnik/go-musthave-diploma/api/rest"
	"github.com/vanamelnik/go-musthave-diploma/cmd/gophermart/config"
	appContext "github.com/vanamelnik/go-musthave-diploma/pkg/ctx"
	"github.com/vanamelnik/go-musthave-diploma/pkg/logging"
	"github.com/vanamelnik/go-musthave-diploma/provider/accrual"
	"github.com/vanamelnik/go-musthave-diploma/service/gophermart"
	"github.com/vanamelnik/go-musthave-diploma/storage/psql"
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
		psql.WithDSN(cfg.DatabaseURI),
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

	// Setup routes
	router := rest.SetupRoutes(service, db, log)
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
