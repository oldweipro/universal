package data

import (
	"context"
	"time"

	pb "universal/api/ai/v1"
	"universal/app/ai/internal/biz"
	"universal/app/ai/internal/data/model"

	"github.com/go-kratos/kratos/v2/log"
	"gorm.io/gorm"
)

type providerRepo struct {
	data *Data
	log  *log.Helper
}

// NewProviderRepo 创建提供商仓库
func NewProviderRepo(data *Data, logger log.Logger) biz.ProviderRepo {
	return &providerRepo{
		data: data,
		log:  log.NewHelper(logger),
	}
}

func (r *providerRepo) CreateProvider(ctx context.Context, provider *biz.Provider) (*biz.Provider, error) {
	po := &model.Provider{
		Name:          provider.Name,
		DisplayName:   provider.DisplayName,
		Description:   provider.Description,
		APIBaseURL:    provider.APIBaseURL,
		DefaultAPIKey: provider.DefaultAPIKey,
		Status:        provider.Status,
	}
	po.SetDefaultHeaders(provider.DefaultHeaders)
	po.SetConfig(provider.Config)

	if err := r.data.db.WithContext(ctx).Create(po).Error; err != nil {
		return nil, err
	}

	return r.convertProviderModelToBiz(po), nil
}

func (r *providerRepo) UpdateProvider(ctx context.Context, provider *biz.Provider) (*biz.Provider, error) {
	po := &model.Provider{}
	if err := r.data.db.WithContext(ctx).First(po, provider.ID).Error; err != nil {
		return nil, err
	}

	po.DisplayName = provider.DisplayName
	po.Description = provider.Description
	po.APIBaseURL = provider.APIBaseURL
	po.DefaultAPIKey = provider.DefaultAPIKey
	po.Status = provider.Status
	po.SetDefaultHeaders(provider.DefaultHeaders)
	po.SetConfig(provider.Config)

	if err := r.data.db.WithContext(ctx).Save(po).Error; err != nil {
		return nil, err
	}

	return r.convertProviderModelToBiz(po), nil
}

func (r *providerRepo) DeleteProvider(ctx context.Context, id int64) error {
	return r.data.db.WithContext(ctx).Delete(&model.Provider{}, id).Error
}

func (r *providerRepo) GetProvider(ctx context.Context, id int64) (*biz.Provider, error) {
	po := &model.Provider{}
	if err := r.data.db.WithContext(ctx).First(po, id).Error; err != nil {
		return nil, err
	}
	return r.convertProviderModelToBiz(po), nil
}

func (r *providerRepo) GetProviderByName(ctx context.Context, name string) (*biz.Provider, error) {
	po := &model.Provider{}
	if err := r.data.db.WithContext(ctx).Where("name = ?", name).First(po).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return r.convertProviderModelToBiz(po), nil
}

