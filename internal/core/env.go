package core

import (
	"fmt"
	flag "github.com/spf13/pflag"
	"github.com/spf13/viper"
	"go.uber.org/fx"
	"log"
)

type ModuleEnv interface {
	ModuleName() string
}

func AsModuleEnv(i interface{}) interface{} {
	return fx.Annotate(i, fx.As(new(ModuleEnv)), fx.ResultTags(`group:"ModuleEnv"`))
}

type Env struct {
	CoreProfile string `mapstructure:"CORE_PROFILE"`
	ModuleEnvs  map[string]ModuleEnv
}

func NewEnv(envs []ModuleEnv, loader EnvLoader) *Env {
	moduleEnvs := make(map[string]ModuleEnv)
	for _, e := range envs {
		moduleEnvs[e.ModuleName()] = e
	}
	env := &Env{ModuleEnvs: moduleEnvs}
	loader.Unmarshal(env)
	return env
}

type EnvLoader interface {
	Unmarshal(env any)
}

type envLoader struct {
	envFile *string
}

func (e envLoader) Unmarshal(env any) {
	if err := viper.Unmarshal(env); err != nil {
		log.Fatal(fmt.Sprintf("Can't Load Env %s", *e.envFile), err)
	}
}

func NewEnvLoader() EnvLoader {
	envFile := flag.String("cfg", ".env", "set env filename")
	flag.Parse()
	viper.SetConfigFile(*envFile)
	viper.SetConfigType("env")
	if err := viper.ReadInConfig(); err != nil {
		log.Fatal(fmt.Sprintf("Can't find file %s", *envFile), err)
	}
	return envLoader{envFile}
}
