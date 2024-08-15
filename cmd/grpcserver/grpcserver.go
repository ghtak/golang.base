package main

import (
	"fmt"
	"github.com/ghtak/golang.grpc.base/internal/adapter/grpcadapter"
	"github.com/ghtak/golang.grpc.base/internal/core"
	"go.uber.org/fx"
	"google.golang.org/grpc"
	"reflect"
)

func printRepo(envs *core.EnvRepository) {
	for _, i := range envs.Envs {
		value := reflect.ValueOf(i)
		for index := 0; index < value.NumField(); index++ {
			fmt.Printf("Field name  : %s\n", value.Type().Field(index).Name) // 1번
			fmt.Printf("Field type  : %v\n", value.Field(index).Type())      // 2번
			fmt.Printf("Field Value : %v\n\n", value.Field(index))
		}
	}
}

func main() {
	fx.New(
		core.Module,
		grpcadapter.Module,
		fx.Invoke(func(*grpc.Server) {}),
		fx.Invoke(printRepo),
	).Run()
}
