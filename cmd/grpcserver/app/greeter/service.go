package greeter

import (
	"context"
	pb "github.com/ghtak/golang.grpc.base/gen/go"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
	"log"
)

func NewService() *Server {
	return &Server{}
}

type Server struct {
	pb.UnimplementedGreeterServer
}

func (s *Server) Register(svr *grpc.Server) {
	pb.RegisterGreeterServer(svr, s)
}

func (s *Server) RegisterGateway(ctx context.Context,
	mux *runtime.ServeMux,
	endpoint string,
	dialOpts []grpc.DialOption) {
	pb.RegisterGreeterHandlerFromEndpoint(ctx, mux, endpoint, dialOpts)
}

func (s *Server) SayHello(_ context.Context, in *pb.HelloRequest) (*pb.HelloReply, error) {
	log.Printf("Received: %v", in.GetName())
	return &pb.HelloReply{Message: "Hello " + in.GetName()}, nil
}
