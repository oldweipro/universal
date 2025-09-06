package data

import (
	"context"
	"time"
	pb "universal/api/gateway/v1"
	"universal/app/gateway/internal/biz"

	"github.com/go-kratos/kratos/v2/log"
)

var gatewayStartTime = time.Now()

type gatewayRepo struct {
	data *Data
	log  *log.Helper
}

// NewGatewayRepo 创建网关仓储
func NewGatewayRepo(data *Data, logger log.Logger) biz.GatewayRepo {
	return &gatewayRepo{
		data: data,
		log:  log.NewHelper(logger),
	}
}

// GetServiceStatuses 获取后端服务状态
func (r *gatewayRepo) GetServiceStatuses(ctx context.Context) ([]*pb.ServiceStatus, error) {
	services := make([]*pb.ServiceStatus, 0)

	// 检查用户服务状态
	userStatus := r.checkUserServiceHealth(ctx)
	services = append(services, userStatus)

	// 检查AI服务状态
	aiStatus := r.checkAiServiceHealth(ctx)
	services = append(services, aiStatus)

	return services, nil
}

// GetStartTime 获取网关启动时间
func (r *gatewayRepo) GetStartTime(ctx context.Context) time.Time {
	return gatewayStartTime
}

// 检查用户服务健康状态
func (r *gatewayRepo) checkUserServiceHealth(ctx context.Context) *pb.ServiceStatus {
	start := time.Now()
	status := &pb.ServiceStatus{
		Name:     "user-service",
		Endpoint: "universal.user.service",
	}

	// 这里应该调用用户服务的健康检查接口，暂时模拟
	// 在实际实现中，可以调用用户服务的健康检查端点
	defer func() {
		status.ResponseTimeMs = int32(time.Since(start).Milliseconds())
	}()

	// 尝试连接用户服务
	if r.data.uc != nil {
		status.Status = "healthy"
		status.Message = "User service is responding"
	} else {
		status.Status = "unhealthy"
		status.Message = "User service client not available"
	}

	return status
}

// 检查AI服务健康状态
func (r *gatewayRepo) checkAiServiceHealth(ctx context.Context) *pb.ServiceStatus {
	start := time.Now()
	status := &pb.ServiceStatus{
		Name:     "ai-service",
		Endpoint: "universal.ai.service",
	}

	defer func() {
		status.ResponseTimeMs = int32(time.Since(start).Milliseconds())
	}()

	// 检查AI服务各个客户端的可用性
	healthyCount := 0
	totalClients := 5

	if r.data.aic != nil {
		healthyCount++
	}
	if r.data.cc != nil {
		healthyCount++
	}
	if r.data.mc != nil {
		healthyCount++
	}
	if r.data.kc != nil {
		healthyCount++
	}
	if r.data.tc != nil {
		healthyCount++
	}

	if healthyCount == totalClients {
		status.Status = "healthy"
		status.Message = "AI service is fully available"
	} else if healthyCount > 0 {
		status.Status = "degraded"
		status.Message = "AI service is partially available"
	} else {
		status.Status = "unhealthy"
		status.Message = "AI service is not available"
	}

	return status
}
