package core

import (
	"fmt"
	flag "github.com/spf13/pflag"
	"github.com/spf13/viper"
	"go.uber.org/fx"
	"log"
)

type EnvRepository struct {
	Envs map[string]interface{}
}

type NamedEnv interface {
	Name() string
}

func AsNamedEnv(i interface{}) interface{} {
	return fx.Annotate(i, fx.As(new(NamedEnv)), fx.ResultTags(`group:"NamedEnv"`))
}

func NewEnvRepository(envs []NamedEnv) *EnvRepository {
	envRepo := &EnvRepository{
		Envs: make(map[string]interface{}),
	}
	envFile := flag.String("cfg", ".env", "set env filename")
	flag.Parse()
	viper.SetConfigFile(*envFile)
	viper.SetConfigType("env")
	if err := viper.ReadInConfig(); err != nil {
		log.Fatal(fmt.Sprintf("Can't find file %s", *envFile), err)
	}
	for _, env := range envs {
		if err := viper.Unmarshal(&env); err != nil {
			log.Fatal(fmt.Sprintf("Can't Load Env %s", *envFile), err)
		}
		envRepo.Envs[env.Name()] = env
	}
	return envRepo
}

type Env struct {
	Profile    string `mapstructure:"PROFILE"`
	LogLevel   string `mapstructure:"LOG_LEVEL"`
	LogEncoder string `mapstructure:"LOG_ENCODER"`
}

func (Env) Name() string {
	return moduleName
}

func NewEnv() Env {
	return Env{}
}
