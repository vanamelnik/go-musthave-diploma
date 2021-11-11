package psql

import (
	"context"
	"testing"
	"time"

	"github.com/vanamelnik/go-musthave-diploma/model"
	"github.com/vanamelnik/go-musthave-diploma/pkg/bcrypt"
	"github.com/vanamelnik/go-musthave-diploma/pkg/logging"
	"github.com/vanamelnik/go-musthave-diploma/storage"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"

	"github.com/google/uuid"
	"github.com/stretchr/testify/suite"
)

const (
	migrationsPath = "file://migration"
)

type (
	TestSuite struct {
		suite.Suite
		storage storage.Storage
		ctx     context.Context

		bob   fixture
		alice fixture
	}

	fixture struct {
		user   *model.User
		orders []*model.Order
	}
)

func (ts *TestSuite) SetupSuite() {
	storage, err := New(WithConfig(defaultConfig),
		WithAutoMigrate(logging.NewLogger(logging.WithConsoleOutput(true)),
			migrationsPath))
	ts.Require().NoError(err)
	ts.storage = storage
	ts.ctx = context.TODO()

	ts.bob = fixture{
		user: &model.User{
			Login: "bobmarley@rambler.ru",
		},
		orders: []*model.Order{
			{
				ID:     "018",
				Status: "PROCESSING",
			},
			{
				ID:     "026",
				Status: "REGISTERED",
			},
		},
	}
	ts.alice = fixture{
		user: &model.User{
			Login: "alicecooper@yandex.cn",
		},
		orders: []*model.Order{
			{
				ID:     "117",
				Status: "PROCESSING",
			},
			{
				ID:     "125",
				Status: "PROCESSED",
			},
		},
	}
	ts.loadFixtures()
}

func (ts *TestSuite) loadFixtures() {
	for _, f := range []fixture{ts.alice, ts.bob} {
		// Setup and load users
		hash, err := bcrypt.BcryptPassword(f.user.Login, "")
		ts.Require().NoError(err)
		f.user.ID = uuid.New()
		f.user.PasswordHash = hash
		f.user.CreatedAt = time.Now()
		err = ts.storage.NewUser(ts.ctx, f.user)
		ts.Require().NoError(err)

		// Setup and load orders
		for _, o := range f.orders {
			o.UserID = f.user.ID
			o.UploadedAt = time.Now()
			err := ts.storage.NewOrder(ts.ctx, o)
			ts.Require().NoError(err)
		}
	}
	ts.T().Log("Fixtures for Bob and Alice successfully created.")
}

func TestPsql(t *testing.T) {
	suite.Run(t, new(TestSuite))
}

func (ts *TestSuite) TearDownSuite() {
	ts.Require().NoError(ts.storage.Close())

	// Migrate down. TODO: add test container.
	m, err := migrate.New(migrationsPath, defaultDSN)
	ts.Require().NoError(err)
	defer m.Close()

	ts.Require().NoError(m.Down())
}
