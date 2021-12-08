package gophermart

import (
	"context"
	"sync"
	"time"

	appContext "github.com/vanamelnik/gophermart/pkg/ctx"
	"github.com/vanamelnik/gophermart/pkg/logging"
	"github.com/vanamelnik/gophermart/provider/accrual"
	"github.com/vanamelnik/gophermart/storage"

	"github.com/rs/zerolog"
)

const (
	defaultAccrualURL = "/"
)

// Ensure service implements interface.
var _ Service = (*GopherMart)(nil)

// GopherMart is an implementation of gophermart.Service interface.
type (
	GopherMart struct {
		// withWorkers is used for testing. If false, accrualServicePoller and balanceUpdater
		// doesn't run.
		withWorkers bool
		workersStop chan struct{}
		workersWg   sync.WaitGroup
		db          storage.Storage
		pwPepper    string
		// accrualClient calls for sending a request to accrual service.
		accrualClient      accrual.AccrualClient
		balanceUpdInterval time.Duration
	}

	Config struct {
		PasswordPepper string        `mapstructure:"password_pepper"`
		UpdateInterval time.Duration `mapstructure:"update_interval"`
	}

	ServiceOption func(*GopherMart)
)

// WithPepper overrides default pepper for users' passwords.
func WithConfig(cfg Config) ServiceOption {
	return func(g *GopherMart) {
		g.pwPepper = cfg.PasswordPepper
		g.balanceUpdInterval = cfg.UpdateInterval
	}
}

// WithAccrualClient sets the provider for GopherAccrualService.
func WithAccrualClient(client accrual.AccrualClient) ServiceOption {
	return func(g *GopherMart) {
		g.accrualClient = client
	}
}

// WithoutWorkers used for testing. It turns off accrualServicePoller and balanceUpdater workers.
func WithoutWorkers() ServiceOption {
	return func(g *GopherMart) {
		g.withWorkers = false
	}
}

// New creates a new Gophermart object with provided database and other custom options.
// Contract: expected logger in context.
func New(ctx context.Context, db storage.Storage, opts ...ServiceOption) (*GopherMart, error) {
	g := &GopherMart{
		workersStop:   make(chan struct{}),
		db:            db,
		accrualClient: accrual.New(defaultAccrualURL),
		withWorkers:   true,
	}
	for _, opt := range opts {
		opt(g)
	}
	if g.balanceUpdInterval == 0 {
		g.withWorkers = false // do not start workers if update interval isn't set.
	}

	if g.withWorkers {
		// Start AccrualService poller and balance updater.
		g.workersWg.Add(2)
		go g.accrualServicePoller(ctx)
		go g.balanceUpdater(ctx)
	}

	return g, nil
}

func (g *GopherMart) Close() {
	if g.workersStop != nil {
		close(g.workersStop)
	}

	g.workersWg.Wait() // wait until workers stopped
	if g.db != nil {
		g.db.Close()
		g.db = nil
	}
}

// userLogger gets a loger from context and expands it by user data fields.
func userLogger(ctx context.Context) zerolog.Logger {
	log := appContext.Logger(ctx)
	c := appContext.User(ctx)
	if c != nil {
		// We use user's login to improve readability as log logs into the console.
		log = log.With().Str(logging.UserKey, c.Login).Logger()
	}

	return log
}
