package main

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"time"

	"github.com/gamoch/otel-examples/internal/tracing"

	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
)

const (
	serviceName    = "http-server"
	serviceVersion = "v0.1.0"

	serverAddr = "0.0.0.0:8000"
)

func SayHello(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	fmt.Println(r.Header)

	<-time.After(time.Millisecond * 500)
	ctx, greeting := GenerateGreeting(ctx)
	<-time.After(time.Millisecond * 500)

	name := r.URL.Query().Get("name")
	message := fmt.Sprintf("%s %s", greeting, name)
	w.Write([]byte(message))
}

func init() {
	rand.Seed(time.Now().Unix())
}

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	traceProvider, err := tracing.NewTraceProvider(ctx, serviceName, serviceVersion)
	if err != nil {
		log.Fatalf("failed to create trace provider: %v", err)
	}
	defer tracing.Shutdown(ctx, traceProvider, time.Second*5)

	server := &http.Server{
		Addr:    serverAddr,
		Handler: otelhttp.NewHandler(http.HandlerFunc(SayHello), "SayHello"),
	}

	log.Printf("server listening at %s", serverAddr)
	if err = server.ListenAndServe(); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

func GenerateGreeting(ctx context.Context) (context.Context, string) {
	ctx, span := tracing.Tracer.Start(ctx, "GenerateGreeting")
	defer span.End()

	greetings := []string{"Hello", "Hey", "Hi"}

	<-time.After(time.Millisecond * 500)
	return ctx, greetings[rand.Intn(len(greetings))]
}
