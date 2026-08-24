package metrics

import (
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	PagesCrawled = promauto.NewCounter(prometheus.CounterOpts{
		Name: "aranea_pages_crawled_total",
		Help: "Total number of pages successfully fetched and processed.",
	})

	FetchErrors = promauto.NewCounter(prometheus.CounterOpts{
		Name: "aranea_fetch_errors_total",
		Help: "Total number of page fetch failures.",
	})

	QueueDepth = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "aranea_queue_depth",
		Help: "Number of URLs waiting in the crawl queue.",
	})

	FetchDuration = promauto.NewHistogram(prometheus.HistogramOpts{
		Name:    "aranea_fetch_duration_seconds",
		Help:    "Time spent fetching a page.",
		Buckets: prometheus.DefBuckets,
	})
)

func Serve(addr string) error {
	mux := http.NewServeMux()
	mux.Handle("/metrics", promhttp.Handler())
	return http.ListenAndServe(addr, mux)
}
