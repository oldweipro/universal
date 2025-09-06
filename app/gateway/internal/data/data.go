package data

import (
	"context"
	aiv1 "universal/api/ai/v1"
	userv1 "universal/api/user/v1"
	"universal/app/gateway/internal/conf"

	"github.com/go-kratos/kratos/contrib/registry/consul/v2"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware/recovery"
	"github.com/go-kratos/kratos/v2/registry"
	"github.com/go-kratos/kratos/v2/transport/grpc"
	"github.com/go-redis/redis/v8"
	"github.com/google/wire"
	consulAPI "github.com/hashicorp/consul/api"
)

// ProviderSet is data providers.
var ProviderSet = wire.NewSet(
	NewData, NewGreeterRepo, NewUserRepo, NewGatewayRepo,
	NewDiscovery,
	NewRegistrar,
	NewUserServiceClient,
	NewAiServiceClient,
	NewConversationServiceClient,
	NewModelServiceClient,
	NewKnowledgeServiceClient,
	NewToolServiceClient,
)

// Data .
type Data struct {
	rdb *redis.Client
	uc  userv1.UserClient
	aic aiv1.AiClient
	cc  aiv1.ConversationClient
	mc  aiv1.ModelClient
	kc  aiv1.KnowledgeClient
	tc  aiv1.ToolClient
}

// NewData .
func NewData(
	c *conf.Data,
	logger log.Logger,
	uc userv1.UserClient,
	aic aiv1.AiClient,
	cc aiv1.ConversationClient,
	mc aiv1.ModelClient,
	kc aiv1.KnowledgeClient,
	tc aiv1.ToolClient,
) (*Data, func(), error) {
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
		uc:  uc,
		aic: aic,
		cc:  cc,
		mc:  mc,
		kc:  kc,
		tc:  tc,
	}

	//cleanup := func() {
	//	log.NewHelper(logger).Info("closing the data resources")
	//}
	return d, cleanup, nil
}

// 客户端访问方法
func (d *Data) UserClient() userv1.UserClient {
	return d.uc
}

func (d *Data) AiClient() aiv1.AiClient {
	return d.aic
}

func (d *Data) ConversationClient() aiv1.ConversationClient {
	return d.cc
}

func (d *Data) ModelClient() aiv1.ModelClient {
	return d.mc
}

func (d *Data) KnowledgeClient() aiv1.KnowledgeClient {
	return d.kc
}

func (d *Data) ToolClient() aiv1.ToolClient {
	return d.tc
}

func NewDiscovery(conf *conf.Registry) registry.Discovery {
	c := consulAPI.DefaultConfig()
	c.Address = conf.Consul.Address
	c.Scheme = conf.Consul.Scheme
	cli, err := consulAPI.NewClient(c)
	if err != nil {
		panic(err)
	}
	r := consul.New(cli, consul.WithHealthCheck(false))
	return r
}

func NewRegistrar(conf *conf.Registry) registry.Registrar {
	c := consulAPI.DefaultConfig()
	c.Address = conf.Consul.Address
	c.Scheme = conf.Consul.Scheme
	cli, err := consulAPI.NewClient(c)
	if err != nil {
		panic(err)
	}
	r := consul.New(cli, consul.WithHealthCheck(false))
	return r
}

func NewUserServiceClient(r registry.Discovery) userv1.UserClient {
	conn, err := grpc.DialInsecure(
		context.Background(),
		grpc.WithEndpoint("discovery:///universal.user.service"),
		grpc.WithDiscovery(r),
		grpc.WithMiddleware(
			recovery.Recovery(),
		),
	)
	if err != nil {
		panic(err)
	}
	c := userv1.NewUserClient(conn)
	return c
}

// AI服务客户端
func NewAiServiceClient(r registry.Discovery) aiv1.AiClient {
	conn, err := grpc.DialInsecure(
		context.Background(),
		grpc.WithEndpoint("discovery:///universal.ai.service"),
		grpc.WithDiscovery(r),
		grpc.WithMiddleware(
			recovery.Recovery(),
		),
	)
	if err != nil {
		panic(err)
	}
	return aiv1.NewAiClient(conn)
}

// 对话服务客户端
func NewConversationServiceClient(r registry.Discovery) aiv1.ConversationClient {
	conn, err := grpc.DialInsecure(
		context.Background(),
		grpc.WithEndpoint("discovery:///universal.ai.service"),
		grpc.WithDiscovery(r),
		grpc.WithMiddleware(
			recovery.Recovery(),
		),
	)
	if err != nil {
		panic(err)
	}
	return aiv1.NewConversationClient(conn)
}

// 模型服务客户端
func NewModelServiceClient(r registry.Discovery) aiv1.ModelClient {
	conn, err := grpc.DialInsecure(
		context.Background(),
		grpc.WithEndpoint("discovery:///universal.ai.service"),
		grpc.WithDiscovery(r),
		grpc.WithMiddleware(
			recovery.Recovery(),
		),
	)
	if err != nil {
		panic(err)
	}
	return aiv1.NewModelClient(conn)
}

// 知识库服务客户端
func NewKnowledgeServiceClient(r registry.Discovery) aiv1.KnowledgeClient {
	conn, err := grpc.DialInsecure(
		context.Background(),
		grpc.WithEndpoint("discovery:///universal.ai.service"),
		grpc.WithDiscovery(r),
		grpc.WithMiddleware(
			recovery.Recovery(),
		),
	)
	if err != nil {
		panic(err)
	}
	return aiv1.NewKnowledgeClient(conn)
}

// 工具服务客户端
func NewToolServiceClient(r registry.Discovery) aiv1.ToolClient {
	conn, err := grpc.DialInsecure(
		context.Background(),
		grpc.WithEndpoint("discovery:///universal.ai.service"),
		grpc.WithDiscovery(r),
		grpc.WithMiddleware(
			recovery.Recovery(),
		),
	)
	if err != nil {
		panic(err)
	}
	return aiv1.NewToolClient(conn)
}
