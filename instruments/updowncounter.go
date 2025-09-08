package instruments

import (
	"context"
	"github.com/romansin312/metrics-exporter/tags"
	"go.opentelemetry.io/otel/metric"
)

type UpDownCounter struct {
	counter metric.Int64UpDownCounter
	meter   metric.Meter
}

func NewUpDownCounterCounter(meter metric.Meter) *UpDownCounter {
	return &UpDownCounter{meter: meter}
}

func (x *UpDownCounter) CreateInnerInstrument(key string, description string) error {
	opts := make([]metric.Int64UpDownCounterOption, 0)
	if description != "" {
		opts = append(opts, metric.WithDescription(description))
	}

	counter, err := x.meter.Int64UpDownCounter(key, opts...)
	if err != nil {
		return err
	}

	x.counter = counter
	return nil
}

func (x *UpDownCounter) Apply(key string, n interface{}, tags ...*tags.TagModel) error {
	value, err := ConvertToInt64(n)
	if err != nil {
		return err
	}

	if x.counter == nil {
		err := x.CreateInnerInstrument(key, "")
		if err != nil {
			return err
		}
	}

	attrs := ConvertTagsToAttributes(tags)
	x.counter.Add(context.Background(), value, metric.WithAttributes(attrs...))

	return nil
}

var _ IInstrument = &UpDownCounter{}
