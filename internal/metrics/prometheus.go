package metrics

import "github.com/prometheus/client_golang/prometheus"

var (
	EmailsSent = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "emails_sent_total",
	})

	EmailsFailed = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "emails_failed_total",
	})

	EmailProcessingDuration = prometheus.NewHistogram(prometheus.HistogramOpts{
		Name: "email_processing_duration_seconds",
	})
)

func init() {
	prometheus.MustRegister(EmailsSent)
	prometheus.MustRegister(EmailsFailed)
	prometheus.MustRegister(EmailProcessingDuration)
}
