package gophermart

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/vanamelnik/go-musthave-diploma/model"
	"github.com/vanamelnik/go-musthave-diploma/pkg/bcrypt"
	appContext "github.com/vanamelnik/go-musthave-diploma/pkg/ctx"
	"github.com/vanamelnik/go-musthave-diploma/storage"

	"github.com/google/uuid"
)

// Create implements Service interface.
func (g *GopherMart) Create(ctx context.Context, login, password string) (model.User, error) {
	log := appContext.Logger(ctx).With().Str("service:", "create").Logger()

	id, err := uuid.NewRandom()
	if err != nil {
		log.Trace().Err(err).Msg("")
		return model.User{}, fmt.Errorf("service: create: %w", err)
	}
	user := model.User{
		ID:             id,
		Login:          login,
		Password:       password,
		CreatedAt:      time.Now(),
		GPointsBalance: 0,
	}

	if err := user.Validate(); err != nil {
		log.Trace().Err(err).Msg("")
		return model.User{}, err
	}

	user.PasswordHash, err = bcrypt.BcryptPassword(password, g.pwPepper)
	if err != nil {
		log.Trace().Err(err).Msg("")
		return model.User{}, fmt.Errorf("service: create: %w", err)
	}
	user.Password = ""

	if err := g.db.CreateUser(ctx, user); err != nil {
		log.Trace().Err(err).Msg("")
		return model.User{}, fmt.Errorf("service: create: %w", err)
	}
	log.Info().
		Str("login", user.Login).
		Str("id", user.ID.String()).
		Msg("successfully created a new user")

	return user, nil
}

// Authenticate implements Service interface.
func (g *GopherMart) Authenticate(ctx context.Context, login, password string) (model.User, error) {
	log := appContext.Logger(ctx).With().Str("service:", "authenticate:").Logger()

	// we don't need to validate login & password - in the DB all is OK.
	user, err := g.db.UserByLogin(ctx, login)
	if err != nil {
		log.Trace().Err(err).Msg("")
		return model.User{}, fmt.Errorf("service: authenticate: %w", err)
	}

	if err = bcrypt.CompareHashAndPassword(password+g.pwPepper, user.PasswordHash); err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			log.Trace().Err(err).Msg("")
			return model.User{}, ErrWrongPassword
		}
		log.Trace().Err(err).Msg("")
		return model.User{}, fmt.Errorf("service: authenticate: %w", err)
	}

	log.Info().
		Str("login", user.Login).
		Str("id", user.ID.String()).
		Msg("successfully authenticated the user")

	return *user, nil
}

// GetOrders implements Service interface.
func (g *GopherMart) GetOrders(ctx context.Context) ([]model.Order, error) {
	log := userLogger(ctx).With().Str("service:", "getOrders").Logger()

	user := appContext.User(ctx)
	if user == nil {
		log.Trace().Err(ErrNotAuthenticated).Msg("")
		return nil, ErrNotAuthenticated
	}

	orders, err := g.db.UserOrders(ctx, user.ID)
	if err != nil {
		log.Trace().Err(err).Msg("")
		return nil, fmt.Errorf("service: GetOrders: %w", err)
	}

	return orders, nil
}

// GetWithdrawals implements Service interface.
func (g *GopherMart) GetWithdrawals(ctx context.Context) ([]model.Withdrawal, error) {
	log := userLogger(ctx).With().Str("service:", "getWithdrawals:").Logger()

	user := appContext.User(ctx)
	if user == nil {
		log.Trace().Err(ErrNotAuthenticated).Msg("")
		return nil, ErrNotAuthenticated
	}

	withdrawals, err := g.db.WithdrawalsByUserID(ctx, user.ID)
	if err != nil {
		log.Trace().Err(err).Msg("")
		return nil, fmt.Errorf("service: GetWithdrawals: %w", err)
	}

	return withdrawals, nil
}

// GetBalance implements Service interface.
func (g *GopherMart) GetBalance(ctx context.Context) (UserBalance, error) {
	log := userLogger(ctx).With().Str("service:", "GetBalance:").Logger()

	user := appContext.User(ctx)
	if user == nil {
		log.Trace().Err(ErrNotAuthenticated).Msg("")
		return UserBalance{}, ErrNotAuthenticated
	}

	cb := UserBalance{}

	withdrawals, err := g.db.WithdrawalsByUserID(ctx, user.ID)
	if err != nil && !errors.Is(err, storage.ErrNotFound) { // It's OK if the user has not performed any withdraw operations.
		log.Trace().Err(err).Msg("")
		return UserBalance{}, fmt.Errorf("service: GetBalance: %w", err)
	}

	// Collect information about total withdrawn bonus amount.
	for _, w := range withdrawals {
		if w.Status == model.StatusProcessed {
			cb.Withdrawn += w.Sum
		}
	}

	// Update balance information
	if _, err := g.db.UpdateBalance(ctx); err != nil {
		log.Trace().Err(err).Msg("")
		return UserBalance{}, err
	}

	// Update current user balance information
	user, err = g.db.UserByLogin(ctx, user.Login)
	if err != nil {
		log.Trace().Err(err).Msg("")
		return UserBalance{}, err
	}
	cb.Current = user.GPointsBalance

	log.Info().
		Float32("current", cb.Current).
		Float32("withdrawn", cb.Withdrawn).
		Msg("information about user's balance successfully received")

	return cb, nil
}
