package gormfx

import (
	"github.com/ghtak/golang.grpc.base/internal/core"
	"go.uber.org/zap"
)

func newDataSources(srcs []envDatasource, cfgs map[string]ConnCfg) []DataSource {
	sources := make([]DataSource, len(srcs))
	for i := 0; i < len(srcs); i++ {
		sources[i] = DataSource{
			Driver: srcs[i].Driver,
			Dsn:    srcs[i].Dsn,
			Cfg:    cfgs[srcs[i].ConnCfg],
		}
	}
	return sources
}

func NewEnv(env core.Env, logger *zap.Logger) Env {
	egorm := envGorm{}
	err := env.Viper.Sub("gorm").Unmarshal(&egorm)
	if err != nil {
		logger.Error("gorm conf fail", zap.Error(err))
		return Env{}
	}
	cfgMap := map[string]ConnCfg{}
	for _, c := range egorm.ConnCfgs {
		cfgMap[c.Name] = c
	}
	infos := map[string]DbConnInfo{}
	for _, envConnInfo := range egorm.DbConnInfos {
		infos[envConnInfo.Name] = DbConnInfo{
			Name:     envConnInfo.Name,
			Sources:  newDataSources(envConnInfo.Sources, cfgMap),
			Replicas: newDataSources(envConnInfo.Replicas, cfgMap),
		}
	}
	return Env{
		DbConnInfos: infos,
	}
}

type DataSource struct {
	Driver string
	Dsn    string
	Cfg    ConnCfg
}

type DbConnInfo struct {
	Name     string
	Sources  []DataSource
	Replicas []DataSource
}

type Env struct {
	DbConnInfos map[string]DbConnInfo
}

type ConnCfg struct {
	Name            string `mapstructure:"name"`
	MaxOpenConns    int    `mapstructure:"max_open_conns"`
	MaxIdleConns    int    `mapstructure:"max_idle_conns"`
	ConnMaxIdleTime int    `mapstructure:"conn_max_idle_time"`
	ConnMaxLifeTime int    `mapstructure:"conn_max_life_time"`
}

type envDatasource struct {
	Driver  string `mapstructure:"driver"`
	Dsn     string `mapstructure:"dsn"`
	ConnCfg string `mapstructure:"conn_cfg"`
}

type envDbConnInfo struct {
	Name     string          `mapstructure:"name"`
	Sources  []envDatasource `mapstructure:"sources"`
	Replicas []envDatasource `mapstructure:"replicas"`
}

type envGorm struct {
	ConnCfgs    []ConnCfg       `mapstructure:"conn_cfgs"`
	DbConnInfos []envDbConnInfo `mapstructure:"db_conn_infos"`
}
