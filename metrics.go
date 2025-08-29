package metrics_exporter

import (
	"github.com/romansin312/metrics-exporter/instruments"
	"github.com/romansin312/metrics-exporter/tags"
	"go.opentelemetry.io/otel/metric"
	"log"
	"time"
)

type Metrics struct {
	meter      metric.Meter
	counters   map[string]*instruments.Counter
	gauges     map[string]*instruments.Gauge
	histograms map[string]*instruments.Histogram
}

func NewMetrics(meter metric.Meter) *Metrics {
	return &Metrics{
		meter:      meter,
		counters:   make(map[string]*instruments.Counter),
		gauges:     make(map[string]*instruments.Gauge),
		histograms: make(map[string]*instruments.Histogram),
	}
}

func (m *Metrics) Increment(key string) {
	m.IncrementWithTags(key)
}

func (m *Metrics) IncrementWithTags(key string, tags ...*tags.TagModel) {
	m.CountWithTags(key, 1, tags...)
}

func (m *Metrics) CreateCounter(key string, description string) error {
	newCounter := instruments.NewCounter(m.meter)
	err := newCounter.CreateInnerInstrument(key, description)
	if err != nil {
		return err
	}

	m.counters[key] = newCounter

	return nil
}

func (m *Metrics) CreateGauge(key string, description string) error {
	newGauge := instruments.NewGauge(m.meter)
	err := newGauge.CreateInnerInstrument(key, description)
	if err != nil {
		return err
	}

	m.gauges[key] = newGauge

	return nil
}

func (m *Metrics) CreateHistogram(key string, description string) error {
	newHistogram := instruments.NewHistogram(m.meter)
	err := newHistogram.CreateInnerInstrument(key, description)
	if err != nil {
		return err
	}

	m.histograms[key] = newHistogram

	return nil
}

func (m *Metrics) CountWithTags(key string, n interface{}, tags ...*tags.TagModel) {
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

func (m *Metrics) Gauge(key string, n interface{}) {
	m.GaugeWithTags(key, n)
}

func (m *Metrics) GaugeWithTags(key string, n interface{}, tags ...*tags.TagModel) {
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

func (m *Metrics) Histogram(bucket string, v interface{}) {
	m.HistogramWithTags(bucket, v)
}

func (m *Metrics) HistogramWithTags(bucket string, v interface{}, tags ...*tags.TagModel) {
	if m.histograms[bucket] == nil {
		m.histograms[bucket] = instruments.NewHistogram(m.meter)
	}

	histogram := m.histograms[bucket]
	err := histogram.Apply(bucket, v, tags...)
	if err != nil {
		log.Printf("metrics: failed to apply histogram %s: %v", bucket, err)
	}
}

func (m *Metrics) Timing(key string, ms interface{}) {
	m.TimingWithTags(key, ms)
}

func (m *Metrics) TimingWithTags(key string, ms interface{}, tags ...*tags.TagModel) {
	m.HistogramWithTags(key+".timing", ms, tags...)
}

func (m *Metrics) Duration(key string, d time.Duration) {
	m.DurationWithTags(key, d)
}

func (m *Metrics) DurationWithTags(key string, d time.Duration, tags ...*tags.TagModel) {
	m.HistogramWithTags(key+".duration", d.Milliseconds(), tags...)
}

func (m *Metrics) Timer(key string) func() {
	return m.TimerWithTags(key)
}

func (m *Metrics) TimerWithTags(key string, tags ...*tags.TagModel) func() {
	start := time.Now()
	return func() {
		duration := time.Since(start)
		m.DurationWithTags(key, duration, tags...)
	}
}

func (m *Metrics) Flush() {
	m.histograms = make(map[string]*instruments.Histogram)
	m.counters = make(map[string]*instruments.Counter)
	m.gauges = make(map[string]*instruments.Gauge)
}

func (m *Metrics) Close() {
}
