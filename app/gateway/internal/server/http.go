package server

import (
	gatewayv1 "universal/api/gateway/v1"
	v1 "universal/api/helloworld/v1"
	"universal/app/gateway/internal/conf"
	"universal/app/gateway/internal/service"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware/recovery"
	"github.com/go-kratos/kratos/v2/transport/http"
)

// NewHTTPServer new an HTTP server.
func NewHTTPServer(c *conf.Server,
	greeter *service.GreeterService,
	gatewayService *service.GatewayService,
	userService *service.UserService,
	aiService *service.AiService,
	conversationService *service.ConversationService,
	knowledgeService *service.KnowledgeService,
	toolService *service.ToolService,
	logger log.Logger,
) *http.Server {
	var opts = []http.ServerOption{
		http.Middleware(
			recovery.Recovery(),
		),
	}
	if c.Http.Network != "" {
		opts = append(opts, http.Network(c.Http.Network))
	}
	if c.Http.Addr != "" {
		opts = append(opts, http.Address(c.Http.Addr))
	}
	if c.Http.Timeout != nil {
		opts = append(opts, http.Timeout(c.Http.Timeout.AsDuration()))
	}
	srv := http.NewServer(opts...)
	v1.RegisterGreeterHTTPServer(srv, greeter)
	gatewayv1.RegisterGatewayHTTPServer(srv, gatewayService)
	gatewayv1.RegisterUserHTTPServer(srv, userService)

	// 注册Gateway AI服务的HTTP接口 - 这是主要的HTTP API入口
	gatewayv1.RegisterAiHTTPServer(srv, aiService)
	gatewayv1.RegisterConversationHTTPServer(srv, conversationService)
	gatewayv1.RegisterKnowledgeHTTPServer(srv, knowledgeService)
	gatewayv1.RegisterToolHTTPServer(srv, toolService)

	return srv
}
