package monitor

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
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
		[]string{"method", "errorCode"},
	)
	http_request_input_total := promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_request_input_total",
			Help: "The total number of http request",
		},
		[]string{"method", "errorCode"},
	)
	// 注册指标，不使用默认的注册器
	reg := prometheus.NewRegistry()
	reg.MustRegister(http_request_duration_seconds, http_request_input_total)
	monitor.RequestHistogram = http_request_duration_seconds
	monitor.RequestCounter = http_request_input_total

}

func (monitor *PromMonitor) ReportHttpCounter(method string, errorCode string) {
	monitor.RequestCounter.WithLabelValues(method, errorCode).Inc()
}

func (monitor *PromMonitor) ReportHttpHistogram(method string, errorCode string, time float64) {
	monitor.RequestHistogram.WithLabelValues(method, errorCode).Observe(time)
}
