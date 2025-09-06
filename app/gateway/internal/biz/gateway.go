package biz

import (
	"context"
	"time"
	pb "universal/api/gateway/v1"

	"github.com/go-kratos/kratos/v2/log"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// GatewayRepo 网关仓储接口
type GatewayRepo interface {
	// 获取后端服务状态
	GetServiceStatuses(ctx context.Context) ([]*pb.ServiceStatus, error)
	// 获取网关启动时间
	GetStartTime(ctx context.Context) time.Time
}

// GatewayUsecase 网关用例
type GatewayUsecase struct {
	repo GatewayRepo
	log  *log.Helper
}

// NewGatewayUsecase 创建网关用例
func NewGatewayUsecase(repo GatewayRepo, logger log.Logger) *GatewayUsecase {
	return &GatewayUsecase{
		repo: repo,
		log:  log.NewHelper(logger),
	}
}

// GetGatewayInfo 获取网关信息
func (uc *GatewayUsecase) GetGatewayInfo(ctx context.Context) (*pb.GatewayInfo, error) {
	startTime := uc.repo.GetStartTime(ctx)
	uptime := time.Since(startTime)

	info := &pb.GatewayInfo{
		Name:          "Universal Gateway",
		Version:       "1.0.0",
		Description:   "Universal microservice gateway providing unified access to all backend services",
		StartTime:     timestamppb.New(startTime),
		UptimeSeconds: int64(uptime.Seconds()),
		SupportedApis: []string{
			"/api/user/v1",             // 用户服务API
			"/api/ai/v1",               // AI服务API (统计分析)
			"/api/ai/v1/conversations", // 对话服务API
			"/api/ai/v1/knowledge",     // 知识库服务API
			"/api/ai/v1/tools",         // 工具服务API
			"/api/gateway/v1",          // 网关管理API
		},
	}

	return info, nil
}

// GetGatewayHealth 获取网关健康状态
func (uc *GatewayUsecase) GetGatewayHealth(ctx context.Context) (*pb.GetGatewayHealthReply, error) {
	// 获取后端服务状态
	services, err := uc.repo.GetServiceStatuses(ctx)
	if err != nil {
		uc.log.Errorf("Failed to get service statuses: %v", err)
		return &pb.GetGatewayHealthReply{
			Status:    "unhealthy",
			Message:   "Failed to check backend services",
			CheckTime: timestamppb.Now(),
			Services:  []*pb.ServiceStatus{},
		}, nil
	}

	// 计算整体健康状态
	status := "healthy"
	message := "All services are healthy"
	unhealthyCount := 0

	for _, service := range services {
		if service.Status != "healthy" {
			unhealthyCount++
		}
	}

	if unhealthyCount > 0 {
		if unhealthyCount == len(services) {
			status = "unhealthy"
			message = "All backend services are unhealthy"
		} else {
			status = "degraded"
			message = "Some backend services are unhealthy"
		}
	}

	return &pb.GetGatewayHealthReply{
		Status:    status,
		Message:   message,
		CheckTime: timestamppb.Now(),
		Services:  services,
	}, nil
}
