package main

import (
	"github.com/ghtak/golang.grpc.base/internal/adapter/gormfx"
	"github.com/ghtak/golang.grpc.base/internal/core"
	"go.uber.org/fx"
	"go.uber.org/zap"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type Product struct {
	gorm.Model
	Code  string
	Price uint
}

type MainDB struct {
	DB gormfx.DB
}

type MainDBResults struct {
	fx.Out
	MainDB MainDB
}

func NewMainDB(env gormfx.Env, logger *zap.Logger) MainDBResults {
	return MainDBResults{
		MainDB: MainDB{DB: gormfx.NewDB(env.DbConnInfos["main"], logger)},
	}
}

func main() {
	fx.New(
		core.Module,
		gormfx.Module,
		fx.Provide(NewMainDB),
		fx.Invoke(
			func(p MainDB) {
				// 테이블 자동 생성
				mainDB := p.DB.Replicas.Next()
				mainDB.AutoMigrate(&Product{})

				// 생성
				mainDB.Create(&Product{Code: "D42", Price: 100})

				// 읽기
				var product Product
				mainDB.First(&product, 2)                 // primary key기준으로 product 찾기
				mainDB.First(&product, "code = ?", "D42") // code가 D42인 product 찾기

				// 수정 - product의 price를 200으로
				mainDB.Model(&product).Update("Price", 200)
				// 수정 - 여러개의 필드를 수정하기
				mainDB.Model(&product).Updates(Product{Price: 200, Code: "F42"})
				mainDB.Model(&product).Updates(map[string]interface{}{"Price": 200, "Code": "F42"})

				// 삭제 - product 삭제하기
				mainDB.Delete(&product, 2)
			}),
	).Run()
}

func hello() {
	db, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
	if err != nil {
		panic("Db 연결에 실패하였습니다.")
	}
	sqlDb, _ := db.DB()
	sqlDb.SetMaxIdleConns(10)
	sqlDb.SetMaxOpenConns(100)

	// 테이블 자동 생성
	db.AutoMigrate(&Product{})

	// 생성
	db.Create(&Product{Code: "D42", Price: 100})

	// 읽기
	var product Product
	db.First(&product, 1)                 // primary key기준으로 product 찾기
	db.First(&product, "code = ?", "D42") // code가 D42인 product 찾기

	// 수정 - product의 price를 200으로
	db.Model(&product).Update("Price", 200)
	// 수정 - 여러개의 필드를 수정하기
	db.Model(&product).Updates(Product{Price: 200, Code: "F42"})
	db.Model(&product).Updates(map[string]interface{}{"Price": 200, "Code": "F42"})

	// 삭제 - product 삭제하기
	//db.Delete(&product, 1)
}
