package instruments

import (
	"context"
	"github.com/romansin312/metrics-exporter/tags"
	"go.opentelemetry.io/otel/metric"
)

type Gauge struct {
	meter metric.Meter
	gauge metric.Float64Gauge
}

func (g *Gauge) Apply(key string, n interface{}, tags ...*tags.TagModel) error {
	value, err := convertToFloat64(n)
	if err != nil {
		return err
	}

	if g.gauge == nil {
		err := g.CreateInnerInstrument(key, "")
		if err != nil {
			return err
		}
	}

	attrs := convertTagsToAttributes(tags)
	g.gauge.Record(context.Background(), value, metric.WithAttributes(attrs...))

	return nil
}

func (g *Gauge) CreateInnerInstrument(key string, description string) error {
	newGauge, err := g.meter.Float64Gauge(key, metric.WithDescription(description))
	if err != nil {
		return err
	}

	g.gauge = newGauge
	return nil
}

func NewGauge(meter metric.Meter) *Gauge {
	return &Gauge{meter: meter}
}

var _ IInstrument = &Gauge{}
