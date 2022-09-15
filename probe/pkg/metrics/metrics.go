package metrics

import (
	"sync"

	"github.com/prometheus/client_golang/prometheus"
)

const (
	prometheusNamespace = "acs"
	prometheusSubsystem = "probe"
)

var (
	metrics *Metrics
	once    sync.Once
)

// Metrics holds the prometheus.Collector instances for the probe's custom metrics
// and provides methods to interact with them.
type Metrics struct {
	startedRuns    prometheus.Counter
	successfulRuns prometheus.Counter
	failedRuns     prometheus.Counter
}

// Register registers the metrics with the given prometheus.Registerer.
func (m *Metrics) Register(r prometheus.Registerer) {
	r.MustRegister(m.startedRuns)
	r.MustRegister(m.successfulRuns)
	r.MustRegister(m.failedRuns)
}

// IncStartedRuns increments the metric counter for started probe runs.
func (m *Metrics) IncStartedRuns() {
	m.startedRuns.Inc()
}

// IncSuccessfulRuns increments the metric counter for successful probe runs.
func (m *Metrics) IncSuccessfulRuns() {
	m.successfulRuns.Inc()
}

// IncFailedRuns increments the metric counter for failed probe runs.
func (m *Metrics) IncFailedRuns() {
	m.failedRuns.Inc()
}

// MetricsInstance return the global Singleton instance for Metrics
func MetricsInstance() *Metrics {
	once.Do(initMetricsInstance)
	return metrics
}

func initMetricsInstance() {
	metrics = newMetrics()
}

// TODO: Add more metrics
func newMetrics() *Metrics {
	return &Metrics{
		startedRuns: prometheus.NewCounter(prometheus.CounterOpts{
			Namespace: prometheusNamespace,
			Subsystem: prometheusSubsystem,
			Name:      "started_runs",
			Help:      "The number of started probe runs.",
		}),
		successfulRuns: prometheus.NewCounter(prometheus.CounterOpts{
			Namespace: prometheusNamespace,
			Subsystem: prometheusSubsystem,
			Name:      "successful_runs",
			Help:      "The number of successful probe runs.",
		}),
		failedRuns: prometheus.NewCounter(prometheus.CounterOpts{
			Namespace: prometheusNamespace,
			Subsystem: prometheusSubsystem,
			Name:      "failed_runs",
			Help:      "The number of failed probe runs.",
		}),
	}
}
