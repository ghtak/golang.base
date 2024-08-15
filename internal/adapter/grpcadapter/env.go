package grpcadapter

import (
	"github.com/ghtak/golang.grpc.base/internal/core"
)

type Env struct {
	GrpcadapterAddress string `mapstructure:"GRPCADAPTER_ADDRESS"`
}

func (e Env) ModuleName() string {
	return moduleName
}

func NewEnv(loader core.EnvLoader) Env {
	env := Env{}
	loader.Unmarshal(&env)
	return env
}
