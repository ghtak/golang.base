package gormfx

import (
	"errors"
	"fmt"
	"github.com/ghtak/golang.grpc.base/internal/core"
	"go.uber.org/zap"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"time"
)

type DB struct {
	Name     string
	Sources  *core.RoundRobin[gorm.DB]
	Replicas *core.RoundRobin[gorm.DB]
}

func newDialector(source DataSource) (gorm.Dialector, error) {
	switch source.Driver {
	case "sqlite":
		return sqlite.Open(source.Dsn), nil
	case "postgres":
		//"host=localhost user=gorm password=gorm dbname=gorm port=9920 sslmode=disable TimeZone=Asia/Shanghai"
		return postgres.Open(source.Dsn), nil
	case "mysql":
		//"user:pass@tcp(127.0.0.1:3306)/dbname?charset=utf8mb4&parseTime=True&loc=Local"
		return mysql.Open(source.Dsn), nil
	default:
	}
	return nil, errors.New(fmt.Sprintf("%s are not supported", source.Driver))
}

func newDBSlice(sources []DataSource, logger *zap.Logger) []*gorm.DB {
	var retSlice []*gorm.DB
	for _, source := range sources {
		dialector, err := newDialector(source)
		if err != nil {
			logger.Error("newDialector() fail", zap.Error(err))
			continue
		}
		db, err := gorm.Open(dialector, &gorm.Config{})
		if err != nil {
			logger.Error("gorm.Open() fail", zap.Error(err))
			continue
		}
		sqlDB, err := db.DB()
		if err != nil {
			logger.Error("db.DB() fail", zap.Error(err))
			continue
		}
		sqlDB.SetMaxOpenConns(source.Cfg.MaxOpenConns)
		sqlDB.SetMaxIdleConns(source.Cfg.MaxIdleConns)
		sqlDB.SetConnMaxIdleTime(time.Hour * time.Duration(source.Cfg.ConnMaxIdleTime))
		sqlDB.SetConnMaxLifetime(time.Hour * time.Duration(source.Cfg.ConnMaxLifeTime))
		retSlice = append(retSlice, db)
	}
	return retSlice
}

func NewDB(dbConnInfo DbConnInfo, logger *zap.Logger) DB {
	return DB{
		Name:     dbConnInfo.Name,
		Sources:  core.NewRoundRobin(newDBSlice(dbConnInfo.Sources, logger)...),
		Replicas: core.NewRoundRobin(newDBSlice(dbConnInfo.Replicas, logger)...),
	}
}
