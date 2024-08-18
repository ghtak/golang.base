package gormfx

import (
	"github.com/ghtak/golang.grpc.base/internal/core"
	"go.uber.org/zap"
)

type envDatabase struct {
	ConnectionConfigs []envConnectionConfigsItem `mapstructure:"connection_configs"`
	Databases         []envDatabasesItem         `mapstructure:"databases"`
}

type envConnectionConfigsItem struct {
	Name            string `mapstructure:"name"`
	MaxOpenConns    int    `mapstructure:"max_open_conns"`
	MaxIdleConns    int    `mapstructure:"max_idle_conns"`
	ConnMaxIdleTime int    `mapstructure:"conn_max_idle_time"`
	ConnMaxLifeTime int    `mapstructure:"conn_max_life_time"`
}

type envDatabaseConnectionInfo struct {
	Driver string `mapstructure:"driver"`
	Dsn    string `mapstructure:"dsn"`
}

type envPrimary struct {
	Config string `mapstructure:"connection_config"`
	Driver string `mapstructure:"driver"`
	Dsn    string `mapstructure:"dsn"`
}

type envSecondary struct {
	Config   string                      `mapstructure:"connection_config"`
	Sources  []envDatabaseConnectionInfo `mapstructure:"sources"`
	Replicas []envDatabaseConnectionInfo `mapstructure:"replicas"`
}

type envDatabasesItem struct {
	Name      string       `mapstructure:"name"`
	Primary   envPrimary   `mapstructure:"primary"`
	Secondary envSecondary `mapstructure:"secondary"`
}

type Env struct {
	envDatabase
	ConnectionConfigsMap map[string]envConnectionConfigsItem
}

func NewEnv(env core.Env, logger *zap.Logger) Env {
	envDB := envDatabase{}
	err := env.Viper.Sub("database").Unmarshal(&envDB)
	if err != nil {
		logger.Error("database conf fail", zap.Error(err))
	}
	connCfgMap := map[string]envConnectionConfigsItem{}
	for _, c := range envDB.ConnectionConfigs {
		connCfgMap[c.Name] = c
	}
	for _, db := range envDB.Databases {
		_, ok := connCfgMap[db.Primary.Config]
		if ok == false {
			logger.Error("database conf connection_config not exist",
				zap.String("database.primary", db.Name),
				zap.String("database.config", db.Primary.Config))
		}
		_, ok = connCfgMap[db.Secondary.Config]
		if ok == false {
			logger.Error("database conf connection_config not exist",
				zap.String("database.primary", db.Name),
				zap.String("database.config", db.Secondary.Config))
		}
	}
	return Env{
		envDatabase:          envDB,
		ConnectionConfigsMap: connCfgMap,
	}
}
