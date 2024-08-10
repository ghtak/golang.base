package main

import (
	"context"
	pb "github.com/ghtak/golang.grpc.base/gen"
	"github.com/ghtak/golang.grpc.base/internal/core"
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

func (s *GreeterServerImpl) ServiceDesc() *grpc.ServiceDesc {
	return &pb.Greeter_ServiceDesc
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

func (s *UserServerImpl) ServiceDesc() *grpc.ServiceDesc {
	return &pb.User_ServiceDesc
}

func NewUserService() *UserServerImpl {
	return &UserServerImpl{}
}

func main() {
	fx.New(core.Module,
		fx.Provide(
			core.AsGrpcService(NewGreetService),
			core.AsGrpcService(NewUserService),
		),
	).Run()
}
