package server

import (
	gatewayv1 "universal/api/gateway/v1"
	v1 "universal/api/helloworld/v1"
	"universal/app/gateway/internal/conf"
	"universal/app/gateway/internal/service"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware/recovery"
	"github.com/go-kratos/kratos/v2/transport/grpc"
)

// NewGRPCServer new a gRPC server.
func NewGRPCServer(c *conf.Server,
	greeter *service.GreeterService,
	gatewayService *service.GatewayService,
	userService *service.UserService,
	logger log.Logger,
) *grpc.Server {
	var opts = []grpc.ServerOption{
		grpc.Middleware(
			recovery.Recovery(),
		),
	}
	if c.Grpc.Network != "" {
		opts = append(opts, grpc.Network(c.Grpc.Network))
	}
	if c.Grpc.Addr != "" {
		opts = append(opts, grpc.Address(c.Grpc.Addr))
	}
	if c.Grpc.Timeout != nil {
		opts = append(opts, grpc.Timeout(c.Grpc.Timeout.AsDuration()))
	}
	srv := grpc.NewServer(opts...)
	v1.RegisterGreeterServer(srv, greeter)
	gatewayv1.RegisterGatewayServer(srv, gatewayService)
	gatewayv1.RegisterUserServer(srv, userService)
	return srv
}
