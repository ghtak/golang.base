package grpcadapter

type Env struct {
	GrpcadapterAddress string `mapstructure:"GRPCADAPTER_ADDRESS"`
}

func (Env) Name() string {
	return moduleName
}

func NewEnv() Env {
	return Env{}
}
