package gophermart

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/vanamelnik/gophermart/model"
	appContext "github.com/vanamelnik/gophermart/pkg/ctx"
	"github.com/vanamelnik/gophermart/storage"
)

// ProcessOrder implements Service interface.
func (g *GopherMart) ProcessOrder(ctx context.Context, orderID model.OrderID) error {
	log := userLogger(ctx).With().Str("service:", "ProcessOrder").Logger()

	user := appContext.User(ctx)
	if user == nil {
		log.Trace().Err(ErrNotAuthenticated).Msg("internal:")
		return ErrNotAuthenticated
	}
	log = log.With().Str("order ID", string(orderID)).Logger()

	// Check if the order is already stored in DB.
	o, err := g.db.OrderByID(ctx, orderID)
	switch {
	case err == nil:
		if o.UserID == user.ID {
			log.Trace().Err(ErrOrderExecutedBySameUser).Msg("")
			return ErrOrderExecutedBySameUser
		}
		log.Trace().
			Str("order owner", o.UserID.String()).
			Err(ErrOrderExecutedByAnotherUser).
			Msg("")
		return ErrOrderExecutedByAnotherUser

	case !errors.Is(err, storage.ErrNotFound):
		log.Trace().Err(err).Msg("")
		return fmt.Errorf("service: ProcessOrder: %w", err)
	}

	order := &model.Order{
		ID:            orderID,
		UserID:        user.ID,
		Status:        model.StatusNew,
		AccrualPoints: 0,
		UploadedAt:    time.Now(),
	}
	if err := g.db.CreateOrder(ctx, order); err != nil {
		// storage.ErrAlreadyProcessed is an internal server error because the presence
		// of an already loaded order in the database should have been determined by OrderById method.
		log.Trace().Err(err).Msg("")
		return fmt.Errorf("service: ProcessOrder: %w", err)
	}
	log.Trace().Msg("the order has been successfully stored in DB")

	return nil
}

// Withdraw implements Service interface.
func (g *GopherMart) Withdraw(ctx context.Context, orderID model.OrderID, sum float32) error {
	log := userLogger(ctx).With().Str("service:", "withdraw").Logger()

	user := appContext.User(ctx)
	if user == nil {
		log.Trace().Err(ErrNotAuthenticated).Msg("")
		return ErrNotAuthenticated
	}

	if err := g.db.ProcessWithdraw(ctx, &model.Withdrawal{
		UserID:      user.ID,
		OrderID:     orderID,
		Sum:         sum,
		Status:      model.StatusProcessing,
		ProcessedAt: time.Now(),
	}); err != nil {
		log.Trace().Err(err).Msg("")
		return fmt.Errorf("service: withdraw: %w", err)
	}

	return nil
}
