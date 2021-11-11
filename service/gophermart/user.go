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
func (g *GopherMart) Create(ctx context.Context, login, password string) (*model.User, error) {
	const logPrefix = "service: create:"
	log := appContext.Logger(ctx)

	id, err := uuid.NewRandom()
	if err != nil {
		log.Trace().Err(err).Msg(logPrefix)
		return nil, fmt.Errorf("%s %w", logPrefix, err)
	}
	user := &model.User{
		ID:             id,
		Login:          login,
		Password:       password,
		CreatedAt:      time.Now(),
		GPointsBalance: 0,
	}

	if err := user.Validate(); err != nil {
		log.Trace().Err(err).Msg(logPrefix)
		return nil, err
	}

	user.PasswordHash, err = bcrypt.BcryptPassword(password, g.pwPepper)
	if err != nil {
		log.Trace().Err(err).Msg(logPrefix)
		return nil, fmt.Errorf("%s %w", logPrefix, err)
	}
	user.Password = ""

	if err := g.db.NewUser(ctx, user); err != nil {
		log.Trace().Err(err).Msg(logPrefix)
		return nil, fmt.Errorf("%s %w", logPrefix, err)
	}
	log.Info().
		Str("login", user.Login).
		Str("id", user.ID.String()).
		Msg(logPrefix + " successfully created a new user")

	return user, nil
}

// Authenticate implements Service interface.
func (g *GopherMart) Authenticate(ctx context.Context, login, password string) (*model.User, error) {
	log := appContext.Logger(ctx)
	const logPrefix = "service: authenticate:"

	// we don't need to validate login & password - in the DB all is OK.
	user, err := g.db.UserByLogin(ctx, login)
	if err != nil {
		log.Trace().Err(err).Msg(logPrefix)
		return nil, fmt.Errorf("service: authenticate: %w", err)
	}

	if err = bcrypt.CompareHashAndPassword(password+g.pwPepper, user.PasswordHash); err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			log.Trace().Err(err).Msg(logPrefix)
			return nil, ErrWrongPassword
		}
		log.Trace().Err(err).Msg(logPrefix)
		return nil, fmt.Errorf("service: authenticate: %w", err)
	}

	log.Info().
		Str("login", user.Login).
		Str("id", user.ID.String()).
		Msg(logPrefix + " successfully authenticated the user")

	return user, nil
}

// GetOrders implements Service interface.
func (g *GopherMart) GetOrders(ctx context.Context) ([]model.Order, error) {
	log := userLogger(ctx)
	const logPrefix = "service: GetOrders:"

	user := appContext.User(ctx)
	if user == nil {
		log.Trace().Err(ErrNotAuthenticated).Msg(logPrefix)
		return nil, ErrNotAuthenticated
	}

	orders, err := g.db.UserOrders(ctx, user.ID)
	if err != nil {
		log.Trace().Err(err).Msg(logPrefix)
		return nil, fmt.Errorf("service: GetOrders: %w", err)
	}

	return orders, nil
}

// GetWithdrawals implements Service interface.
func (g *GopherMart) GetWithdrawals(ctx context.Context) ([]model.Withdrawal, error) {
	log := userLogger(ctx)
	const logPrefix = "service: GetWithdrawals:"

	user := appContext.User(ctx)
	if user == nil {
		log.Trace().Err(ErrNotAuthenticated).Msg(logPrefix)
		return nil, ErrNotAuthenticated
	}

	withdrawals, err := g.db.WithdrawalsByUserID(ctx, user.ID)
	if err != nil {
		log.Trace().Err(err).Msg(logPrefix)
		return nil, fmt.Errorf("service: GetWithdrawals: %w", err)
	}

	return withdrawals, nil
}

// GetBalance implements Service interface.
func (g *GopherMart) GetBalance(ctx context.Context) (*UserBalance, error) {
	log := userLogger(ctx)
	const logPrefix = "service: GetBalance:"

	user := appContext.User(ctx)
	if user == nil {
		log.Trace().Err(ErrNotAuthenticated).Msg(logPrefix)
		return nil, ErrNotAuthenticated
	}

	cb := &UserBalance{}

	withdrawals, err := g.db.WithdrawalsByUserID(ctx, user.ID)
	if err != nil && !errors.Is(err, storage.ErrNotFound) { // It's OK if the user has not performed any withdraw operations.
		log.Trace().Err(err).Msg(logPrefix)
		return nil, fmt.Errorf("service: GetWithdrawals: %w", err)
	}

	// Collect information about total withdrawn bonus amount.
	for _, w := range withdrawals {
		if w.Status == model.StatusProcessed {
			cb.Withdrawn += w.Sum
		}
	}

	// Update balance information
	if _, err := g.db.UpdateBalance(ctx); err != nil {
		log.Trace().Err(err).Msg(logPrefix)
		return nil, err
	}

	// Update current user balance information
	user, err = g.db.UserByLogin(ctx, user.Login)
	if err != nil {
		log.Trace().Err(err).Msg(logPrefix)
		return nil, err
	}
	cb.Current = user.GPointsBalance

	log.Info().
		Float32("current", cb.Current).
		Float32("withdrawn", cb.Withdrawn).
		Msg("information about user's balance successfully received")

	return cb, nil
}
