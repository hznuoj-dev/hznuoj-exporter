package collector

import (
	"context"
	"database/sql"

	"github.com/go-kit/log"
	"github.com/prometheus/client_golang/prometheus"
)

const userCountQuery = `
	SELECT COUNT(*) FROM users
`

type UserInfo struct{}

func (UserInfo) Name() string {
	return "User Info"
}

func (UserInfo) Help() string {
	return "User Info"
}

func (UserInfo) Version() float64 {
	return 5.7
}

func (UserInfo) Scrape(ctx context.Context, db *sql.DB, ch chan<- prometheus.Metric, logger log.Logger) error {
	userCountRows, err := db.QueryContext(ctx, userCountQuery)
	if err != nil {
		return err
	}
	defer userCountRows.Close()

	var userCount int64

	for userCountRows.Next() {
		err = userCountRows.Scan(&userCount)
		if err != nil {
			return err
		}

		ch <- prometheus.MustNewConstMetric(
			newDesc(userSubsystem, "total", "total number of users"), prometheus.CounterValue, float64(userCount),
		)

	}

	return nil
}

// check interface
var _ Scraper = ProblemInfo{}
