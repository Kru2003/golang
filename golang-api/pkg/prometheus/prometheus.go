package prometheus

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

const Namespace = "golang_api"

type PrometheusMetrics struct {
	MoviesMetrics   prometheus.Gauge
	RequestsMetrics *prometheus.CounterVec
}

var metrics *PrometheusMetrics = nil

func InitPrometheusMetrics() *PrometheusMetrics {
	if metrics == nil {
		metrics = &PrometheusMetrics{
			MoviesMetrics: promauto.NewGauge(prometheus.GaugeOpts{
				Namespace: Namespace,
				Name:      "movies_total",
				Help:      "Total movies",
			}),
			RequestsMetrics: promauto.NewCounterVec(prometheus.CounterOpts{
				Namespace: Namespace,
				Name:      "requests_total",
				Help:      "Total http requests",
			}, []string{"code"}),
		}
	}

	return metrics
}
