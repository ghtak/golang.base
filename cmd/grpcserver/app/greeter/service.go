package greeter

import (
	"context"
	pb "github.com/ghtak/golang.grpc.base/gen/go"
	"github.com/ghtak/golang.grpc.base/internal/adapter/grpcfx"
	"google.golang.org/grpc"
	"log"
)

func NewService() grpcfx.Service {
	return &server{}
}

type server struct {
	pb.UnimplementedGreeterServer
}

func (s *server) Register(svr *grpc.Server) {
	pb.RegisterGreeterServer(svr, s)
}

func (s *server) SayHello(_ context.Context, in *pb.HelloRequest) (*pb.HelloReply, error) {
	log.Printf("Received: %v", in.GetName())
	return &pb.HelloReply{Message: "Hello " + in.GetName()}, nil
}
