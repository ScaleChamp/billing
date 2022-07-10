package config

import (
	"context"
	"database/sql"
	"fmt"
	"gitlab.com/scalablespace/billing/app/services"
	"gitlab.com/scalablespace/billing/lib/components"
	"go.uber.org/fx"
	"go.uber.org/multierr"
	"sync"
	"time"
)

func boot(lc fx.Lifecycle, db *sql.DB, repository components.UsageRepository) {
	ticker := time.NewTicker(30 * time.Minute)
	var wg sync.WaitGroup
	stop := make(chan struct{})

	scheduler := services.Scheduler{
		Ticker:     ticker,
		Repository: repository,
		Publisher:  new(services.Publisher),
	}
	lc.Append(fx.Hook{
		OnStart: func(context.Context) error {
			wg.Add(1)
			go func() {
				defer wg.Done()
				scheduler.Start(stop)
			}()
			return nil
		},
		OnStop: func(ctx context.Context) error {
			fmt.Println("closing")
			ticker.Stop()
			close(stop)
			fmt.Println("closing")
			wg.Wait()
			fmt.Println("closing")
			return multierr.Combine(db.Close())
		},
	})
}
