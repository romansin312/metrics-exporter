package metrics_exporter

import (
	"context"
	"fmt"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetrichttp"
	"go.opentelemetry.io/otel/exporters/stdout/stdoutmetric"
	otelMetric "go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/resource"
	semconv "go.opentelemetry.io/otel/semconv/v1.17.0"
	"log"
	"time"
)

type exporter struct {
	meterProvider *metric.MeterProvider
	shutdownFunc  func(context.Context) error
}

func newExporter(ctx context.Context, cfg *config) (*exporter, error) {
	var exporters []metric.Exporter

	otlpExp, err := otlpmetrichttp.New(ctx, cfg.build()...)

	if err != nil {
		return nil, fmt.Errorf("failed to create OTLP exporter: %w", err)
	}
	exporters = append(exporters, otlpExp)

	if cfg.enableConsoleExporter {
		consoleExp, err := stdoutmetric.New()
		if err != nil {
			log.Printf("failed to create stdout metric exporter: %v", err)
		} else {
			exporters = append(exporters, consoleExp)
		}
	}

	res, err := resource.New(
		ctx,
		resource.WithAttributes(semconv.ServiceName(cfg.serviceName)))

	if err != nil {
		return nil, fmt.Errorf("failed to create resource: %w", err)
	}

	sendingInterval := time.Second * 3
	if cfg.sendingInterval > 0 {
		sendingInterval = cfg.sendingInterval
	}

	readers := make([]metric.Option, 0)
	readers = append(readers, metric.WithResource(res))
	for _, exp := range exporters {
		readers = append(readers, metric.WithReader(metric.NewPeriodicReader(exp, metric.WithInterval(sendingInterval))))
	}

	meterProvider := metric.NewMeterProvider(readers...)

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

	return &exporter{
		meterProvider: meterProvider,
		shutdownFunc:  shutdownFunc,
	}, nil
}

func (e *exporter) getMeter(name string) otelMetric.Meter {
	return e.meterProvider.Meter(name)
}

func (e *exporter) Shutdown(ctx context.Context) error {
	if e.shutdownFunc != nil {
		return e.shutdownFunc(ctx)
	}
	return nil
}
