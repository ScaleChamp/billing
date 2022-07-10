package components

import (
	uuid "github.com/satori/go.uuid"
	"gitlab.com/scalablespace/billing/app/models"
)

type UsageRepository interface {
	CalculateUsagePerProject() ([]*models.Project, error)
	UsagesCount() (int, error)
	ProjectUsagePerHour(id uuid.UUID) (float64, error)
}
