package fiberfx

import "github.com/ghtak/golang.grpc.base/internal/core"

type Env struct {
	Address string
}

func NewEnv(env core.Env) Env {
	return Env{
		Address: env.GetString("fiber.address", "0.0.0.0:9999"),
	}
}
