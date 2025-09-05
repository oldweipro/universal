package service

import (
	"context"

	"google.golang.org/protobuf/types/known/timestamppb"

	pb "universal/api/ai/v1"
	"universal/app/ai/internal/biz"
)

type ModelService struct {
	pb.UnimplementedModelServer
	modelUc *biz.ModelUsecase
}

func NewModelService(modelUc *biz.ModelUsecase) *ModelService {
	return &ModelService{
		modelUc: modelUc,
	}
}

func (s *ModelService) CreateProvider(ctx context.Context, req *pb.CreateProviderRequest) (*pb.CreateProviderReply, error) {
	provider, err := s.modelUc.CreateProvider(ctx, req)
	if err != nil {
		return nil, err
	}

	return &pb.CreateProviderReply{
		Provider: convertBizProviderToPb(provider),
	}, nil
}
func (s *ModelService) UpdateProvider(ctx context.Context, req *pb.UpdateProviderRequest) (*pb.UpdateProviderReply, error) {
	provider, err := s.modelUc.UpdateProvider(ctx, req)
	if err != nil {
		return nil, err
	}

	return &pb.UpdateProviderReply{
		Provider: convertBizProviderToPb(provider),
	}, nil
}
func (s *ModelService) DeleteProvider(ctx context.Context, req *pb.DeleteProviderRequest) (*pb.DeleteProviderReply, error) {
	err := s.modelUc.DeleteProvider(ctx, req.Id)
	if err != nil {
		return nil, err
	}

	return &pb.DeleteProviderReply{}, nil
}
func (s *ModelService) ListProviders(ctx context.Context, req *pb.ListProvidersRequest) (*pb.ListProvidersReply, error) {
	providers, total, err := s.modelUc.ListProviders(ctx, req.Status, req.Page, req.PageSize)
	if err != nil {
		return nil, err
	}

	pbProviders := make([]*pb.Provider, len(providers))
	for i, provider := range providers {
		pbProviders[i] = convertBizProviderToPb(provider)
	}

	return &pb.ListProvidersReply{
		Providers: pbProviders,
		Total:     total,
	}, nil
}
func (s *ModelService) TestProvider(ctx context.Context, req *pb.TestProviderRequest) (*pb.TestProviderReply, error) {
	isAvailable, responseTime, errorMessage, err := s.modelUc.TestProvider(ctx, req.Id, req.ApiKey)
	if err != nil {
		return nil, err
	}

	return &pb.TestProviderReply{
		IsAvailable:  isAvailable,
		ResponseTime: responseTime,
		ErrorMessage: errorMessage,
	}, nil
}
func (s *ModelService) CreateModel(ctx context.Context, req *pb.CreateModelRequest) (*pb.CreateModelReply, error) {
	model, err := s.modelUc.CreateModel(ctx, req)
	if err != nil {
		return nil, err
	}

	return &pb.CreateModelReply{
		Model: convertBizModelToPb(model),
	}, nil
}
func (s *ModelService) UpdateModel(ctx context.Context, req *pb.UpdateModelRequest) (*pb.UpdateModelReply, error) {
	model, err := s.modelUc.UpdateModel(ctx, req)
	if err != nil {
		return nil, err
	}

	return &pb.UpdateModelReply{
		Model: convertBizModelToPb(model),
	}, nil
}
func (s *ModelService) DeleteModel(ctx context.Context, req *pb.DeleteModelRequest) (*pb.DeleteModelReply, error) {
	err := s.modelUc.DeleteModel(ctx, req.Id)
	if err != nil {
		return nil, err
	}

	return &pb.DeleteModelReply{}, nil
}
func (s *ModelService) ListModels(ctx context.Context, req *pb.ListModelsRequest) (*pb.ListModelsReply, error) {
	models, total, err := s.modelUc.ListModels(ctx, req.ProviderId, req.Status, req.RequiredCapabilities, req.Page, req.PageSize)
	if err != nil {
		return nil, err
	}

	pbModels := make([]*pb.ModelInfo, len(models))
	for i, model := range models {
		pbModels[i] = convertBizModelToPb(model)
	}

	return &pb.ListModelsReply{
		Models: pbModels,
		Total:  total,
	}, nil
}
func (s *ModelService) GetModel(ctx context.Context, req *pb.GetModelRequest) (*pb.GetModelReply, error) {
	model, err := s.modelUc.GetModel(ctx, req.Id)
	if err != nil {
		return nil, err
	}

	return &pb.GetModelReply{
		Model: convertBizModelToPb(model),
	}, nil
}
func (s *ModelService) SwitchModel(ctx context.Context, req *pb.SwitchModelRequest) (*pb.SwitchModelReply, error) {
	err := s.modelUc.SwitchModel(ctx, req.UserId, req.ModelId)
	if err != nil {
		return nil, err
	}

	return &pb.SwitchModelReply{
		Success: true,
	}, nil
}
func (s *ModelService) GetUserQuota(ctx context.Context, req *pb.GetUserQuotaRequest) (*pb.GetUserQuotaReply, error) {
	quota, err := s.modelUc.GetUserQuota(ctx, req.UserId, req.ModelId)
	if err != nil {
		return nil, err
	}

	return &pb.GetUserQuotaReply{
		Quotas: []*pb.UserQuota{convertBizUserQuotaToPb(quota)},
	}, nil
}
func (s *ModelService) UpdateUserQuota(ctx context.Context, req *pb.UpdateUserQuotaRequest) (*pb.UpdateUserQuotaReply, error) {
	quota := &biz.UserQuota{
		UserID:              req.UserId,
		ModelID:             req.ModelId,
		DailyTokenLimit:     req.DailyTokenLimit,
		MonthlyTokenLimit:   req.MonthlyTokenLimit,
		DailyRequestLimit:   req.DailyRequestLimit,
		MonthlyRequestLimit: req.MonthlyRequestLimit,
	}

	updatedQuota, err := s.modelUc.UpdateUserQuota(ctx, quota)
	if err != nil {
		return nil, err
	}

	return &pb.UpdateUserQuotaReply{
		Quota: convertBizUserQuotaToPb(updatedQuota),
	}, nil
}
func (s *ModelService) GetUsageStats(ctx context.Context, req *pb.GetUsageStatsRequest) (*pb.GetUsageStatsReply, error) {
	stats, err := s.modelUc.GetUsageStats(ctx, req.UserId, req.ModelId, req.StartDate, req.EndDate, req.GroupBy)
	if err != nil {
		return nil, err
	}

	pbStats := make([]*pb.UsageStats, len(stats))
	totalRequests := int64(0)
	totalTokens := int64(0)
	totalCost := float64(0)

	for i, stat := range stats {
		pbStats[i] = convertBizUsageStatsToPb(stat)
		totalRequests += stat.TotalRequests
		totalTokens += stat.TotalTokens
		totalCost += stat.TotalCost
	}

	return &pb.GetUsageStatsReply{
		Stats:         pbStats,
		TotalRequests: totalRequests,
		TotalTokens:   totalTokens,
		TotalCost:     totalCost,
	}, nil
}
func (s *ModelService) ResetUsage(ctx context.Context, req *pb.ResetUsageRequest) (*pb.ResetUsageReply, error) {
	affectedUsers, err := s.modelUc.ResetUsage(ctx, req.UserId, req.ModelId, req.ResetType)
	if err != nil {
		return nil, err
	}

	return &pb.ResetUsageReply{
		Success:       true,
		AffectedUsers: affectedUsers,
	}, nil
}
func (s *ModelService) CheckRateLimit(ctx context.Context, req *pb.CheckRateLimitRequest) (*pb.CheckRateLimitReply, error) {
	allowed, reason, retryAfter, remainingRequests, remainingTokens, err := s.modelUc.CheckRateLimit(ctx, req.UserId, req.ModelId, req.RequestedTokens)
	if err != nil {
		return nil, err
	}

	return &pb.CheckRateLimitReply{
		Allowed:           allowed,
		Reason:            reason,
		RetryAfterSeconds: retryAfter,
		RemainingRequests: remainingRequests,
		RemainingTokens:   remainingTokens,
	}, nil
}
func (s *ModelService) GetRateLimitConfig(ctx context.Context, req *pb.GetRateLimitConfigRequest) (*pb.GetRateLimitConfigReply, error) {
	configs, err := s.modelUc.GetRateLimitConfig(ctx, req.ModelId, req.UserLevel)
	if err != nil {
		return nil, err
	}

	pbConfigs := make([]*pb.RateLimitConfig, len(configs))
	for i, config := range configs {
		pbConfigs[i] = convertBizRateLimitConfigToPb(config)
	}

	return &pb.GetRateLimitConfigReply{
		Configs: pbConfigs,
	}, nil
}
func (s *ModelService) UpdateRateLimitConfig(ctx context.Context, req *pb.UpdateRateLimitConfigRequest) (*pb.UpdateRateLimitConfigReply, error) {
	config := &biz.RateLimitConfig{
		ID:                 req.Id,
		ModelID:            req.ModelId,
		UserLevel:          req.UserLevel,
		RequestsPerMinute:  req.RequestsPerMinute,
		RequestsPerHour:    req.RequestsPerHour,
		RequestsPerDay:     req.RequestsPerDay,
		TokensPerMinute:    req.TokensPerMinute,
		ConcurrentRequests: req.ConcurrentRequests,
		BurstLimit:         req.BurstLimit,
	}

	updatedConfig, err := s.modelUc.UpdateRateLimitConfig(ctx, config)
	if err != nil {
		return nil, err
	}

	return &pb.UpdateRateLimitConfigReply{
		Config: convertBizRateLimitConfigToPb(updatedConfig),
	}, nil
}
func (s *ModelService) HealthCheck(ctx context.Context, req *pb.HealthCheckRequest) (*pb.HealthCheckReply, error) {
	health, err := s.modelUc.HealthCheck(ctx, req.ModelId)
	if err != nil {
		return nil, err
	}

	return &pb.HealthCheckReply{
		HealthStatus: []*pb.ModelHealth{convertBizModelHealthToPb(health)},
	}, nil
}
func (s *ModelService) GetModelMetrics(ctx context.Context, req *pb.GetModelMetricsRequest) (*pb.GetModelMetricsReply, error) {
	return s.modelUc.GetModelMetrics(ctx, req.ModelId, req.StartTime, req.EndTime)
}

