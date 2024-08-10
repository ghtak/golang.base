package main

import (
	"context"
	"flag"
	"fmt"
	pb "github.com/ghtak/golang.grpc.base/gen"
	"github.com/ghtak/golang.grpc.base/internal/core"
	"go.uber.org/fx"
	"google.golang.org/grpc"
	"log"
	"net"
)

func NewGrpcServer(lx fx.Lifecycle, listener *net.Listener) *grpc.Server {
	s := grpc.NewServer()
	lx.Append(
		fx.Hook{
			OnStart: func(ctx context.Context) error {
				go func() {
					err := s.Serve(*listener)
					if err != nil {
						log.Fatalf("failed to serve: %v", err)
					}
				}()
				return nil
			},
			OnStop: func(ctx context.Context) error {
				return nil
			},
		})
	return s
}

func NewListener() *net.Listener {
	port := flag.Int("port", 50051, "The server port")
	flag.Parse()
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	return &lis
}

type GreeterServerImpl struct {
	pb.UnimplementedGreeterServer
}

// SayHello implements helloworld.GreeterServer
func (s *GreeterServerImpl) SayHello(_ context.Context, in *pb.HelloRequest) (*pb.HelloReply, error) {
	log.Printf("Received: %v", in.GetName())
	return &pb.HelloReply{Message: "Hello " + in.GetName()}, nil
}

func (s *GreeterServerImpl) ServiceDesc() *grpc.ServiceDesc {
	return &pb.Greeter_ServiceDesc
}

func NewGreetService() *GreeterServerImpl {
	return &GreeterServerImpl{}
}

type InitCtx struct {
}

func RegisterServices(services []core.GrpcService, s *grpc.Server) *InitCtx {
	for _, svc := range services {
		if t, ok := svc.(interface{ testEmbeddedByValue() }); ok {
			t.testEmbeddedByValue()
		}
		s.RegisterService(svc.ServiceDesc(), svc)
	}
	return nil
}

func main() {
	fx.New(
		fx.Provide(
			NewListener,
			NewGrpcServer,
			core.AsGrpcService(NewGreetService),
			fx.Annotate(RegisterServices, fx.ParamTags(`group:"GrpcService"`))),
		fx.Invoke(func(*InitCtx) {}),
	).Run()
}
