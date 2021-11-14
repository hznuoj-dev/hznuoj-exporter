package collector

import (
	"context"
	"database/sql"

	"github.com/go-kit/log"
	"github.com/prometheus/client_golang/prometheus"
)

const submissionCountQuery = `
	SELECT COUNT(*) FROM solution
`

type SubmissionInfo struct{}

func (SubmissionInfo) Name() string {
	return "Submission Info"
}

func (SubmissionInfo) Help() string {
	return "Submission Info"
}

func (SubmissionInfo) Version() float64 {
	return 5.7
}

func (SubmissionInfo) Scrape(ctx context.Context, db *sql.DB, ch chan<- prometheus.Metric, logger log.Logger) error {
	CountRows, err := db.QueryContext(ctx, submissionCountQuery)
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
			newDesc(submissionSubsystem, "total", "total number of submission"), prometheus.CounterValue, float64(count),
		)

	}

	return nil
}

// check interface
var _ Scraper = SubmissionInfo{}
