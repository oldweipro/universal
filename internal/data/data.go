package data

import (
	"universal/internal/conf"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-redis/redis/v8"
	"github.com/google/wire"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	gormlogger "gorm.io/gorm/logger"
)

// ProviderSet is data providers.
var ProviderSet = wire.NewSet(NewData, NewGreeterRepo)

// Data .
type Data struct {
	db  *gorm.DB
	rdb *redis.Client
}

// NewData .
func NewData(c *conf.Data, logger log.Logger) (*Data, func(), error) {
	helper := log.NewHelper(logger)
	// 初始化数据库连接
	db, err := gorm.Open(mysql.Open(c.Database.Source), &gorm.Config{
		Logger: gormlogger.Default.LogMode(gormlogger.Info),
	})
	if err != nil {
		helper.Fatalf("failed to connect database: %v", err)
		return nil, nil, err
	}
	// 数据库迁移
	//if err = db.AutoMigrate(&User{}, &APIKey{}, &File{}); err != nil {
	//	helper.Fatalf("failed to migrate database: %v", err)
	//	return nil, nil, err
	//}

	// 初始化Redis连接
	rdb := redis.NewClient(&redis.Options{
		Addr:         c.Redis.Addr,
		ReadTimeout:  c.Redis.ReadTimeout.AsDuration(),
		WriteTimeout: c.Redis.WriteTimeout.AsDuration(),
	})

	// TODO 初始化MinIO客户端

	cleanup := func() {
		helper.Info("closing the data resources")
	}
	d := &Data{
		db:  db,
		rdb: rdb,
	}
	return d, cleanup, nil
}