func (r *providerRepo) ListProviders(ctx context.Context, status int32, page, pageSize int32) ([]*biz.Provider, int64, error) {
	var providers []*model.Provider
	var total int64

	query := r.data.db.WithContext(ctx).Model(&model.Provider{})
	if status >= 0 {
		query = query.Where("status = ?", status)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	if err := query.Offset(int(offset)).Limit(int(pageSize)).Find(&providers).Error; err != nil {
		return nil, 0, err
	}

	bizProviders := make([]*biz.Provider, len(providers))
	for i, po := range providers {
		bizProviders[i] = r.convertProviderModelToBiz(po)
	}

	return bizProviders, total, nil
}

func (r *providerRepo) convertProviderModelToBiz(po *model.Provider) *biz.Provider {
	return &biz.Provider{
		ID:             po.ID,
		Name:           po.Name,
		DisplayName:    po.DisplayName,
		Description:    po.Description,
		APIBaseURL:     po.APIBaseURL,
		DefaultAPIKey:  po.DefaultAPIKey,
		DefaultHeaders: po.GetDefaultHeaders(),
		Config:         po.GetConfig(),
		Status:         po.Status,
		CreatedAt:      po.CreatedAt,
		UpdatedAt:      po.UpdatedAt,
	}
}

// 模型仓库实现
type modelRepo struct {
	data *Data
	log  *log.Helper
}

// NewModelRepo 创建模型仓库
func NewModelRepo(data *Data, logger log.Logger) biz.ModelRepo {
	return &modelRepo{
		data: data,
		log:  log.NewHelper(logger),
	}
}

func (r *modelRepo) CreateModel(ctx context.Context, bizModel *biz.Model) (*biz.Model, error) {
	po := &model.Model{
		ProviderID:  bizModel.ProviderID,
		Name:        bizModel.Name,
		DisplayName: bizModel.DisplayName,
		Description: bizModel.Description,
		Version:     bizModel.Version,
		Status:      bizModel.Status,
	}
	po.SetCapabilities(bizModel.Capabilities)
	po.SetLimits(bizModel.Limits)
	po.SetPricing(bizModel.Pricing)
	po.SetDefaultParams(bizModel.DefaultParams)

	if err := r.data.db.WithContext(ctx).Create(po).Error; err != nil {
		return nil, err
	}

	return r.convertModelModelToBiz(po), nil
}

func (r *modelRepo) UpdateModel(ctx context.Context, modelBiz *biz.Model) (*biz.Model, error) {
	po := &model.Model{}
	if err := r.data.db.WithContext(ctx).First(po, modelBiz.ID).Error; err != nil {
		return nil, err
	}

	po.DisplayName = modelBiz.DisplayName
	po.Description = modelBiz.Description
	po.Status = modelBiz.Status
	po.SetCapabilities(modelBiz.Capabilities)
	po.SetLimits(modelBiz.Limits)
	po.SetPricing(modelBiz.Pricing)
	po.SetDefaultParams(modelBiz.DefaultParams)

	if err := r.data.db.WithContext(ctx).Save(po).Error; err != nil {
		return nil, err
	}

	return r.convertModelModelToBiz(po), nil
}

func (r *modelRepo) DeleteModel(ctx context.Context, id int64) error {
	return r.data.db.WithContext(ctx).Delete(&model.Model{}, id).Error
}

func (r *modelRepo) GetModel(ctx context.Context, id int64) (*biz.Model, error) {
	po := &model.Model{}
	if err := r.data.db.WithContext(ctx).First(po, id).Error; err != nil {
		return nil, err
	}
	return r.convertModelModelToBiz(po), nil
}

func (r *modelRepo) GetModelByName(ctx context.Context, name string) (*biz.Model, error) {
	po := &model.Model{}
	if err := r.data.db.WithContext(ctx).Where("name = ?", name).First(po).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return r.convertModelModelToBiz(po), nil
}

func (r *modelRepo) ListModels(ctx context.Context, providerID int64, status int32, capabilities []string, page, pageSize int32) ([]*biz.Model, int64, error) {
	var models []*model.Model
	var total int64

	query := r.data.db.WithContext(ctx).Model(&model.Model{})

	if providerID > 0 {
		query = query.Where("provider_id = ?", providerID)
	}

	if status >= 0 {
		query = query.Where("status = ?", status)
	}

	// 能力过滤（简单实现，生产环境可能需要更复杂的JSON查询）
	if len(capabilities) > 0 {
		for _, capability := range capabilities {
			query = query.Where("capabilities LIKE ?", "%"+capability+"%")
		}
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	if err := query.Offset(int(offset)).Limit(int(pageSize)).Find(&models).Error; err != nil {
		return nil, 0, err
	}

	bizModels := make([]*biz.Model, len(models))
	for i, po := range models {
		bizModels[i] = r.convertModelModelToBiz(po)
	}

	return bizModels, total, nil
}

func (r *modelRepo) SwitchUserDefaultModel(ctx context.Context, userID, modelID int64) error {
	// 使用事务确保数据一致性
	return r.data.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// 删除旧的默认模型记录
		if err := tx.Where("user_id = ?", userID).Delete(&model.UserDefaultModel{}).Error; err != nil {
			return err
		}

		// 插入新的默认模型记录
		userDefault := &model.UserDefaultModel{
			UserID:  userID,
			ModelID: modelID,
		}
		return tx.Create(userDefault).Error
	})
}

