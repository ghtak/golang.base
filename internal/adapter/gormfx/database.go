package gormfx

import (
	"errors"
	"fmt"
	"go.uber.org/zap"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/plugin/dbresolver"
	"time"
)

var (
	driverSqlite   = "sqlite"
	driverMySQL    = "mysql"
	driverPostgres = "postgres"
)

type Database struct {
	dbMap map[string]*gorm.DB
}

func (d *Database) GetDatabase(name string) *gorm.DB {
	val, ok := d.dbMap[name]
	if ok {
		return val
	}
	return nil
}

func newDialector(driver, dsn string) (gorm.Dialector, error) {
	switch driver {
	case driverSqlite:
		return sqlite.Open(dsn), nil
	case driverPostgres:
		//"host=localhost user=gorm password=gorm dbname=gorm port=9920 sslmode=disable TimeZone=Asia/Shanghai"
		return postgres.Open(dsn), nil
	case driverMySQL:
		//"user:pass@tcp(127.0.0.1:3306)/dbname?charset=utf8mb4&parseTime=True&loc=Local"
		return mysql.Open(dsn), nil
	default:
	}
	return nil, errors.New(fmt.Sprintf("%s are not supported", driver))
}

func newDialectors(infos []envDatabaseConnectionInfo) []gorm.Dialector {
	var dials []gorm.Dialector
	for _, info := range infos {
		dial, err := newDialector(info.Driver, info.Dsn)
		if err != nil {
			continue
		}
		dials = append(dials, dial)
	}
	return dials
}

func newDatabase(env Env, database envDatabasesItem) (*gorm.DB, error) {
	dial, err := newDialector(database.Primary.Driver, database.Primary.Dsn)
	if err != nil {
		return nil, err
	}

	db, err := gorm.Open(dial, &gorm.Config{
		//Logger: logger.New(
		//	log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
		//	logger.Config{
		//		SlowThreshold:             time.Second, // Slow SQL threshold
		//		LogLevel:                  logger.Info, // Log level
		//		IgnoreRecordNotFoundError: true,        // Ignore ErrRecordNotFound error for logger
		//		ParameterizedQueries:      true,        // Don't include params in the SQL log
		//		Colorful:                  false,       // Disable color
		//	},
		//),
	})
	if err != nil {
		return db, err
	}
	sqlDB, err := db.DB()
	if err != nil {
		return db, err
	}
	connCfg := env.ConnectionConfigsMap[database.Primary.Config]
	sqlDB.SetMaxOpenConns(connCfg.MaxOpenConns)
	sqlDB.SetMaxIdleConns(connCfg.MaxIdleConns)
	sqlDB.SetConnMaxIdleTime(time.Hour * time.Duration(connCfg.ConnMaxIdleTime))
	sqlDB.SetConnMaxLifetime(time.Hour * time.Duration(connCfg.ConnMaxLifeTime))

	connCfg = env.ConnectionConfigsMap[database.Secondary.Config]
	db.Use(dbresolver.Register(
		dbresolver.Config{
			Sources:           newDialectors(database.Secondary.Sources),
			Replicas:          newDialectors(database.Secondary.Replicas),
			Policy:            dbresolver.RandomPolicy{},
			TraceResolverMode: true,
		}).
		SetMaxOpenConns(connCfg.MaxOpenConns).
		SetMaxIdleConns(connCfg.MaxIdleConns).
		SetConnMaxIdleTime(time.Hour * time.Duration(connCfg.ConnMaxIdleTime)).
		SetConnMaxLifetime(time.Hour * time.Duration(connCfg.ConnMaxLifeTime)))
	return db, nil
}

func NewDatabase(env Env, logger *zap.Logger) *Database {
	dbMap := map[string]*gorm.DB{}
	for _, database := range env.Databases {
		db, err := newDatabase(env, database)
		if err != nil {
			logger.Error("newDatabase Fail", zap.Any("config", database), zap.Error(err))
			continue
		}
		dbMap[database.Name] = db
	}
	return &Database{dbMap: dbMap}
}
