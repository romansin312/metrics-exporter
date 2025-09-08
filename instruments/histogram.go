package instruments

import (
	"context"
	"github.com/romansin312/metrics-exporter/tags"
	"go.opentelemetry.io/otel/metric"
)

type Histogram struct {
	meter     metric.Meter
	histogram metric.Float64Histogram
}

func (h *Histogram) Apply(key string, v interface{}, tags ...*tags.TagModel) error {
	value, err := ConvertToFloat64(v)
	if err != nil {
		return err
	}

	if h.histogram == nil {
		err = h.CreateInnerInstrument(key, "")
		if err != nil {
			return err
		}
	}

	attrs := ConvertTagsToAttributes(tags)
	h.histogram.Record(context.Background(), value, metric.WithAttributes(attrs...))

	return nil
}

func (h *Histogram) CreateInnerInstrument(key string, description string) error {
	hist, err := h.meter.Float64Histogram(key, metric.WithDescription(description))
	if err != nil {
		return err
	}

	h.histogram = hist
	return nil
}

func NewHistogram(meter metric.Meter) *Histogram {
	return &Histogram{meter: meter}
}

var _ IInstrument = &Histogram{}
