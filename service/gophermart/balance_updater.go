package gophermart

import (
	"context"
	"time"

	appContext "github.com/vanamelnik/gophermart/pkg/ctx"
)

// balanceUpdater looks for unfinished acrrual operations and updates user balance with accrual points.
func (g *GopherMart) balanceUpdater(ctx context.Context) {
	log := appContext.Logger(ctx).With().Str("service:", "balanceUpdater").Logger()
	log.Info().Msg("balanceUpdater started")
	t := time.NewTicker(g.balanceUpdInterval)
loop:
	for {
		select {
		case <-t.C:
			n, err := g.db.UpdateBalance(ctx)
			if err != nil {
				log.Error().Err(err).Msg("")

				continue
			}
			if n > 0 {
				log.Info().Int("number of accrual operations processed", n).Msg("")
			}
		case <-g.workersStop:
			break loop
		}
	}
	g.workersWg.Done()
	log.Info().Msg("balanceUpdater stopped")
}
