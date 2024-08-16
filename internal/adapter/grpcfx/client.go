package grpcfx

import "google.golang.org/grpc"

func NewClient() *grpc.ClientConn {
	conn, _ := grpc.NewClient("")
	return conn
}
