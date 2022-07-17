package main

import (
	"context"
	"io"
	"log"
	"math/rand"
	"net/http"
	"net/http/httptrace"
	"time"

	"github.com/gamoch/otel-examples/internal/tracing"

	"go.opentelemetry.io/contrib/instrumentation/net/http/httptrace/otelhttptrace"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
)

const (
	serviceName    = "http-client"
	serviceVersion = "v0.1.0"

	serverAddr = "http://localhost:8000"
)

func init() {
	rand.Seed(time.Now().Unix())
}

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	traceProvider, err := tracing.NewTraceProvider(ctx, serviceName, serviceVersion)
	if err != nil {
		log.Fatalf("failed to create trace provider: %v", err)
	}
	defer tracing.Shutdown(ctx, traceProvider, time.Second*5)

	ctx, name := GenerateName(ctx)
	req, _ := http.NewRequestWithContext(ctx, http.MethodGet, serverAddr, nil)

	query := req.URL.Query()
	query.Set("name", name)
	req.URL.RawQuery = query.Encode()

	client := http.Client{Transport: otelhttp.NewTransport(http.DefaultTransport)}
	ctx = httptrace.WithClientTrace(ctx, otelhttptrace.NewClientTrace(ctx))

	res, err := client.Do(req)
	if err != nil {
		log.Fatalf("could not greet: %v", err)
	}
	defer res.Body.Close()

	message, _ := io.ReadAll(res.Body)
	log.Printf("Greeting: %s", string(message))
}

func GenerateName(ctx context.Context) (context.Context, string) {
	ctx, span := tracing.Tracer.Start(ctx, "GenerateName")
	defer span.End()

	names := []string{"John", "Jane", "Johnny"}

	<-time.After(time.Millisecond * 500)
	return ctx, names[rand.Intn(len(names))]
}
