package config

import (
	"github.com/caarlos0/env/v6"
	"gitlab.com/scalablespace/billing/app/models"
)

func newEnvironment() (models.Environment, error) {
	var e models.Environment
	return e, env.Parse(&e)
}
