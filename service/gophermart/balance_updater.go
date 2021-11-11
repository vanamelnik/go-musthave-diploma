package gophermart

import (
	"context"
	"time"

	appContext "github.com/vanamelnik/go-musthave-diploma/pkg/ctx"
)

// balanceUpdater looks for unfinished acrrual operations and updates user balance with accrual points.
func (g *GopherMart) balanceUpdater(ctx context.Context) {
	log := appContext.Logger(ctx)
	const logPrefix = "balanceUpdater:"
	log.Info().Msg("balanceUpdater started")
	t := time.NewTicker(g.balanceUpdInterval)
loop:
	for {
		select {
		case <-t.C:
			n, err := g.db.UpdateBalance(ctx)
			if err != nil {
				log.Error().Err(err).Msgf("%s", logPrefix)

				continue
			}
			log.Info().Int("number of accrual operations processed", n).Msgf("%s", logPrefix)
		case <-ctx.Done():
			break loop
		}
	}
	log.Info().Msg("balanceUpdater stopped")
}
