package main

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"net"
	"time"

	"github.com/gamoch/otel-examples/internal/grpc/pb"

	"github.com/gamoch/otel-examples/internal/tracing"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"

	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
)

const (
	serviceName    = "grpc-server"
	serviceVersion = "v0.1.0"

	serverAddr = "0.0.0.0:8001"
)

type Server struct {
	pb.UnimplementedGreeterServer
}

func (s *Server) SayHello(ctx context.Context, in *pb.HelloRequest) (*pb.HelloReply, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if ok {
		fmt.Println(md)
	}

	<-time.After(time.Millisecond * 500)
	ctx, greeting := GenerateGreeting(ctx)
	<-time.After(time.Millisecond * 500)

	name := in.GetName()
	message := fmt.Sprintf("%s %s", greeting, name)
	return &pb.HelloReply{Message: message}, nil
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

	listener, err := net.Listen("tcp", serverAddr)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	server := grpc.NewServer(
		grpc.UnaryInterceptor(otelgrpc.UnaryServerInterceptor()),
		grpc.StreamInterceptor(otelgrpc.StreamServerInterceptor()),
	)
	pb.RegisterGreeterServer(server, &Server{})

	log.Printf("server listening at %s", serverAddr)
	if err = server.Serve(listener); err != nil {
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
