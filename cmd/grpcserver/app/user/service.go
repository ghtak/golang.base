package user

import (
	"context"
	pb "github.com/ghtak/golang.grpc.base/gen/go"
	"github.com/ghtak/golang.grpc.base/internal/adapter/grpcfx"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

func NewService(logger *zap.Logger) grpcfx.Service {
	return &server{
		logger: logger,
	}
}

type server struct {
	pb.UnimplementedUserServer
	logger *zap.Logger
}

func (s *server) Register(svr *grpc.Server) {
	pb.RegisterUserServer(svr, s)
}

func (s *server) GetUser(_ context.Context, in *pb.UserRequest) (*pb.UserResponse, error) {
	s.logger.Info("Received", zap.String("name", in.GetName()))
	return &pb.UserResponse{Id: "0"}, nil
}