// 类型转换函数

func convertBizProviderToPb(provider *biz.Provider) *pb.Provider {
	return &pb.Provider{
		Id:             provider.ID,
		Name:           provider.Name,
		DisplayName:    provider.DisplayName,
		Description:    provider.Description,
		ApiBaseUrl:     provider.APIBaseURL,
		DefaultApiKey:  provider.DefaultAPIKey,
		DefaultHeaders: provider.DefaultHeaders,
		Config:         provider.Config,
		Status:         provider.Status,
		CreatedAt:      timestamppb.New(provider.CreatedAt),
		UpdatedAt:      timestamppb.New(provider.UpdatedAt),
	}
}

func convertBizModelToPb(model *biz.Model) *pb.ModelInfo {
	return &pb.ModelInfo{
		Id:            model.ID,
		ProviderId:    model.ProviderID,
		Name:          model.Name,
		DisplayName:   model.DisplayName,
		Description:   model.Description,
		Version:       model.Version,
		Capabilities:  model.Capabilities,
		Limits:        model.Limits,
		Pricing:       model.Pricing,
		DefaultParams: model.DefaultParams,
		Status:        model.Status,
		CreatedAt:     timestamppb.New(model.CreatedAt),
		UpdatedAt:     timestamppb.New(model.UpdatedAt),
	}
}

