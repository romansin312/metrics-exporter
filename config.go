package metrics_exporter

import "time"

type Config struct {
	Endpoint              string
	Headers               map[string]string
	Timeout               time.Duration
	ServiceName           string
	EnableConsoleExporter bool
}
