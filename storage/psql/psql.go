package psql

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/golang-migrate/migrate/v4"
	"github.com/rs/zerolog"
	"github.com/vanamelnik/go-musthave-diploma/storage"

	// Register postgres and file drivers.
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"

	// Register pgx driver.
	_ "github.com/jackc/pgx/stdlib"
)

// Ensure service implements interface.
var _ storage.Storage = (*Psql)(nil)

type (
	Psql struct {
		config Config
		db     *sql.DB
	}

	Option func(*Psql) error
)

// WithConfig overrides default config.
func WithConfig(config Config) Option {
	return func(p *Psql) error {
		p.config = config

		return nil
	}
}

// WithAutoMigrate applies migrate against db. Should run after setting up correct DSN string in the config.
func WithAutoMigrate(log zerolog.Logger, path string) Option {
	return func(p *Psql) error {
		m, err := migrate.New(path, p.config.DSN)
		if err != nil {
			return fmt.Errorf("psql: WithAutoMigrate: %w", err)
		}
		defer m.Close()

		if err := m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
			return fmt.Errorf("psql: WithAutoMigrate: %w", err)
		}

		log.Info().Str("service", "psql").Msg("auto migrate applied")
		return nil
	}
}

// New creates a new connection to postgres database.
func New(opts ...Option) (*Psql, error) {
	p := &Psql{config: defaultConfig}
	for i, opt := range opts {
		if err := opt(p); err != nil {
			return nil, fmt.Errorf("storage: applying option [%d]: %w", i, err)
		}
	}
	if err := p.config.Validate(); err != nil {
		return nil, err
	}

	db, err := sql.Open("pgx", p.config.DSN)
	if err != nil {
		return nil, err
	}
	p.db = db

	if err := p.db.Ping(); err != nil {
		return nil, fmt.Errorf("storage: ping for DSN (%s) failed: %w", p.config.DSN, err)
	}

	return p, nil
}

func (p Psql) Close() error {
	if p.db == nil {
		return nil
	}

	return p.db.Close()
}
