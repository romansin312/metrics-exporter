package metrics_exporter

import (
	"context"
	"go.opentelemetry.io/otel/metric"
)

type Counter struct {
	counter metric.Int64Counter
	meter   metric.Meter
}

func NewCounter(meter metric.Meter) *Counter {
	return &Counter{meter: meter}
}

func (x Counter) Create(key string, description string) error {
	counter, err := x.meter.Int64Counter(key, metric.WithDescription(description))
	if err != nil {
		return err
	}

	x.counter = counter
	return nil
}

func (x Counter) Apply(key string, n interface{}, tags ...*TagModel) error {
	value, err := convertToInt64(n)
	if err != nil {
		return err
	}

	if x.counter == nil {
		err := x.Create(key, "")
		if err != nil {
			return err
		}
	}

	attrs := convertTagsToAttributes(tags)
	x.counter.Add(context.Background(), value, metric.WithAttributes(attrs...))

	return nil
}

var _ IInstrument = Counter{}
