package gophermart

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/vanamelnik/go-musthave-diploma/model"
	appContext "github.com/vanamelnik/go-musthave-diploma/pkg/ctx"
	"github.com/vanamelnik/go-musthave-diploma/storage"
)

// ProcessOrder implements Service interface.
func (g *GopherMart) ProcessOrder(ctx context.Context, orderID model.OrderID) error {
	log := userLogger(ctx)
	const logMsg = "service: ProcessOrder:"

	user := appContext.User(ctx)
	if user == nil {
		log.Trace().Err(ErrNotAuthenticated).Msg("internal:")
		return ErrNotAuthenticated
	}
	// if !orderID.Valid() {
	// 	log.Trace().Err(ErrInvalidOrderNumber).Msg(logMsg)
	// 	return ErrInvalidOrderNumber
	// }
	log = log.With().Str("order ID", string(orderID)).Logger()

	// Check if the order is already stored in DB.
	o, err := g.db.OrderByID(ctx, orderID)
	switch {
	case err == nil:
		if o.UserID == user.ID {
			log.Trace().Err(ErrOrderExecutedBySameUser).Msg(logMsg)
			return ErrOrderExecutedBySameUser
		}
		log.Trace().
			Str("order owner", o.UserID.String()).
			Err(ErrOrderExecutedByAnotherUser).
			Msg(logMsg)
		return ErrOrderExecutedByAnotherUser

	case !errors.Is(err, storage.ErrNotFound):
		log.Trace().Err(err).Msg(logMsg)
		return fmt.Errorf("%s %w", logMsg, err)
	}

	order := &model.Order{
		ID:            orderID,
		UserID:        user.ID,
		Status:        model.StatusNew,
		AccrualPoints: 0,
		UploadedAt:    time.Now(),
	}
	if err := g.db.NewOrder(ctx, order); err != nil {
		// storage.ErrAlreadyProcessed is an internal server error because the presence
		// of an already loaded order in the database should have been determined by OrderById method.
		log.Trace().Err(err).Msg(logMsg)
		return fmt.Errorf("%s %w", logMsg, err)
	}
	log.Trace().Msgf("%s the order has been successfully stored in DB", logMsg)

	return nil
}

// Withdraw implements Service interface.
func (g *GopherMart) Withdraw(ctx context.Context, orderID model.OrderID, sum float32) error {
	log := userLogger(ctx)
	const logMsg = "service: withdraw:"

	user := appContext.User(ctx)
	if user == nil {
		log.Trace().Err(ErrNotAuthenticated).Msg(logMsg)
		return ErrNotAuthenticated
	}

	if err := g.db.NewWithdraw(ctx, &model.Withdrawal{
		UserID:      user.ID,
		OrderID:     orderID,
		Sum:         sum,
		Status:      model.StatusProcessing,
		ProcessedAt: time.Now(),
	}); err != nil {
		log.Trace().Err(err).Msg(logMsg)
		return fmt.Errorf("%s %w", logMsg, err)
	}

	log.Info().
		Float32("withdrawed", sum).
		Msg("successfully withdrawed")

	return nil
}
