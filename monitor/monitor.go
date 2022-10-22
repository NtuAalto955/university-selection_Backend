package monitor

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"net/http"
)

func RefPromMonitor() *PromMonitor {
	return monitor
}

type PromMonitor struct {
	RequestCounter   *prometheus.CounterVec
	RequestHistogram *prometheus.HistogramVec
}

var monitor = &PromMonitor{}

func init() {
	// new一个counter，包含有taskName，status，errCode，msg等标签
	http_request_duration_seconds := promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name: "http_request_duration_seconds",
			Help: "http request_process time",
		},
		[]string{"caller", "method", "error_code"},
	)
	http_request_input_total := promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_request_input_total",
			Help: "The total number of http request",
		},
		[]string{"caller", "method", "error_code"},
	)
	// 注册指标，不使用默认的注册器
	reg := prometheus.NewRegistry()
	reg.MustRegister(http_request_duration_seconds, http_request_input_total)
	monitor.RequestHistogram = http_request_duration_seconds
	monitor.RequestCounter = http_request_input_total
	//
	http.Handle("/metrics", promhttp.Handler())
	http.ListenAndServe("localhost:8000", nil)

}

func (monitor *PromMonitor) ReportHttpCounter(caller, method string, errorCode string) {
	monitor.RequestCounter.WithLabelValues(caller, method, errorCode).Inc()
}

func (monitor *PromMonitor) ReportHttpHistogram(caller, method string, errorCode string, time float64) {
	monitor.RequestHistogram.WithLabelValues(caller, method, errorCode).Observe(time)
}
