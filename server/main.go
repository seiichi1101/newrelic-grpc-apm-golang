package main

import (
	"context"
	"log"
	"net"
	"os"

	pb "newrelic-go-grpc/proto"

	newrelic "github.com/newrelic/go-agent"
	"github.com/newrelic/go-agent/_integrations/nrgrpc"
	"google.golang.org/grpc"
)

const (
	port = ":50051"
)

type server struct {
	pb.UnimplementedGreeterServer
}

func (s *server) SayHello(ctx context.Context, in *pb.HelloRequest) (*pb.HelloReply, error) {
	return &pb.HelloReply{Message: "Hello World"}, nil
}

func main() {
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	cfg := newrelic.NewConfig("gRPC Server", os.Getenv("NEW_RELIC_LICENSE_KEY"))
	app, _ := newrelic.NewApplication(cfg)

	s := grpc.NewServer(
		grpc.UnaryInterceptor(
			nrgrpc.UnaryServerInterceptor(app),
		),
	)
	pb.RegisterGreeterServer(s, &server{})
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
