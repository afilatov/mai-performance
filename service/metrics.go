package service

import (
	"github.com/prometheus/client_golang/prometheus"
)

type MetricsCollectioner interface {
	RegisterEndpoint(string) error
	CollectRequest(string)
	CollectRequestTime(string, float64)
}

type MetricsCollector struct {
	reqCounters map[string]prometheus.Counter
	reqTimingsHistograms map[string]prometheus.Histogram
}

func NewMetricsCollector() (*MetricsCollector, error) {
	return &MetricsCollector{
		reqCounters: make(map[string]prometheus.Counter),
		reqTimingsHistograms: make(map[string]prometheus.Histogram),
	}, nil
}

func (mc *MetricsCollector) RegisterEndpoint(endpoint string) error {
	if _, ok := mc.reqCounters[endpoint]; !ok {
		mc.reqCounters[endpoint] = prometheus.NewCounter(prometheus.CounterOpts{
			Name: "req_cnt_total",
			Help: "The total number of processed events",
			ConstLabels: map[string]string{"endpoint": endpoint},
		})
		if err := prometheus.Register(mc.reqCounters[endpoint]); err != nil {
			return err
		}
	}

	if _, ok := mc.reqTimingsHistograms[endpoint]; !ok {
		mc.reqTimingsHistograms[endpoint] = prometheus.NewHistogram(prometheus.HistogramOpts{
			Name: "resp_time_histogram_ms",
			Help: "Response time histogram",
			Buckets: prometheus.LinearBuckets(0, 10, 20),
			ConstLabels: map[string]string{"endpoint": endpoint},
		})
		if err := prometheus.Register(mc.reqTimingsHistograms[endpoint]); err != nil {
			return err
		}
	}

	return nil
}

func (mc *MetricsCollector) CollectRequest(endpoint string) {
	c, ok := mc.reqCounters[endpoint]
	if !ok {
		return
	}

	c.Inc()
}

func (mc *MetricsCollector) CollectRequestTime(endpoint string, t float64) {
	c, ok := mc.reqTimingsHistograms[endpoint]
	if !ok {
		return
	}

	c.Observe(t)
}