func (r *modelRepo) convertModelModelToBiz(po *model.Model) *biz.Model {
	return &biz.Model{
		ID:            po.ID,
		ProviderID:    po.ProviderID,
		Name:          po.Name,
		DisplayName:   po.DisplayName,
		Description:   po.Description,
		Version:       po.Version,
		Capabilities:  po.GetCapabilities(),
		Limits:        po.GetLimits(),
		Pricing:       po.GetPricing(),
		DefaultParams: po.GetDefaultParams(),
		Status:        po.Status,
		CreatedAt:     po.CreatedAt,
		UpdatedAt:     po.UpdatedAt,
	}
}

// 配额仓库实现
type quotaRepo struct {
	data *Data
	log  *log.Helper
}

// NewQuotaRepo 创建配额仓库
func NewQuotaRepo(data *Data, logger log.Logger) biz.QuotaRepo {
	return &quotaRepo{
		data: data,
		log:  log.NewHelper(logger),
	}
}

func (r *quotaRepo) GetUserQuota(ctx context.Context, userID, modelID int64) (*biz.UserQuota, error) {
	po := &model.UserQuota{}
	if err := r.data.db.WithContext(ctx).Where("user_id = ? AND model_id = ?", userID, modelID).First(po).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			// 创建默认配额
			po = &model.UserQuota{
				UserID:              userID,
				ModelID:             modelID,
				DailyTokenLimit:     10000,  // 默认每日10K tokens
				MonthlyTokenLimit:   300000, // 默认每月300K tokens
				DailyRequestLimit:   100,    // 默认每日100次请求
				MonthlyRequestLimit: 3000,   // 默认每月3000次请求
				ResetDailyAt:        time.Now().AddDate(0, 0, 1),
				ResetMonthlyAt:      time.Now().AddDate(0, 1, 0),
			}
			if err := r.data.db.WithContext(ctx).Create(po).Error; err != nil {
				return nil, err
			}
		} else {
			return nil, err
		}
	}

	return r.convertUserQuotaModelToBiz(po), nil
}

func (r *quotaRepo) UpdateUserQuota(ctx context.Context, quota *biz.UserQuota) (*biz.UserQuota, error) {
	po := &model.UserQuota{
		UserID:              quota.UserID,
		ModelID:             quota.ModelID,
		DailyTokenLimit:     quota.DailyTokenLimit,
		MonthlyTokenLimit:   quota.MonthlyTokenLimit,
		DailyRequestLimit:   quota.DailyRequestLimit,
		MonthlyRequestLimit: quota.MonthlyRequestLimit,
		DailyTokensUsed:     quota.DailyTokensUsed,
		MonthlyTokensUsed:   quota.MonthlyTokensUsed,
		DailyRequestsUsed:   quota.DailyRequestsUsed,
		MonthlyRequestsUsed: quota.MonthlyRequestsUsed,
		ResetDailyAt:        quota.ResetDailyAt,
		ResetMonthlyAt:      quota.ResetMonthlyAt,
	}

	if err := r.data.db.WithContext(ctx).Where("user_id = ? AND model_id = ?", quota.UserID, quota.ModelID).Save(po).Error; err != nil {
		return nil, err
	}

	return r.convertUserQuotaModelToBiz(po), nil
}

