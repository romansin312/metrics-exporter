package metrics_exporter

import (
	"context"
	"fmt"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetrichttp"
	otelMetric "go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/resource"
	semconv "go.opentelemetry.io/otel/semconv/v1.17.0"
	"time"
)

type Exporter struct {
	meterProvider *metric.MeterProvider
	shutdownFunc  func(context.Context) error
}

func InitExporter(ctx context.Context, cfg Config) (*Exporter, error) {
	var exporters []metric.Exporter

	otlpExp, err := otlpmetrichttp.New(
		ctx,
		otlpmetrichttp.WithEndpoint(cfg.Endpoint),
		otlpmetrichttp.WithURLPath("/v1/metrics"),
		otlpmetrichttp.WithInsecure(),
		otlpmetrichttp.WithHeaders(cfg.Headers),
		otlpmetrichttp.WithTimeout(cfg.Timeout))

	if err != nil {
		return nil, fmt.Errorf("failed to create OTLP exporter: %w", err)
	}
	exporters = append(exporters, otlpExp)

	//if cfg.EnableConsoleExporter {
	//	consoleExp, err := stdoutmetric.New()
	//	if err != nil {
	//		log.Printf("failed to create stdout metric exporter: %v", err)
	//	} else {
	//		exporters = append(exporters, consoleExp)
	//	}
	//}

	res, err := resource.New(
		ctx,
		resource.WithAttributes(semconv.ServiceName(cfg.ServiceName)))
	if err != nil {
		return nil, fmt.Errorf("failed to create resource: %w", err)
	}

	meterProvider := metric.NewMeterProvider(
		metric.WithResource(res),
		metric.WithReader(metric.NewPeriodicReader(otlpExp, metric.WithInterval(3*time.Second))))

	shutdownFunc := func(ctx context.Context) error {
		var err error
		for _, exp := range exporters {
			if shutdownErr := exp.Shutdown(ctx); shutdownErr != nil {
				if err == nil {
					err = shutdownErr
				} else {
					err = fmt.Errorf("%v; %v", err, shutdownErr)
				}
			}
		}
		return err
	}

	return &Exporter{
		meterProvider: meterProvider,
		shutdownFunc:  shutdownFunc,
	}, nil
}

func (e *Exporter) GetMeterProvider() *metric.MeterProvider {
	return e.meterProvider
}

func (e *Exporter) GetMeter(name string) otelMetric.Meter {
	return e.meterProvider.Meter(name)
}

func (e *Exporter) CreateMetrics(meterName string) *MetricsImpl2 {
	meter := e.GetMeter(meterName)
	return NewMetrics(meter)
}

func (e *Exporter) Shutdown(ctx context.Context) error {
	if e.shutdownFunc != nil {
		return e.shutdownFunc(ctx)
	}
	return nil
}
