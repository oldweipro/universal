package data

import (
	"universal/app/system/internal/conf"
	"universal/app/system/internal/data/model"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/google/wire"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	gormlogger "gorm.io/gorm/logger"
)

// ProviderSet is data providers.
var ProviderSet = wire.NewSet(NewData, NewUserRepo)

// Data .
type Data struct {
	db *gorm.DB
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
	if err = db.AutoMigrate(&model.User{}); err != nil {
		helper.Fatalf("failed to migrate database: %v", err)
		return nil, nil, err
	}
	cleanup := func() {
		helper.Info("closing the data resources")
	}
	d := &Data{
		db: db,
	}
	return d, cleanup, nil
}
