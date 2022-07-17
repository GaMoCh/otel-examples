module github.com/gamoch/otel-examples

go 1.16

require (
	go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc v0.33.0
	go.opentelemetry.io/contrib/instrumentation/net/http/httptrace/otelhttptrace v0.33.0
	go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp v0.33.0
	go.opentelemetry.io/otel v1.8.0
	go.opentelemetry.io/otel/exporters/otlp/otlptrace v1.8.0
	go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc v1.8.0
	go.opentelemetry.io/otel/sdk v1.8.0
	go.opentelemetry.io/otel/trace v1.8.0
	google.golang.org/grpc v1.48.0
	google.golang.org/protobuf v1.28.0
)
