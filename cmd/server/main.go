package main

import (
	"context"
	pb "github.com/ghtak/golang.grpc.base/gen/go"
	"github.com/ghtak/golang.grpc.base/internal/application"
	"github.com/ghtak/golang.grpc.base/internal/core0"
	"github.com/ghtak/golang.grpc.base/internal/middleware"
	"go.uber.org/fx"
	"google.golang.org/grpc"
	"log"
)

type GreeterServerImpl struct {
	pb.UnimplementedGreeterServer
}

func (s *GreeterServerImpl) SayHello(_ context.Context, in *pb.HelloRequest) (*pb.HelloReply, error) {
	log.Printf("Received: %v", in.GetName())
	return &pb.HelloReply{Message: "Hello " + in.GetName()}, nil
}

func (s *GreeterServerImpl) Register(svr *grpc.Server) {
	pb.RegisterGreeterServer(svr, s)
}

func NewGreetService() *GreeterServerImpl {
	return &GreeterServerImpl{}
}

type UserServerImpl struct {
	pb.UnimplementedUserServer
}

func (s *UserServerImpl) GetUser(_ context.Context, in *pb.UserRequest) (*pb.UserResponse, error) {
	log.Printf("Received: %v", in.GetName())
	return &pb.UserResponse{Id: "0"}, nil
}

func (s *UserServerImpl) Register(svr *grpc.Server) {
	pb.RegisterUserServer(svr, s)
}

func NewUserService() *UserServerImpl {
	return &UserServerImpl{}
}

func main() {
	fx.New(
		core0.Module,
		middleware.Module,
		application.ModuleAuth,
		fx.Provide(
			core0.AsGrpcService(NewGreetService),
			core0.AsGrpcService(NewUserService),
		),
	).Run()
}
