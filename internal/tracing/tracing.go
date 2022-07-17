package tracing

import (
	"context"
	"fmt"
	"log"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.10.0"
	"go.opentelemetry.io/otel/trace"
)

const (
	otelCollectorAddr = "localhost:4317"

	instrumentationName    = "otel/tracing"
	instrumentationVersion = "semver:0.1.0"
)

var Tracer = otel.GetTracerProvider().Tracer(
	instrumentationName,
	trace.WithInstrumentationVersion(instrumentationVersion),
	trace.WithSchemaURL(semconv.SchemaURL),
)

func NewTraceProvider(ctx context.Context, serviceName, serviceVersion string, attrs ...attribute.KeyValue) (*sdktrace.TracerProvider, error) {
	client := otlptracegrpc.NewClient(
		otlptracegrpc.WithInsecure(),
		otlptracegrpc.WithEndpoint(otelCollectorAddr),
	)

	exporter, err := otlptrace.New(ctx, client)
	if err != nil {
		return nil, fmt.Errorf("failed to create OTLP trace exporter: %w", err)
	}

	resourceAttributes := append([]attribute.KeyValue{
		semconv.ServiceNameKey.String(serviceName),
		semconv.ServiceVersionKey.String(serviceVersion),
	}, attrs...)

	tracerProvider := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(exporter),
		sdktrace.WithSampler(sdktrace.AlwaysSample()),
		sdktrace.WithResource(resource.NewWithAttributes(
			semconv.SchemaURL,
			resourceAttributes...,
		)),
	)

	textMapPropagator := propagation.NewCompositeTextMapPropagator(
		propagation.TraceContext{},
		propagation.Baggage{},
	)

	otel.SetTracerProvider(tracerProvider)
	otel.SetTextMapPropagator(textMapPropagator)

	return tracerProvider, nil
}

func Shutdown(ctx context.Context, provider *sdktrace.TracerProvider, duration time.Duration) {
	ctx, cancel := context.WithTimeout(ctx, duration)
	defer cancel()
	if err := provider.Shutdown(ctx); err != nil {
		log.Fatal(err)
	}
}
