package collector

import (
	"context"
	"database/sql"

	"github.com/go-kit/log"
	"github.com/prometheus/client_golang/prometheus"
)

const contestCountQuery = `
	SELECT COUNT(*) FROM contest
`

type ContestInfo struct{}

func (ContestInfo) Name() string {
	return "Contest Info"
}

func (ContestInfo) Help() string {
	return "Contest Info"
}

func (ContestInfo) Version() float64 {
	return 5.7
}

func (ContestInfo) Scrape(ctx context.Context, db *sql.DB, ch chan<- prometheus.Metric, logger log.Logger) error {
	CountRows, err := db.QueryContext(ctx, contestCountQuery)
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
			newDesc(contestSubsystem, "total", "total number of contest"), prometheus.CounterValue, float64(count),
		)

	}

	return nil
}

// check interface
var _ Scraper = ContestInfo{}
