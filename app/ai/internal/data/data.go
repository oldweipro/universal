package data

import (
	"universal/app/ai/internal/conf"
	"universal/app/ai/internal/data/model"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/google/wire"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	gormlogger "gorm.io/gorm/logger"
)

// ProviderSet is data providers.
var ProviderSet = wire.NewSet(NewData, NewAiRepo, NewProviderRepo, NewModelRepo, NewQuotaRepo, NewRateLimitRepo, NewHealthRepo, NewConversationRepo, NewKnowledgeRepo)

// Data .
type Data struct {
	db *gorm.DB
}

// GetDB returns the database instance
func (d *Data) GetDB() *gorm.DB {
	return d.db
}

// NewTestData creates a test Data instance
func NewTestData(db *gorm.DB) *Data {
	return &Data{db: db}
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
	if err = db.AutoMigrate(
		&model.Provider{},
		&model.Model{},
		&model.UserQuota{},
		&model.UsageStats{},
		&model.RateLimitConfig{},
		&model.ModelHealth{},
		&model.UserDefaultModel{},
		&model.Conversation{},
		&model.ConversationMemory{},
		&model.Message{},
		&model.ToolCall{},
		&model.MessageAttachment{},
		&model.KnowledgeBase{},
		&model.Document{},
		&model.KnowledgeChunk{},
		&model.ProcessingJob{},
		&model.SearchHistory{},
	); err != nil {
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
