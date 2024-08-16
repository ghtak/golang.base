package user

import (
	"context"
	pb "github.com/ghtak/golang.grpc.base/gen/go"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

func NewService(logger *zap.Logger) *Server {
	return &Server{
		logger: logger,
	}
}

type Server struct {
	pb.UnimplementedUserServer
	logger *zap.Logger
}

func (s *Server) Register(svr *grpc.Server) {
	pb.RegisterUserServer(svr, s)
}

func (s *Server) RegisterGateway(ctx context.Context,
	mux *runtime.ServeMux,
	endpoint string,
	dialOpts []grpc.DialOption) {
	pb.RegisterUserHandlerFromEndpoint(ctx, mux, endpoint, dialOpts)
}

func (s *Server) GetUser(_ context.Context, in *pb.UserRequest) (*pb.UserResponse, error) {
	s.logger.Info("Received", zap.String("name", in.GetName()))
	return &pb.UserResponse{Id: "0"}, nil
}
