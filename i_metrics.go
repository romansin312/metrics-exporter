package metrics_exporter

import (
	"github.com/romansin312/metrics-exporter/tags"
	"time"
)

type IMetrics interface {
	Increment(key string)
	IncrementWithTags(key string, tags ...*tags.TagModel)
	Count(key string, n interface{})
	CountWithTags(key string, n interface{}, tags ...*tags.TagModel)
	Gauge(key string, n interface{})
	GaugeWithTags(key string, n interface{}, tags ...*tags.TagModel)
	Add(key string, n interface{})
	AddWithTags(key string, n interface{}, tags ...*tags.TagModel)
	Subtract(key string, n interface{})
	SubtractWithTags(key string, n interface{}, tags ...*tags.TagModel)
	Histogram(bucket string, v interface{})
	HistogramWithTags(bucket string, v interface{}, tags ...*tags.TagModel)
	Timing(key string, ms interface{})
	TimingWithTags(key string, ms interface{}, tags ...*tags.TagModel)
	Duration(key string, d time.Duration)
	DurationWithTags(key string, d time.Duration, tags ...*tags.TagModel)
	Timer(key string) func()
	TimerWithTags(key string, tags ...*tags.TagModel) func()
	Flush()
	Close() error
}