func (r *quotaRepo) GetUsageStats(ctx context.Context, userID, modelID int64, startDate, endDate, groupBy string) ([]*biz.UsageStats, error) {
	var stats []*model.UsageStats

	query := r.data.db.WithContext(ctx).Model(&model.UsageStats{})

	if userID > 0 {
		query = query.Where("user_id = ?", userID)
	}

	if modelID > 0 {
		query = query.Where("model_id = ?", modelID)
	}

	if startDate != "" && endDate != "" {
		query = query.Where("date BETWEEN ? AND ?", startDate, endDate)
	}

	if err := query.Find(&stats).Error; err != nil {
		return nil, err
	}

	bizStats := make([]*biz.UsageStats, len(stats))
	for i, po := range stats {
		bizStats[i] = r.convertUsageStatsModelToBiz(po)
	}

	return bizStats, nil
}

func (r *quotaRepo) ResetUsage(ctx context.Context, userID, modelID int64, resetType string) (int64, error) {
	query := r.data.db.WithContext(ctx).Model(&model.UserQuota{})

	if userID > 0 {
		query = query.Where("user_id = ?", userID)
	}

	if modelID > 0 {
		query = query.Where("model_id = ?", modelID)
	}

	updates := make(map[string]interface{})

	switch resetType {
	case "daily":
		updates["daily_tokens_used"] = 0
		updates["daily_requests_used"] = 0
		updates["reset_daily_at"] = time.Now().AddDate(0, 0, 1)
	case "monthly":
		updates["monthly_tokens_used"] = 0
		updates["monthly_requests_used"] = 0
		updates["reset_monthly_at"] = time.Now().AddDate(0, 1, 0)
	case "all":
		updates["daily_tokens_used"] = 0
		updates["daily_requests_used"] = 0
		updates["monthly_tokens_used"] = 0
		updates["monthly_requests_used"] = 0
		updates["reset_daily_at"] = time.Now().AddDate(0, 0, 1)
		updates["reset_monthly_at"] = time.Now().AddDate(0, 1, 0)
	}

	result := query.Updates(updates)
	return result.RowsAffected, result.Error
}

func (r *quotaRepo) CheckRateLimit(ctx context.Context, userID, modelID int64, requestedTokens int32) (bool, string, int32, int32, int32, error) {
	// 获取用户配额
	quota, err := r.GetUserQuota(ctx, userID, modelID)
	if err != nil {
		return false, "failed to get user quota", 0, 0, 0, err
	}

	// 检查每日token限制
	if quota.DailyTokensUsed+int64(requestedTokens) > quota.DailyTokenLimit {
		remaining := quota.DailyTokenLimit - quota.DailyTokensUsed
		if remaining < 0 {
			remaining = 0
		}
		return false, "daily token limit exceeded", 0, int32(quota.DailyRequestLimit - quota.DailyRequestsUsed), int32(remaining), nil
	}

	// 检查每日请求限制
	if quota.DailyRequestsUsed >= quota.DailyRequestLimit {
		return false, "daily request limit exceeded", 0, 0, int32(quota.DailyTokenLimit - quota.DailyTokensUsed), nil
	}

	// 检查每月限制
	if quota.MonthlyTokensUsed+int64(requestedTokens) > quota.MonthlyTokenLimit {
		return false, "monthly token limit exceeded", 0, int32(quota.DailyRequestLimit - quota.DailyRequestsUsed), int32(quota.DailyTokenLimit - quota.DailyTokensUsed), nil
	}

	if quota.MonthlyRequestsUsed >= quota.MonthlyRequestLimit {
		return false, "monthly request limit exceeded", 0, int32(quota.DailyRequestLimit - quota.DailyRequestsUsed), int32(quota.DailyTokenLimit - quota.DailyTokensUsed), nil
	}

	return true, "", 0, int32(quota.DailyRequestLimit - quota.DailyRequestsUsed), int32(quota.DailyTokenLimit - quota.DailyTokensUsed), nil
}

