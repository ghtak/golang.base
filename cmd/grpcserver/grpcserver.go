package main

import (
	"github.com/ghtak/golang.grpc.base/internal/adapter/grpcadapter"
	"go.uber.org/fx"
)

func main() {
	fx.New(
		grpcadapter.Module,
	).Run()
}
