package data

import (
	"universal/internal/conf"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-redis/redis/v8"
	"github.com/google/wire"
)

// ProviderSet is data providers.
var ProviderSet = wire.NewSet(NewData, NewGreeterRepo)

// Data .
type Data struct {
	rdb *redis.Client
}

// NewData .
func NewData(c *conf.Data, logger log.Logger) (*Data, func(), error) {
	helper := log.NewHelper(logger)

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
		rdb: rdb,
	}

	//cleanup := func() {
	//	log.NewHelper(logger).Info("closing the data resources")
	//}
	return d, cleanup, nil
}
