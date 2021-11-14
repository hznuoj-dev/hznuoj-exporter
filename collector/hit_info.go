package collector

import (
	"context"
	"database/sql"

	"github.com/go-kit/log"
	"github.com/prometheus/client_golang/prometheus"
)

const hitCountQuery = `
	SELECT COUNT(*) FROM hit_log
`

type HitInfo struct{}

func (HitInfo) Name() string {
	return "Hit Info"
}

func (HitInfo) Help() string {
	return "Hit Info"
}

func (HitInfo) Version() float64 {
	return 5.7
}

func (HitInfo) Scrape(ctx context.Context, db *sql.DB, ch chan<- prometheus.Metric, logger log.Logger) error {
	hitCountRows, err := db.QueryContext(ctx, hitCountQuery)
	if err != nil {
		return err
	}
	defer hitCountRows.Close()

	var hitCount int64

	for hitCountRows.Next() {
		err = hitCountRows.Scan(&hitCount)
		if err != nil {
			return err
		}

		ch <- prometheus.MustNewConstMetric(
			newDesc(hitSubsystem, "total", "total number of hit"), prometheus.CounterValue, float64(hitCount),
		)

	}

	return nil
}

// check interface
var _ Scraper = HitInfo{}
