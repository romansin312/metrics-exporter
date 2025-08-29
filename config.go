package metrics_exporter

import "time"

type Config struct {
	Endpoint              string
	UrlPath               string
	Headers               map[string]string
	Timeout               time.Duration
	SendingInterval       time.Duration
	ServiceName           string
	EnableConsoleExporter bool
}