func convertBizUserQuotaToPb(quota *biz.UserQuota) *pb.UserQuota {
	return &pb.UserQuota{
		UserId:              quota.UserID,
		ModelId:             quota.ModelID,
		DailyTokenLimit:     quota.DailyTokenLimit,
		MonthlyTokenLimit:   quota.MonthlyTokenLimit,
		DailyRequestLimit:   quota.DailyRequestLimit,
		MonthlyRequestLimit: quota.MonthlyRequestLimit,
		DailyTokensUsed:     quota.DailyTokensUsed,
		MonthlyTokensUsed:   quota.MonthlyTokensUsed,
		DailyRequestsUsed:   quota.DailyRequestsUsed,
		MonthlyRequestsUsed: quota.MonthlyRequestsUsed,
		ResetDailyAt:        timestamppb.New(quota.ResetDailyAt),
		ResetMonthlyAt:      timestamppb.New(quota.ResetMonthlyAt),
	}
}

func convertBizUsageStatsToPb(stats *biz.UsageStats) *pb.UsageStats {
	return &pb.UsageStats{
		UserId:             stats.UserID,
		ModelId:            stats.ModelID,
		Date:               stats.Date,
		TotalRequests:      stats.TotalRequests,
		SuccessfulRequests: stats.SuccessfulRequests,
		FailedRequests:     stats.FailedRequests,
		TotalTokens:        stats.TotalTokens,
		InputTokens:        stats.InputTokens,
		OutputTokens:       stats.OutputTokens,
		TotalCost:          stats.TotalCost,
		AvgResponseTime:    stats.AvgResponseTime,
	}
}

func convertBizRateLimitConfigToPb(config *biz.RateLimitConfig) *pb.RateLimitConfig {
	return &pb.RateLimitConfig{
		Id:                 config.ID,
		ModelId:            config.ModelID,
		UserLevel:          config.UserLevel,
		RequestsPerMinute:  config.RequestsPerMinute,
		RequestsPerHour:    config.RequestsPerHour,
		RequestsPerDay:     config.RequestsPerDay,
		TokensPerMinute:    config.TokensPerMinute,
		ConcurrentRequests: config.ConcurrentRequests,
		BurstLimit:         config.BurstLimit,
		CreatedAt:          timestamppb.New(config.CreatedAt),
		UpdatedAt:          timestamppb.New(config.UpdatedAt),
	}
}

func convertBizModelHealthToPb(health *biz.ModelHealth) *pb.ModelHealth {
	return &pb.ModelHealth{
		ModelId:        health.ModelID,
		IsHealthy:      health.IsHealthy,
		ResponseTime:   health.ResponseTime,
		SuccessRate:    health.SuccessRate,
		TotalRequests:  health.TotalRequests,
		FailedRequests: health.FailedRequests,
		ErrorMessage:   health.ErrorMessage,
		LastCheck:      timestamppb.New(health.LastCheck),
	}
}
