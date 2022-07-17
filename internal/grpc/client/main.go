package main

import (
	"context"
	"log"
	"math/rand"
	"time"

	"github.com/gamoch/otel-examples/internal/grpc/pb"

	"github.com/gamoch/otel-examples/internal/tracing"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
)

const (
	serviceName    = "grpc-client"
	serviceVersion = "v0.1.0"

	serverAddr = "localhost:8001"
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

	conn, err := grpc.Dial(serverAddr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithUnaryInterceptor(otelgrpc.UnaryClientInterceptor()),
		grpc.WithStreamInterceptor(otelgrpc.StreamClientInterceptor()),
	)
	if err != nil {
		log.Fatalf("failed to create client connection: %v", err)
	}
	defer conn.Close()
	client := pb.NewGreeterClient(conn)

	ctx, name := GenerateName(ctx)
	r, err := client.SayHello(ctx, &pb.HelloRequest{Name: name})
	if err != nil {
		log.Fatalf("could not greet: %v", err)
	}

	log.Printf("Greeting: %s", r.GetMessage())
}

func GenerateName(ctx context.Context) (context.Context, string) {
	ctx, span := tracing.Tracer.Start(ctx, "GenerateName")
	defer span.End()

	names := []string{"John", "Jane", "Johnny"}

	<-time.After(time.Millisecond * 500)
	return ctx, names[rand.Intn(len(names))]
}
