package service

import (
	"context"
	aiv1 "universal/api/ai/v1"
	gatewayv1 "universal/api/gateway/v1"
	"universal/app/gateway/internal/biz"
	"universal/app/gateway/internal/data"

	"github.com/go-kratos/kratos/v2/log"
)

// AiService AI服务代理
type AiService struct {
	gatewayv1.UnimplementedAiServer

	data *data.Data
	uc   *biz.UserUsecase // 可以用于权限验证等
	log  *log.Helper
}

// NewAiService 创建AI服务
func NewAiService(data *data.Data, uc *biz.UserUsecase, logger log.Logger) *AiService {
	return &AiService{
		data: data,
		uc:   uc,
		log:  log.NewHelper(logger),
	}
}

// GetConversationAnalytics 获取对话统计分析
func (s *AiService) GetConversationAnalytics(ctx context.Context, req *aiv1.GetConversationAnalyticsRequest) (*aiv1.GetConversationAnalyticsReply, error) {
	s.log.WithContext(ctx).Infof("GetConversationAnalytics called")

	// 转发到AI服务的gRPC客户端
	return s.data.AiClient().GetConversationAnalytics(ctx, req)
}

// GetUserUsageStats 获取用户使用统计
func (s *AiService) GetUserUsageStats(ctx context.Context, req *aiv1.GetUserUsageStatsRequest) (*aiv1.GetUserUsageStatsReply, error) {
	s.log.WithContext(ctx).Infof("GetUserUsageStats called")

	// 转发到AI服务的gRPC客户端
	return s.data.AiClient().GetUserUsageStats(ctx, req)
}

// GetModelPerformanceStats 获取模型性能统计
func (s *AiService) GetModelPerformanceStats(ctx context.Context, req *aiv1.GetModelPerformanceStatsRequest) (*aiv1.GetModelPerformanceStatsReply, error) {
	s.log.WithContext(ctx).Infof("GetModelPerformanceStats called")

	// 转发到AI服务的gRPC客户端
	return s.data.AiClient().GetModelPerformanceStats(ctx, req)
}

// GetConversationTrends 获取对话趋势分析
func (s *AiService) GetConversationTrends(ctx context.Context, req *aiv1.GetConversationTrendsRequest) (*aiv1.GetConversationTrendsReply, error) {
	s.log.WithContext(ctx).Infof("GetConversationTrends called")

	// 转发到AI服务的gRPC客户端
	return s.data.AiClient().GetConversationTrends(ctx, req)
}

// GetTopicAnalysis 获取话题分析
func (s *AiService) GetTopicAnalysis(ctx context.Context, req *aiv1.GetTopicAnalysisRequest) (*aiv1.GetTopicAnalysisReply, error) {
	s.log.WithContext(ctx).Infof("GetTopicAnalysis called")

	// 转发到AI服务的gRPC客户端
	return s.data.AiClient().GetTopicAnalysis(ctx, req)
}

// GetSystemOverview 获取系统总览统计
func (s *AiService) GetSystemOverview(ctx context.Context, req *aiv1.GetSystemOverviewRequest) (*aiv1.GetSystemOverviewReply, error) {
	s.log.WithContext(ctx).Infof("GetSystemOverview called")

	// 转发到AI服务的gRPC客户端
	return s.data.AiClient().GetSystemOverview(ctx, req)
}