func (r *quotaRepo) convertUserQuotaModelToBiz(po *model.UserQuota) *biz.UserQuota {
	return &biz.UserQuota{
		UserID:              po.UserID,
		ModelID:             po.ModelID,
		DailyTokenLimit:     po.DailyTokenLimit,
		MonthlyTokenLimit:   po.MonthlyTokenLimit,
		DailyRequestLimit:   po.DailyRequestLimit,
		MonthlyRequestLimit: po.MonthlyRequestLimit,
		DailyTokensUsed:     po.DailyTokensUsed,
		MonthlyTokensUsed:   po.MonthlyTokensUsed,
		DailyRequestsUsed:   po.DailyRequestsUsed,
		MonthlyRequestsUsed: po.MonthlyRequestsUsed,
		ResetDailyAt:        po.ResetDailyAt,
		ResetMonthlyAt:      po.ResetMonthlyAt,
	}
}

func (r *quotaRepo) convertUsageStatsModelToBiz(po *model.UsageStats) *biz.UsageStats {
	return &biz.UsageStats{
		UserID:             po.UserID,
		ModelID:            po.ModelID,
		Date:               po.Date,
		TotalRequests:      po.TotalRequests,
		SuccessfulRequests: po.SuccessfulRequests,
		FailedRequests:     po.FailedRequests,
		TotalTokens:        po.TotalTokens,
		InputTokens:        po.InputTokens,
		OutputTokens:       po.OutputTokens,
		TotalCost:          po.TotalCost,
		AvgResponseTime:    po.AvgResponseTime,
	}
}

// 限流配置仓库实现
type rateLimitRepo struct {
	data *Data
	log  *log.Helper
}

// NewRateLimitRepo 创建限流配置仓库
func NewRateLimitRepo(data *Data, logger log.Logger) biz.RateLimitRepo {
	return &rateLimitRepo{
		data: data,
		log:  log.NewHelper(logger),
	}
}

func (r *rateLimitRepo) GetRateLimitConfig(ctx context.Context, modelID int64, userLevel string) ([]*biz.RateLimitConfig, error) {
	var configs []*model.RateLimitConfig

	query := r.data.db.WithContext(ctx).Model(&model.RateLimitConfig{})

	if modelID > 0 {
		query = query.Where("model_id = ?", modelID)
	}

	if userLevel != "" {
		query = query.Where("user_level = ?", userLevel)
	}

	if err := query.Find(&configs).Error; err != nil {
		return nil, err
	}

	bizConfigs := make([]*biz.RateLimitConfig, len(configs))
	for i, po := range configs {
		bizConfigs[i] = r.convertRateLimitConfigModelToBiz(po)
	}

	return bizConfigs, nil
}

func (r *rateLimitRepo) UpdateRateLimitConfig(ctx context.Context, config *biz.RateLimitConfig) (*biz.RateLimitConfig, error) {
	po := &model.RateLimitConfig{
		ID:                 config.ID,
		ModelID:            config.ModelID,
		UserLevel:          config.UserLevel,
		RequestsPerMinute:  config.RequestsPerMinute,
		RequestsPerHour:    config.RequestsPerHour,
		RequestsPerDay:     config.RequestsPerDay,
		TokensPerMinute:    config.TokensPerMinute,
		ConcurrentRequests: config.ConcurrentRequests,
		BurstLimit:         config.BurstLimit,
	}

	if config.ID > 0 {
		// 更新
		if err := r.data.db.WithContext(ctx).Save(po).Error; err != nil {
			return nil, err
		}
	} else {
		// 创建
		if err := r.data.db.WithContext(ctx).Create(po).Error; err != nil {
			return nil, err
		}
	}

	return r.convertRateLimitConfigModelToBiz(po), nil
}

