package config

import (
	"gitlab.com/scalablespace/billing/db"
	"gitlab.com/scalablespace/billing/db/repositories"
	"go.uber.org/fx"
)

func NewApp() *fx.App {
	return fx.New(
		fx.Provide(repositories.NewUsageRepository),
		fx.Provide(newEnvironment),
		fx.Provide(db.NewDB),
		fx.Invoke(boot),
	)
}
