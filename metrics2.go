package metrics_exporter

import (
	"go.opentelemetry.io/otel/metric"
	"log"
	"time"
)

type MetricsImpl2 struct {
	meter      metric.Meter
	counters   map[string]*Counter
	gauges     map[string]*Gauge
	histograms map[string]*Histogram
}

func NewMetrics(meter metric.Meter) *MetricsImpl2 {
	return &MetricsImpl2{
		meter:      meter,
		counters:   make(map[string]*Counter),
		gauges:     make(map[string]*Gauge),
		histograms: make(map[string]*Histogram),
	}
}

func (m *MetricsImpl2) Increment(key string) {
	m.IncrementWithTags(key)
}

func (m *MetricsImpl2) IncrementWithTags(key string, tags ...*TagModel) {
	m.CountWithTags(key, 1, tags...)
}

func (m *MetricsImpl2) ConfigureCounter(key string, description string) {

}

func (m *MetricsImpl2) CountWithTags(key string, n interface{}, tags ...*TagModel) {
	counter := m.counters[key]
	if counter == nil {
		m.counters[key] = NewCounter(m.meter)
		counter = m.counters[key]
	}

	err := counter.Apply(key, n, tags...)
	if err != nil {
		log.Printf("metrics: failed to apply counter %s: %v", key, err)
	}
}

func (m *MetricsImpl2) Gauge(key string, n interface{}) {
	m.GaugeWithTags(key, n)
}

func (m *MetricsImpl2) GaugeWithTags(key string, n interface{}, tags ...*TagModel) {
	gauge := m.gauges[key]
	if gauge == nil {
		gauge = NewGauge(m.meter)
		m.gauges[key] = gauge
	}

	err := gauge.Apply(key, n, tags...)
	if err != nil {
		log.Printf("metrics: failed to apply gauge %s: %v", key, err)
	}
}

func (m *MetricsImpl2) Histogram(bucket string, v interface{}) {
	m.HistogramWithTags(bucket, v)
}

func (m *MetricsImpl2) HistogramWithTags(bucket string, v interface{}, tags ...*TagModel) {
	if m.histograms[bucket] == nil {
		m.histograms[bucket] = NewHistogram(m.meter)
	}

	histogram := m.histograms[bucket]
	err := histogram.Apply(bucket, v, tags...)
	if err != nil {
		log.Printf("metrics: failed to apply histogram %s: %v", bucket, err)
	}
}

func (m *MetricsImpl2) Timing(key string, ms interface{}) {
	m.TimingWithTags(key, ms)
}

func (m *MetricsImpl2) TimingWithTags(key string, ms interface{}, tags ...*TagModel) {
	m.HistogramWithTags(key+".timing", ms, tags...)
}

func (m *MetricsImpl2) Duration(key string, d time.Duration) {
	m.DurationWithTags(key, d)
}

func (m *MetricsImpl2) DurationWithTags(key string, d time.Duration, tags ...*TagModel) {
	m.HistogramWithTags(key+".duration", d.Milliseconds(), tags...)
}

func (m *MetricsImpl2) Timer(key string) func() {
	return m.TimerWithTags(key)
}

func (m *MetricsImpl2) TimerWithTags(key string, tags ...*TagModel) func() {
	start := time.Now()
	return func() {
		duration := time.Since(start)
		m.DurationWithTags(key, duration, tags...)
	}
}

func (m *MetricsImpl2) Flush() {
	m.histograms = make(map[string]*Histogram)
	m.counters = make(map[string]*Counter)
	m.gauges = make(map[string]*Gauge)
}

func (m *MetricsImpl2) Close() {
}
