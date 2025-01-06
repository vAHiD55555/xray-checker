package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	proxyStatus = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "xray_proxy_status",
			Help: "Status of proxy connection (1: success, 0: failure)",
		},
		[]string{"protocol", "address", "name"},
	)
	proxyLatency = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "xray_proxy_latency_ms",
			Help: "Latency of proxy connection in milliseconds, 0 if failed",
		},
		[]string{"protocol", "address", "name"},
	)
)

func GetProxyStatusMetric() *prometheus.GaugeVec {
	return proxyStatus
}

func GetProxyLatencyMetric() *prometheus.GaugeVec {
    return proxyLatency
}

func RecordProxyStatus(protocol, address, name string, value float64) {
	proxyStatus.WithLabelValues(protocol, address, name).Set(value)
}


func RecordProxyLatency(protocol, address, name string, value float64) {
	proxyLatency.WithLabelValues(protocol, address, name).Set(value)
}

func DeleteProxyStatus(protocol, address, name string) {
	proxyStatus.DeleteLabelValues(protocol, address, name)
}

func DeleteProxyLatency(protocol, address, name string) {
	proxyLatency.DeleteLabelValues(protocol, address, name)
}