func (r *rateLimitRepo) convertRateLimitConfigModelToBiz(po *model.RateLimitConfig) *biz.RateLimitConfig {
	return &biz.RateLimitConfig{
		ID:                 po.ID,
		ModelID:            po.ModelID,
		UserLevel:          po.UserLevel,
		RequestsPerMinute:  po.RequestsPerMinute,
		RequestsPerHour:    po.RequestsPerHour,
		RequestsPerDay:     po.RequestsPerDay,
		TokensPerMinute:    po.TokensPerMinute,
		ConcurrentRequests: po.ConcurrentRequests,
		BurstLimit:         po.BurstLimit,
		CreatedAt:          po.CreatedAt,
		UpdatedAt:          po.UpdatedAt,
	}
}

// 健康检查仓库实现
type healthRepo struct {
	data *Data
	log  *log.Helper
}

// NewHealthRepo 创建健康检查仓库
func NewHealthRepo(data *Data, logger log.Logger) biz.HealthRepo {
	return &healthRepo{
		data: data,
		log:  log.NewHelper(logger),
	}
}

func (r *healthRepo) GetModelHealth(ctx context.Context, modelID int64) (*biz.ModelHealth, error) {
	po := &model.ModelHealth{}
	if err := r.data.db.WithContext(ctx).Where("model_id = ?", modelID).First(po).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return &biz.ModelHealth{
				ModelID:   modelID,
				IsHealthy: false,
			}, nil
		}
		return nil, err
	}

	return r.convertModelHealthModelToBiz(po), nil
}

func (r *healthRepo) UpdateModelHealth(ctx context.Context, health *biz.ModelHealth) error {
	po := &model.ModelHealth{
		ModelID:        health.ModelID,
		IsHealthy:      health.IsHealthy,
		ResponseTime:   health.ResponseTime,
		SuccessRate:    health.SuccessRate,
		TotalRequests:  health.TotalRequests,
		FailedRequests: health.FailedRequests,
		ErrorMessage:   health.ErrorMessage,
		LastCheck:      health.LastCheck,
	}

	return r.data.db.WithContext(ctx).Where("model_id = ?", health.ModelID).Save(po).Error
}

func (r *healthRepo) GetModelMetrics(ctx context.Context, modelID int64, startTime, endTime string) (*pb.GetModelMetricsReply, error) {
	// 这里可以实现复杂的指标查询逻辑
	// 简单实现：从usage_stats表聚合数据
	var result struct {
		AvgResponseTime float64
		SuccessRate     float64
		TotalRequests   int64
		ErrorCount      int64
	}

	query := `
		SELECT 
			AVG(avg_response_time) as avg_response_time,
			SUM(successful_requests) / SUM(total_requests) as success_rate,
			SUM(total_requests) as total_requests,
			SUM(failed_requests) as error_count
		FROM ai_usage_stats 
		WHERE model_id = ?
	`

	args := []interface{}{modelID}

	if startTime != "" && endTime != "" {
		query += " AND date BETWEEN ? AND ?"
		args = append(args, startTime, endTime)
	}

	if err := r.data.db.WithContext(ctx).Raw(query, args...).Scan(&result).Error; err != nil {
		return nil, err
	}

	// 获取主要错误类型（简单实现）
	var topErrors []string
	// 这里可以实现更复杂的错误统计逻辑

	return &pb.GetModelMetricsReply{
		ModelId:         modelID,
		AvgResponseTime: result.AvgResponseTime,
		SuccessRate:     result.SuccessRate,
		TotalRequests:   result.TotalRequests,
		ErrorCount:      result.ErrorCount,
		TopErrors:       topErrors,
	}, nil
}

func (r *healthRepo) convertModelHealthModelToBiz(po *model.ModelHealth) *biz.ModelHealth {
	return &biz.ModelHealth{
		ModelID:        po.ModelID,
		IsHealthy:      po.IsHealthy,
		ResponseTime:   po.ResponseTime,
		SuccessRate:    po.SuccessRate,
		TotalRequests:  po.TotalRequests,
		FailedRequests: po.FailedRequests,
		ErrorMessage:   po.ErrorMessage,
		LastCheck:      po.LastCheck,
	}
}
