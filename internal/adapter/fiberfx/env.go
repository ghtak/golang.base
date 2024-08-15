package fiberfx

import "github.com/ghtak/golang.grpc.base/internal/core"

type Env struct {
	Address string
	Port    int
}

func NewEnv(env core.Env) Env {
	return Env{
		Address: env.GetString("fiber.address", "0.0.0.0"),
		Port:    env.GetInt("fiber.port", 9999),
	}
}
