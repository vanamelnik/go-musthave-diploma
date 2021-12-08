package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"

	"github.com/vanamelnik/gophermart/api/rest"
	"github.com/vanamelnik/gophermart/cmd/gophermart/config"
	appContext "github.com/vanamelnik/gophermart/pkg/ctx"
	"github.com/vanamelnik/gophermart/pkg/logging"
	"github.com/vanamelnik/gophermart/provider/accrual"
	"github.com/vanamelnik/gophermart/service/gophermart"
	"github.com/vanamelnik/gophermart/storage/psql"
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
		if err != nil && err != http.ErrServerClosed {
			log.Error().Err(err).Msg("server stopped")
			return
		}
		log.Info().Msg("server stopped")
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
