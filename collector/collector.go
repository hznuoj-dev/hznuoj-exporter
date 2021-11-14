package collector

import "github.com/prometheus/client_golang/prometheus"

var (
	Namespace = ""
)

func newDesc(subSystem, name, help string) *prometheus.Desc {
	return prometheus.NewDesc(
		prometheus.BuildFQName(Namespace, subSystem, name),
		help, nil, nil,
	)
}
