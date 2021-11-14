package collector

import (
	"context"
	"database/sql"

	"github.com/go-kit/log"
	"github.com/prometheus/client_golang/prometheus"
)

const problemCountQuery = `
	SELECT COUNT(*) FROM problem
`

type ProblemInfo struct{}

func (ProblemInfo) Name() string {
	return "Problem Info"
}

func (ProblemInfo) Help() string {
	return "Problem Info"
}

func (ProblemInfo) Version() float64 {
	return 5.7
}

func (ProblemInfo) Scrape(ctx context.Context, db *sql.DB, ch chan<- prometheus.Metric, logger log.Logger) error {
	problemCountRows, err := db.QueryContext(ctx, problemCountQuery)
	if err != nil {
		return err
	}
	defer problemCountRows.Close()

	var problemCount int64

	for problemCountRows.Next() {
		err = problemCountRows.Scan(&problemCount)
		if err != nil {
			return err
		}

		ch <- prometheus.MustNewConstMetric(
			newDesc(problemSubsystem, "total", "total number of problems"), prometheus.CounterValue, float64(problemCount),
		)

	}

	return nil
}

// check interface
var _ Scraper = ProblemInfo{}
