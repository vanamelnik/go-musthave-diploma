package gophermart

import (
	"context"
	"errors"

	"github.com/vanamelnik/go-musthave-diploma/model"
	appContext "github.com/vanamelnik/go-musthave-diploma/pkg/ctx"
	"github.com/vanamelnik/go-musthave-diploma/provider/accrual"
	"github.com/vanamelnik/go-musthave-diploma/storage"
)

// accrualServicePoller looks in the storage for 'NEW' and 'PROCESSING' orders and sends requests to the
// GopherAccrualService. If poller receives a response with a status changed, the order status changedes in DB.
// If accrual calculation is done, a new entry in accruals log creates.
func (g *GopherMart) accrualServicePoller(ctx context.Context) {
	log := appContext.Logger(ctx)
	log.Info().Msg("accrualServicePoller started")
poller:
	for {
		select {
		case <-ctx.Done():
			break poller
		default:
			orders := g.getOrders(ctx)
			for _, order := range orders {
				g.processOrder(ctx, order)
			}
		}
	}
	log.Info().Msg("accrualServicePoller stopped")
}

// getOrders returns the orders with statuses 'NEW', 'REGISTERED', 'PROCESSING' from the storage.
func (g *GopherMart) getOrders(ctx context.Context) []model.Order {
	log := appContext.Logger(ctx)
	const logPrefix = "service: poller: getOrders:"
	orders := make([]model.Order, 0)
	for _, status := range []model.Status{model.StatusNew, model.StatusProcessing, model.StatusRegistered} {
		o, err := g.db.OrdersByStatus(ctx, status)
		if err != nil { // if there're no orders, empty list is returned and err == nil.
			log.Error().Err(err).Msgf("%s could not get orders with status %s", logPrefix, status)
		}
		orders = append(orders, o...)
	}

	return orders
}

// processOrder sends a request to GopherAccrualService and updates order status to 'REGISTERED', 'PROCESSING'
// or 'INVALID'. If the calculation is done, the new entry in accruals log is created.
func (g *GopherMart) processOrder(ctx context.Context, order model.Order) {
	log := appContext.Logger(ctx).With().Str("orderID", order.ID.String()).Logger()
	const logPrefix = "service: poller: process order:"
	// Mark NEW orders as REGISTERED.
	if order.Status == model.StatusNew {
		if err := g.db.UpdateOrderStatus(ctx, order.ID, model.StatusRegistered); err != nil {
			log.Error().Err(err).Msgf("%s could not update order status, operation cancelled", logPrefix)

			return
		}
		order.Status = model.StatusRegistered
	}
	// Send a request to the GopherAccualService
	resp, err := g.accrualClient.Request(ctx, order.ID)
	if err != nil {
		var apiErr *accrual.ErrUnexpectedStatus
		if errors.As(err, &apiErr) {
			log.Warn().Err(err).Msgf("%s accrual service response:", logPrefix)

			return
		}
		log.Error().Err(err).Msgf("%s accrual service response:", logPrefix)

		return
	}
	// If the accruals are calculated, create a new entry in accruals log
	if resp.Status == model.StatusProcessed {
		if resp.Accrual < 0 {
			if err := g.db.UpdateOrderStatus(ctx, order.ID, model.StatusInvalid); err != nil {
				log.Error().Err(err).Msgf("%s could not update order status, operation cancelled", logPrefix)

				return
			}
		}
		err := g.db.NewAccrual(ctx, order.ID, resp.Accrual) // if all is OK, the order status is set to 'PROCESSED' within db transaction.
		switch {
		case err == nil:
			log.Info().Float32("amount", resp.Accrual).Msgf("%s a new entry in accruals log has been created", logPrefix)

			return
		case errors.Is(err, storage.ErrAlreadyProcessed):
			log.Error().Err(err).Msgf("%s internal error! processed order has status %s", logPrefix, order.Status)

			return
		case errors.Is(err, storage.ErrNotFound):
			log.Error().Err(err).Msgf("%s internal error! user not found", logPrefix)

			return
		default:
			log.Error().Err(err).Msgf("%s internal error:", logPrefix)

			return
		}

	}
	// Update order status. Accrual calculation hasn't processed yet.
	if resp.Status != order.Status {
		if err := g.db.UpdateOrderStatus(ctx, order.ID, resp.Status); err != nil {
			log.Error().Err(err).Msgf("%s could not update order status, operation cancelled", logPrefix)

			return
		}
		log.Info().Str("status", string(resp.Status)).Msgf("%s status updated", logPrefix)
	}
}
