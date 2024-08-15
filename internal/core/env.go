package core

import (
	"fmt"
	flag "github.com/spf13/pflag"
	"github.com/spf13/viper"
	"log"
)

func NewEnv() Env {
	cfgFile := flag.String("cfg-file", "dev.conf", "set config filename")
	cfgType := flag.String("cfg-type", "json", "set config type json, yaml. toml")
	flag.Parse()
	env := Env{viper.New()}
	env.viper.SetConfigFile(*cfgFile)
	env.viper.SetConfigType(*cfgType)
	if err := env.viper.ReadInConfig(); err != nil {
		log.Fatal(fmt.Sprintf("ReadInConfig Fail With %s %s", *cfgFile, *cfgType), err)
	}
	return env
}

type Env struct {
	viper *viper.Viper
}

func (e Env) GetString(key string, defValue string) string {
	if e.viper.IsSet(key) {
		return e.viper.GetString(key)
	}
	return defValue
}

func (e Env) GetInt(key string, defValue int) int {
	if e.viper.IsSet(key) {
		return e.viper.GetInt(key)
	}
	return defValue
}

func (e Env) Sub(envPrefix string) SubEnv {
	return SubEnv{
		Env:       e,
		envPrefix: envPrefix,
	}
}

type SubEnv struct {
	Env
	envPrefix string
}

func (e SubEnv) GetString(key string, defValue string) string {
	key = fmt.Sprintf("%s.%s", e.envPrefix, key)
	if e.viper.IsSet(key) {
		return e.viper.GetString(key)
	}
	return defValue
}

func (e SubEnv) GetInt(key string, defValue int) int {
	key = fmt.Sprintf("%s.%s", e.envPrefix, key)
	if e.viper.IsSet(key) {
		return e.viper.GetInt(key)
	}
	return defValue
}
