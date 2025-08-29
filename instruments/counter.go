package instruments

import (
	"context"
	"github.com/romansin312/metrics-exporter/tags"
	"go.opentelemetry.io/otel/metric"
)

type Counter struct {
	counter metric.Int64Counter
	meter   metric.Meter
}

func NewCounter(meter metric.Meter) *Counter {
	return &Counter{meter: meter}
}

func (x *Counter) CreateInnerInstrument(key string, description string) error {
	opts := make([]metric.Int64CounterOption, 0)
	if description != "" {
		opts = append(opts, metric.WithDescription(description))
	}

	counter, err := x.meter.Int64Counter(key, opts...)
	if err != nil {
		return err
	}

	x.counter = counter
	return nil
}

func (x *Counter) Apply(key string, n interface{}, tags ...*tags.TagModel) error {
	value, err := convertToInt64(n)
	if err != nil {
		return err
	}

	if x.counter == nil {
		err := x.CreateInnerInstrument(key, "")
		if err != nil {
			return err
		}
	}

	attrs := convertTagsToAttributes(tags)
	x.counter.Add(context.Background(), value, metric.WithAttributes(attrs...))

	return nil
}

var _ IInstrument = &Counter{}
