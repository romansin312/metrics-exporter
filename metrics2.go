package metrics_exporter

import (
	"github.com/romansin312/metrics-exporter/instruments"
	"github.com/romansin312/metrics-exporter/tags"
	"go.opentelemetry.io/otel/metric"
	"log"
	"time"
)

type MetricsImpl2 struct {
	meter      metric.Meter
	counters   map[string]*instruments.Counter
	gauges     map[string]*instruments.Gauge
	histograms map[string]*instruments.Histogram
}

func NewMetrics(meter metric.Meter) *MetricsImpl2 {
	return &MetricsImpl2{
		meter:      meter,
		counters:   make(map[string]*instruments.Counter),
		gauges:     make(map[string]*instruments.Gauge),
		histograms: make(map[string]*instruments.Histogram),
	}
}

func (m *MetricsImpl2) Increment(key string) {
	m.IncrementWithTags(key)
}

func (m *MetricsImpl2) IncrementWithTags(key string, tags ...*tags.TagModel) {
	m.CountWithTags(key, 1, tags...)
}

func (m *MetricsImpl2) ConfigureCounter(key string, description string) {

}

func (m *MetricsImpl2) CountWithTags(key string, n interface{}, tags ...*tags.TagModel) {
	counter := m.counters[key]
	if counter == nil {
		m.counters[key] = instruments.NewCounter(m.meter)
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

func (m *MetricsImpl2) GaugeWithTags(key string, n interface{}, tags ...*tags.TagModel) {
	gauge := m.gauges[key]
	if gauge == nil {
		gauge = instruments.NewGauge(m.meter)
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

func (m *MetricsImpl2) HistogramWithTags(bucket string, v interface{}, tags ...*tags.TagModel) {
	if m.histograms[bucket] == nil {
		m.histograms[bucket] = instruments.NewHistogram(m.meter)
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

func (m *MetricsImpl2) TimingWithTags(key string, ms interface{}, tags ...*tags.TagModel) {
	m.HistogramWithTags(key+".timing", ms, tags...)
}

func (m *MetricsImpl2) Duration(key string, d time.Duration) {
	m.DurationWithTags(key, d)
}

func (m *MetricsImpl2) DurationWithTags(key string, d time.Duration, tags ...*tags.TagModel) {
	m.HistogramWithTags(key+".duration", d.Milliseconds(), tags...)
}

func (m *MetricsImpl2) Timer(key string) func() {
	return m.TimerWithTags(key)
}

func (m *MetricsImpl2) TimerWithTags(key string, tags ...*tags.TagModel) func() {
	start := time.Now()
	return func() {
		duration := time.Since(start)
		m.DurationWithTags(key, duration, tags...)
	}
}

func (m *MetricsImpl2) Flush() {
	m.histograms = make(map[string]*instruments.Histogram)
	m.counters = make(map[string]*instruments.Counter)
	m.gauges = make(map[string]*instruments.Gauge)
}

func (m *MetricsImpl2) Close() {
}
