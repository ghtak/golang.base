package ginfx

import (
	"github.com/ghtak/golang.grpc.base/internal/core"
)

func NewEnv(env core.Env) Env {
	return Env{
		Address: env.GetString("gin.address", "0.0.0.0:9999"),
	}
}

type Env struct {
	Address string
}
