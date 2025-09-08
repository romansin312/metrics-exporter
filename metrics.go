package metrics_exporter

import (
	"context"
	"github.com/romansin312/metrics-exporter/instruments"
	"github.com/romansin312/metrics-exporter/tags"
	"go.opentelemetry.io/otel/metric"
	"log"
	"sync"
	"time"
)

type Metrics struct {
	exporter       *exporter
	upDownCounters map[string]*instruments.UpDownCounter
	counters       map[string]*instruments.Counter
	gauges         map[string]*instruments.Gauge
	histograms     map[string]*instruments.Histogram
	meter          metric.Meter
	mu             sync.RWMutex
}

func NewMetrics(ctx context.Context, opts ...configOption) (*Metrics, error) {
	config, err := newConfig(opts...)
	if err != nil {
		return nil, err
	}

	exporter, err := newExporter(ctx, config)
	if err != nil {
		return nil, err
	}
	meter := exporter.getMeter(config.scope)
	return &Metrics{
		exporter:       exporter,
		meter:          meter,
		counters:       make(map[string]*instruments.Counter),
		gauges:         make(map[string]*instruments.Gauge),
		histograms:     make(map[string]*instruments.Histogram),
		upDownCounters: make(map[string]*instruments.UpDownCounter),
	}, nil
}

func (m *Metrics) Increment(key string) {
	m.IncrementWithTags(key)
}

func (m *Metrics) IncrementWithTags(key string, tags ...*tags.TagModel) {
	m.CountWithTags(key, 1, tags...)
}

func (m *Metrics) Count(key string, n interface{}) {
	m.CountWithTags(key, n)
}
func (m *Metrics) CountWithTags(key string, n interface{}, tags ...*tags.TagModel) {
	counter := getOrCreateInstrument(
		&m.mu,
		m.counters,
		key,
		func() *instruments.Counter { return instruments.NewCounter(m.meter) },
	)

	err := counter.Apply(key, n, tags...)
	if err != nil {
		log.Printf("metrics: failed to apply counter %s: %v", key, err)
	}
}

func (m *Metrics) Add(key string, n interface{}) {
	m.AddWithTags(key, n)
}
func (m *Metrics) AddWithTags(key string, n interface{}, tags ...*tags.TagModel) {
	upDownCounter := getOrCreateInstrument(
		&m.mu,
		m.upDownCounters,
		key,
		func() *instruments.UpDownCounter { return instruments.NewUpDownCounterCounter(m.meter) },
	)

	err := upDownCounter.Apply(key, n, tags...)
	if err != nil {
		log.Printf("metrics: failed to apply updowncounter %s: %v", key, err)
	}
}

func (m *Metrics) Subtract(key string, n interface{}) {
	m.SubtractWithTags(key, n)
}
func (m *Metrics) SubtractWithTags(key string, n interface{}, tags ...*tags.TagModel) {
	upDownCounter := getOrCreateInstrument(
		&m.mu,
		m.upDownCounters,
		key,
		func() *instruments.UpDownCounter { return instruments.NewUpDownCounterCounter(m.meter) },
	)

	convertedValue, err := instruments.ConvertToInt64(n)
	if err != nil {
		log.Printf("metrics: failed to convert value to int64: %v", err)
		return
	}

	err = upDownCounter.Apply(key, -convertedValue, tags...)
	if err != nil {
		log.Printf("metrics: failed to apply updowncounter %s: %v", key, err)
	}
}

func (m *Metrics) Gauge(key string, n interface{}) {
	m.GaugeWithTags(key, n)
}

func (m *Metrics) GaugeWithTags(key string, n interface{}, tags ...*tags.TagModel) {
	gauge := getOrCreateInstrument(
		&m.mu,
		m.gauges,
		key,
		func() *instruments.Gauge { return instruments.NewGauge(m.meter) },
	)

	err := gauge.Apply(key, n, tags...)
	if err != nil {
		log.Printf("metrics: failed to apply gauge %s: %v", key, err)
	}
}

func (m *Metrics) Histogram(bucket string, v interface{}) {
	m.HistogramWithTags(bucket, v)
}

func (m *Metrics) HistogramWithTags(key string, v interface{}, tags ...*tags.TagModel) {
	histogram := getOrCreateInstrument(
		&m.mu,
		m.histograms,
		key,
		func() *instruments.Histogram { return instruments.NewHistogram(m.meter) },
	)

	err := histogram.Apply(key, v, tags...)
	if err != nil {
		log.Printf("metrics: failed to apply histogram %s: %v", key, err)
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
	m.mu.Lock()
	defer m.mu.Unlock()

	m.histograms = make(map[string]*instruments.Histogram)
	m.counters = make(map[string]*instruments.Counter)
	m.gauges = make(map[string]*instruments.Gauge)
}

func (e *exporter) ForceFlush(ctx context.Context) error {
	if e.meterProvider != nil {
		return e.meterProvider.ForceFlush(ctx)
	}
	return nil
}

func (m *Metrics) Close() {
	m.mu.Lock()
	defer m.mu.Unlock()

	if m.exporter != nil {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		if err := m.exporter.ForceFlush(ctx); err != nil {
			log.Printf("metrics: failed to force flush metrics: %v", err)
		}

		if err := m.exporter.Shutdown(ctx); err != nil {
			log.Printf("metrics: failed to shutdown exporter gracefully: %v", err)
		}
	}
}
func getOrCreateInstrument[T any](
	mu *sync.RWMutex,
	instrumentMap map[string]*T,
	key string,
	createFunc func() *T,
) *T {
	mu.RLock()
	instrument := instrumentMap[key]
	mu.RUnlock()

	if instrument == nil {
		mu.Lock()
		defer mu.Unlock()

		if instrumentMap[key] == nil {
			instrumentMap[key] = createFunc()
		}
		instrument = instrumentMap[key]
	}

	return instrument
}

var _ IMetrics = &Metrics{}
