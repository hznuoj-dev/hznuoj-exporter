package main

import (
	"net/http"
	"os"

	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
	_ "github.com/go-sql-driver/mysql"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/prometheus/common/promlog"
	"github.com/prometheus/common/promlog/flag"
	"github.com/prometheus/common/version"
	"github.com/prometheus/exporter-toolkit/web"
	webflag "github.com/prometheus/exporter-toolkit/web/kingpinflag"
	kingpin "gopkg.in/alecthomas/kingpin.v2"

	"hznuoj_exporter/collector"
)

const (
	Namespace = collector.Namespace + "_exporter"
)

var (
	webConfig = webflag.AddFlags(kingpin.CommandLine)

	listenAddress = kingpin.Flag(
		"web.listen-address",
		"Address to listen on for web interface and telemetry.",
	).Default(":9800").String()

	metricPath = kingpin.Flag(
		"web.telemetry-path",
		"Path under which to expose metrics.",
	).Default("/metrics").String()

	dbType = kingpin.Flag(
		"db.type",
		"DB Type",
	).Default("mysql").String()

	dbConnectString = kingpin.Flag(
		"db.connect.string",
		"DB Connect String",
	).Default("root:root@tcp(127.0.0.1:3306)/jol?charset=utf8&parseTime=True").String()
)

func init() {
	prometheus.MustRegister(version.NewCollector(Namespace))
}

// db, err := sql.Open(*dbType, *dbConnectString)
// if err != nil {
// 	panic(err)
// }

// defer db.Close()

// rows := db.QueryRow("SELECT COUNT(*) FROM hit_log")
// var count int64
// rows.Scan(&count)
// fmt.Println(count)

func newHandler(metrics collector.Metrics, scrapers []collector.Scraper, logger log.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		dsn := *dbConnectString
		dbType := *dbType

		registry := prometheus.NewRegistry()
		registry.MustRegister(collector.New(ctx, dbType, dsn, metrics, scrapers, logger))

		gatherers := prometheus.Gatherers{
			prometheus.DefaultGatherer,
			registry,
		}

		h := promhttp.HandlerFor(gatherers, promhttp.HandlerOpts{})
		h.ServeHTTP(w, r)
	}
}

func main() {
	// Parse flags.
	promlogConfig := &promlog.Config{}
	flag.AddFlags(kingpin.CommandLine, promlogConfig)
	kingpin.Version(version.Print(Namespace))
	kingpin.HelpFlag.Short('h')
	kingpin.Parse()
	logger := promlog.New(promlogConfig)

	level.Info(logger).Log("msg", "Starting hznuoj_exporter", "version", version.Info())
	level.Info(logger).Log("msg", "Build context", version.BuildContext())

	var landingPage = []byte(`<html>
	<head><title>HZNUOJ exporter</title></head>
	<body>
	<h1>HZNUOJ exporter</h1>
	<p><a href='` + *metricPath + `'>Metrics</a></p>
	</body>
	</html>
	`)

	scrapers := []collector.Scraper{
		collector.ProblemInfo{},
	}
	handlerFunc := newHandler(collector.NewMetrics(), scrapers, logger)
	http.Handle(*metricPath, promhttp.InstrumentMetricHandler(prometheus.DefaultRegisterer, handlerFunc))
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write(landingPage)
	})

	level.Info(logger).Log("msg", "Listening on address", "address", *listenAddress)
	srv := &http.Server{Addr: *listenAddress}
	if err := web.ListenAndServe(srv, *webConfig, logger); err != nil {
		level.Error(logger).Log("msg", "Error starting HTTP server", "err", err)
		os.Exit(1)
	}
}
