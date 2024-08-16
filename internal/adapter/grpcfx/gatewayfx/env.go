package gatewayfx

import (
	"github.com/ghtak/golang.grpc.base/internal/core"
)

func NewEnv(env core.Env) Env {
	return Env{
		Address:     env.GetString("grpc_gateway.address", "0.0.0.0:9999"),
		GrpcAddress: env.GetString("grpc_gateway.grpc_address", "0.0.0.0:9999"),
	}
}

type Env struct {
	Address     string
	GrpcAddress string
}
