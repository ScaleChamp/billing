package repositories

import (
	"database/sql"
	uuid "github.com/satori/go.uuid"
	"gitlab.com/scalablespace/billing/app/models"
	"gitlab.com/scalablespace/billing/lib/components"
	"log"
)

type usageRepository struct {
	db *sql.DB
}

func (m *usageRepository) CalculateUsagePerProject() ([]*models.Project, error) {
	const update = `
UPDATE projects
SET usage = (
    	select coalesce(sum(ceil(extract(epoch from (
                        CASE WHEN ended_at BETWEEN date_trunc('month', current_timestamp) and current_timestamp THEN ended_at ELSE current_timestamp END
                        - CASE WHEN started_at BETWEEN date_trunc('month', current_timestamp) and current_timestamp THEN started_at ELSE date_trunc('month', current_timestamp) END)) / 3600) * plans.price), 0) as total
        from usages
            left join plans on plans.id = usages.plan_id
        where (ended_at is null or ended_at >= date_trunc('month', current_timestamp)) and project_id = projects.id
),
    credit = credit + CASE WHEN date_trunc('month', current_timestamp) <> date_trunc('month', coalesce(billed_at, date_trunc('month', current_timestamp - interval '1 month')))
        		THEN (
        		    	select coalesce(sum(ceil(extract(epoch from (
					CASE WHEN ended_at BETWEEN previous_payout and current_payout
					    THEN ended_at
					    ELSE current_payout END
					- CASE WHEN started_at BETWEEN previous_payout and current_payout
					    THEN started_at
					    ELSE previous_payout END)) / 3600) * plans.price), 0) as total
			 	from usages
				  left join plans on plans.id = usages.plan_id
				  left join (select date_trunc('month', current_timestamp - interval '1 month') :: TIMESTAMP as previous_payout,
							date_trunc('month', current_timestamp) as current_payout) as payouts on true
					 where (ended_at is null or ended_at >= previous_payout) and (started_at < current_payout) and project_id = projects.id
				)
        		ELSE 0 END,
    billed_at = current_timestamp
RETURNING projects.id, projects.usage, projects.credit
`
	rows, err := m.db.Query(update)
	if err != nil {
		return nil, err
	}
	projects := make([]*models.Project, 0)
	for rows.Next() {
		p := new(models.Project)
		if err := rows.Scan(&p.Id, &p.Usage, &p.Credit); err != nil {
			return nil, err
		}
		log.Println("project", p.Id.String(), "usage", p.Usage, "credit", p.Credit)
		projects = append(projects, p)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return projects, nil
}

func (m *usageRepository) UsagesCount() (int, error) {
	var count int
	if err := m.db.QueryRow(`SELECT count(*) from usages;`).Scan(&count); err != nil {
		return -1, err
	}
	return count, nil
}

func (m *usageRepository) ProjectUsagePerHour(id uuid.UUID) (float64, error) {
	var usage sql.NullFloat64
	if err := m.db.QueryRow(`SELECT sum((select price from plans where id = instances.plan_id)) FROM instances WHERE project_id = $1 and state not in (2, 5)`, id).Scan(&usage); err != nil {
		return 0, err
	}
	if usage.Valid {
		return usage.Float64, nil
	}
	return 0, nil
}

// collect which nodes available stop all insufficient funds nodes

// send email about stop

func NewUsageRepository(db *sql.DB) components.UsageRepository {
	return &usageRepository{db}
}
