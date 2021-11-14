package collector

import (
	"context"
	"database/sql"

	"github.com/go-kit/log"
	"github.com/prometheus/client_golang/prometheus"
)

const loginCountQuery = `
	SELECT COUNT(*) FROM loginlog
`

type LoginInfo struct{}

func (LoginInfo) Name() string {
	return "Login Info"
}

func (LoginInfo) Help() string {
	return "Login Info"
}

func (LoginInfo) Version() float64 {
	return 5.7
}

func (LoginInfo) Scrape(ctx context.Context, db *sql.DB, ch chan<- prometheus.Metric, logger log.Logger) error {
	CountRows, err := db.QueryContext(ctx, loginCountQuery)
	if err != nil {
		return err
	}
	defer CountRows.Close()

	var count int64

	for CountRows.Next() {
		err = CountRows.Scan(&count)
		if err != nil {
			return err
		}

		ch <- prometheus.MustNewConstMetric(
			newDesc(loginSubsystem, "total", "total number of login"), prometheus.CounterValue, float64(count),
		)

	}

	return nil
}

// check interface
var _ Scraper = LoginInfo{}
