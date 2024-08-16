package ginfx

import (
	"fmt"
	"github.com/ghtak/golang.grpc.base/internal/core"
)

func NewEnv(env core.Env) Env {
	return Env{
		Address: env.GetString("gin.address", "0.0.0.0"),
		Port:    env.GetInt("gin.port", 9999),
	}
}

type Env struct {
	Address string
	Port    int
}

func (e Env) ServerAddress() string {
	return fmt.Sprintf("%s:%d", e.Address, e.Port)
}
