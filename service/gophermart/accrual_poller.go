package gophermart

import (
	"context"
	"errors"

	"github.com/vanamelnik/go-musthave-diploma/model"
	appContext "github.com/vanamelnik/go-musthave-diploma/pkg/ctx"
	"github.com/vanamelnik/go-musthave-diploma/provider/accrual"
	"github.com/vanamelnik/go-musthave-diploma/storage"
)

// accrualServicePoller looks in the storage for 'NEW', and 'PROCESSING' orders and sends requests to the
// GopherAccrualService. If poller receives a response with a status changed, the order status changedes in DB.
// If accrual calculation is done, a new entry in accruals log creates.
func (g *GopherMart) accrualServicePoller(ctx context.Context) {
	log := appContext.Logger(ctx)
	log.Info().Msg("accrualServicePoller started")
poller:
	for {
		select {
		case <-g.workersStop:
			break poller
		default:
			orders := g.getOrders(ctx)
			for _, order := range orders {
				g.processOrder(ctx, order)
			}
		}
	}
	log.Info().Msg("accrualServicePoller stopped")
	g.workersWg.Done()
}

// getOrders returns the orders with statuses 'NEW' and 'PROCESSING' from the storage.
func (g *GopherMart) getOrders(ctx context.Context) []model.Order {
	log := appContext.Logger(ctx).With().Str("service:", "poller: getOrders:").Logger()
	orders := make([]model.Order, 0)
	for _, status := range []model.Status{model.StatusNew, model.StatusProcessing} {
		o, err := g.db.OrdersByStatus(ctx, status)
		if err != nil { // if there're no orders, empty list is returned and err == nil.
			log.Error().Err(err).Msgf("could not get orders with status %s", status)
		}
		orders = append(orders, o...)
	}

	return orders
}

// processOrder sends a request to GopherAccrualService and updates order status to 'PROCESSING'
// or 'INVALID'. If the calculation is done, the new entry in accruals log is created.
func (g *GopherMart) processOrder(ctx context.Context, order model.Order) {
	log := appContext.Logger(ctx).With().
		Str("orderID", order.ID.String()).
		Str("service:", "poller: process order:").
		Logger()

	// Send a request to the GopherAccualService
	resp, err := g.accrualClient.Request(ctx, order.ID)
	if err != nil {
		var apiErr *accrual.ErrUnexpectedStatus
		if errors.As(err, &apiErr) {
			log.Warn().Err(err).Msg("accrual service response:")

			return
		}
		log.Error().Err(err).Msg("accrual service response:")

		return
	}
	// If the accruals are calculated, create a new entry in accruals log
	if resp.Status == model.StatusProcessed {
		if resp.Accrual < 0 {
			if err := g.db.UpdateOrderStatus(ctx, order.ID, model.StatusInvalid); err != nil {
				log.Error().Err(err).Msg("could not update order status, operation canceled")

				return
			}
		}
		// if all is OK, the order status is set to 'PROCESSED' within db transaction.
		if err := g.db.CreateAccrual(ctx, order.ID, resp.Accrual); err != nil {
			if errors.Is(err, storage.ErrAlreadyProcessed) {
				log.Error().Err(err).Msgf("internal error! processed order has status %s", order.Status)
				return
			}
			if errors.Is(err, storage.ErrNotFound) {
				log.Error().Err(err).Msg("internal error! user not found")
				return
			}
			log.Error().Err(err).Msg("internal error:")
			return
		}
		log.Info().Float32("amount", resp.Accrual).Msg("a new entry in accruals log has been created")

		return
	}

	// Update order status. Accrual calculation hasn't processed yet.
	if resp.Status != order.Status && resp.Status != model.StatusRegistered {
		if err := g.db.UpdateOrderStatus(ctx, order.ID, resp.Status); err != nil {
			log.Error().Err(err).Msg("could not update order status, operation canceled")

			return
		}
		log.Info().Str("status", string(resp.Status)).Msg("status updated")
	}
}
