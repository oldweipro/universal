package biz

import (
	"context"
	"fmt"
	"net/http"
	"time"

	pb "universal/api/ai/v1"

	"github.com/go-kratos/kratos/v2/log"
)

// Provider 提供商业务实体
type Provider struct {
	ID             int64             `json:"id"`
	Name           string            `json:"name"`
	DisplayName    string            `json:"display_name"`
	Description    string            `json:"description"`
	APIBaseURL     string            `json:"api_base_url"`
	DefaultAPIKey  string            `json:"default_api_key"`
	DefaultHeaders map[string]string `json:"default_headers"`
	Config         map[string]string `json:"config"`
	Status         int32             `json:"status"`
	CreatedAt      time.Time         `json:"created_at"`
	UpdatedAt      time.Time         `json:"updated_at"`
}

// Model 模型业务实体
type Model struct {
	ID            int64                 `json:"id"`
	ProviderID    int64                 `json:"provider_id"`
	Name          string                `json:"name"`
	DisplayName   string                `json:"display_name"`
	Description   string                `json:"description"`
	Version       string                `json:"version"`
	Capabilities  *pb.ModelCapabilities `json:"capabilities"`
	Limits        *pb.ModelLimits       `json:"limits"`
	Pricing       *pb.ModelPricing      `json:"pricing"`
	DefaultParams map[string]string     `json:"default_params"`
	Status        int32                 `json:"status"`
	CreatedAt     time.Time             `json:"created_at"`
	UpdatedAt     time.Time             `json:"updated_at"`
}

// UserQuota 用户配额业务实体
type UserQuota struct {
	UserID              int64     `json:"user_id"`
	ModelID             int64     `json:"model_id"`
	DailyTokenLimit     int64     `json:"daily_token_limit"`
	MonthlyTokenLimit   int64     `json:"monthly_token_limit"`
	DailyRequestLimit   int64     `json:"daily_request_limit"`
	MonthlyRequestLimit int64     `json:"monthly_request_limit"`
	DailyTokensUsed     int64     `json:"daily_tokens_used"`
	MonthlyTokensUsed   int64     `json:"monthly_tokens_used"`
	DailyRequestsUsed   int64     `json:"daily_requests_used"`
	MonthlyRequestsUsed int64     `json:"monthly_requests_used"`
	ResetDailyAt        time.Time `json:"reset_daily_at"`
	ResetMonthlyAt      time.Time `json:"reset_monthly_at"`
}

// UsageStats 使用统计业务实体
type UsageStats struct {
	UserID             int64   `json:"user_id"`
	ModelID            int64   `json:"model_id"`
	Date               string  `json:"date"`
	TotalRequests      int64   `json:"total_requests"`
	SuccessfulRequests int64   `json:"successful_requests"`
	FailedRequests     int64   `json:"failed_requests"`
	TotalTokens        int64   `json:"total_tokens"`
	InputTokens        int64   `json:"input_tokens"`
	OutputTokens       int64   `json:"output_tokens"`
	TotalCost          float64 `json:"total_cost"`
	AvgResponseTime    float64 `json:"avg_response_time"`
}

// RateLimitConfig 限流配置业务实体
type RateLimitConfig struct {
	ID                 int64     `json:"id"`
	ModelID            int64     `json:"model_id"`
	UserLevel          string    `json:"user_level"`
	RequestsPerMinute  int32     `json:"requests_per_minute"`
	RequestsPerHour    int32     `json:"requests_per_hour"`
	RequestsPerDay     int32     `json:"requests_per_day"`
	TokensPerMinute    int32     `json:"tokens_per_minute"`
	ConcurrentRequests int32     `json:"concurrent_requests"`
	BurstLimit         int32     `json:"burst_limit"`
	CreatedAt          time.Time `json:"created_at"`
	UpdatedAt          time.Time `json:"updated_at"`
}

// ModelHealth 模型健康状态业务实体
type ModelHealth struct {
	ModelID        int64     `json:"model_id"`
	IsHealthy      bool      `json:"is_healthy"`
	ResponseTime   float64   `json:"response_time"`
	SuccessRate    float64   `json:"success_rate"`
	TotalRequests  int64     `json:"total_requests"`
	FailedRequests int64     `json:"failed_requests"`
	ErrorMessage   string    `json:"error_message"`
	LastCheck      time.Time `json:"last_check"`
}

// ProviderRepo 提供商数据访问接口
type ProviderRepo interface {
	CreateProvider(ctx context.Context, provider *Provider) (*Provider, error)
	UpdateProvider(ctx context.Context, provider *Provider) (*Provider, error)
	DeleteProvider(ctx context.Context, id int64) error
	GetProvider(ctx context.Context, id int64) (*Provider, error)
	GetProviderByName(ctx context.Context, name string) (*Provider, error)
	ListProviders(ctx context.Context, status int32, page, pageSize int32) ([]*Provider, int64, error)
}

// ModelRepo 模型数据访问接口
type ModelRepo interface {
	CreateModel(ctx context.Context, model *Model) (*Model, error)
	UpdateModel(ctx context.Context, model *Model) (*Model, error)
	DeleteModel(ctx context.Context, id int64) error
	GetModel(ctx context.Context, id int64) (*Model, error)
	GetModelByName(ctx context.Context, name string) (*Model, error)
	ListModels(ctx context.Context, providerID int64, status int32, capabilities []string, page, pageSize int32) ([]*Model, int64, error)
	SwitchUserDefaultModel(ctx context.Context, userID, modelID int64) error
}

// QuotaRepo 配额数据访问接口
type QuotaRepo interface {
	GetUserQuota(ctx context.Context, userID, modelID int64) (*UserQuota, error)
	UpdateUserQuota(ctx context.Context, quota *UserQuota) (*UserQuota, error)
	GetUsageStats(ctx context.Context, userID, modelID int64, startDate, endDate, groupBy string) ([]*UsageStats, error)
	ResetUsage(ctx context.Context, userID, modelID int64, resetType string) (int64, error)
	CheckRateLimit(ctx context.Context, userID, modelID int64, requestedTokens int32) (bool, string, int32, int32, int32, error)
}

// RateLimitRepo 限流配置数据访问接口
type RateLimitRepo interface {
	GetRateLimitConfig(ctx context.Context, modelID int64, userLevel string) ([]*RateLimitConfig, error)
	UpdateRateLimitConfig(ctx context.Context, config *RateLimitConfig) (*RateLimitConfig, error)
}

// HealthRepo 健康检查数据访问接口
type HealthRepo interface {
	GetModelHealth(ctx context.Context, modelID int64) (*ModelHealth, error)
	UpdateModelHealth(ctx context.Context, health *ModelHealth) error
	GetModelMetrics(ctx context.Context, modelID int64, startTime, endTime string) (*pb.GetModelMetricsReply, error)
}

// ModelUsecase 模型管理业务用例
type ModelUsecase struct {
	providerRepo  ProviderRepo
	modelRepo     ModelRepo
	quotaRepo     QuotaRepo
	rateLimitRepo RateLimitRepo
	healthRepo    HealthRepo
	log           *log.Helper
}

// NewModelUsecase 创建模型管理业务用例
func NewModelUsecase(
	providerRepo ProviderRepo,
	modelRepo ModelRepo,
	quotaRepo QuotaRepo,
	rateLimitRepo RateLimitRepo,
	healthRepo HealthRepo,
	logger log.Logger,
) *ModelUsecase {
	return &ModelUsecase{
		providerRepo:  providerRepo,
		modelRepo:     modelRepo,
		quotaRepo:     quotaRepo,
		rateLimitRepo: rateLimitRepo,
		healthRepo:    healthRepo,
		log:           log.NewHelper(logger),
	}
}

// 提供商管理相关方法

// CreateProvider 创建模型提供商
func (uc *ModelUsecase) CreateProvider(ctx context.Context, req *pb.CreateProviderRequest) (*Provider, error) {
	// 检查提供商名称是否已存在
	existing, _ := uc.providerRepo.GetProviderByName(ctx, req.Name)
	if existing != nil {
		return nil, fmt.Errorf("provider name already exists: %s", req.Name)
	}

	provider := &Provider{
		Name:           req.Name,
		DisplayName:    req.DisplayName,
		Description:    req.Description,
		APIBaseURL:     req.ApiBaseUrl,
		DefaultAPIKey:  req.DefaultApiKey,
		DefaultHeaders: req.DefaultHeaders,
		Config:         req.Config,
		Status:         0, // 默认启用
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}

	return uc.providerRepo.CreateProvider(ctx, provider)
}

// UpdateProvider 更新模型提供商
func (uc *ModelUsecase) UpdateProvider(ctx context.Context, req *pb.UpdateProviderRequest) (*Provider, error) {
	// 检查提供商是否存在
	existing, err := uc.providerRepo.GetProvider(ctx, req.Id)
	if err != nil {
		return nil, fmt.Errorf("provider not found: %w", err)
	}

	// 更新字段
	if req.DisplayName != "" {
		existing.DisplayName = req.DisplayName
	}
	if req.Description != "" {
		existing.Description = req.Description
	}
	if req.ApiBaseUrl != "" {
		existing.APIBaseURL = req.ApiBaseUrl
	}
	if req.DefaultApiKey != "" {
		existing.DefaultAPIKey = req.DefaultApiKey
	}
	if req.DefaultHeaders != nil {
		existing.DefaultHeaders = req.DefaultHeaders
	}
	if req.Config != nil {
		existing.Config = req.Config
	}
	if req.Status >= 0 {
		existing.Status = req.Status
	}
	existing.UpdatedAt = time.Now()

	return uc.providerRepo.UpdateProvider(ctx, existing)
}

// DeleteProvider 删除模型提供商
func (uc *ModelUsecase) DeleteProvider(ctx context.Context, id int64) error {
	// 检查提供商是否存在
	_, err := uc.providerRepo.GetProvider(ctx, id)
	if err != nil {
		return fmt.Errorf("provider not found: %w", err)
	}

	return uc.providerRepo.DeleteProvider(ctx, id)
}

// ListProviders 获取提供商列表
func (uc *ModelUsecase) ListProviders(ctx context.Context, status int32, page, pageSize int32) ([]*Provider, int64, error) {
	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 || pageSize > 100 {
		pageSize = 20
	}

	return uc.providerRepo.ListProviders(ctx, status, page, pageSize)
}

// TestProvider 测试提供商连接
func (uc *ModelUsecase) TestProvider(ctx context.Context, id int64, apiKey string) (bool, float64, string, error) {
	provider, err := uc.providerRepo.GetProvider(ctx, id)
	if err != nil {
		return false, 0, "provider not found", err
	}

	// 使用提供的 API Key 或默认的
	testKey := apiKey
	if testKey == "" {
		testKey = provider.DefaultAPIKey
	}

	// 执行健康检查请求
	startTime := time.Now()

	// 根据提供商类型构建测试请求
	var testURL string
	var testHeaders map[string]string

	switch provider.Name {
	case "openai":
		testURL = provider.APIBaseURL + "/v1/models"
		testHeaders = map[string]string{
			"Authorization": "Bearer " + testKey,
			"Content-Type":  "application/json",
		}
	case "anthropic":
		testURL = provider.APIBaseURL + "/v1/messages"
		testHeaders = map[string]string{
			"x-api-key":    testKey,
			"Content-Type": "application/json",
		}
	default:
		// 通用测试：尝试 GET 请求到基础 URL
		testURL = provider.APIBaseURL
		testHeaders = provider.DefaultHeaders
		if testKey != "" {
			testHeaders["Authorization"] = "Bearer " + testKey
		}
	}

	// 创建 HTTP 客户端并执行请求
	client := &http.Client{Timeout: 10 * time.Second}
	req, err := http.NewRequestWithContext(ctx, "GET", testURL, nil)
	if err != nil {
		return false, 0, fmt.Sprintf("failed to create request: %v", err), err
	}

	for key, value := range testHeaders {
		req.Header.Set(key, value)
	}

	resp, err := client.Do(req)
	responseTime := float64(time.Since(startTime).Nanoseconds()) / 1e6 // 转换为毫秒

	if err != nil {
		return false, responseTime, fmt.Sprintf("request failed: %v", err), err
	}
	defer resp.Body.Close()

	// 检查响应状态
	if resp.StatusCode >= 200 && resp.StatusCode < 400 {
		return true, responseTime, "", nil
	}

	return false, responseTime, fmt.Sprintf("unexpected status code: %d", resp.StatusCode), nil
}

// 模型管理相关方法

// CreateModel 创建模型
func (uc *ModelUsecase) CreateModel(ctx context.Context, req *pb.CreateModelRequest) (*Model, error) {
	// 检查提供商是否存在
	_, err := uc.providerRepo.GetProvider(ctx, req.ProviderId)
	if err != nil {
		return nil, fmt.Errorf("provider not found: %w", err)
	}

	// 检查模型名称是否已存在
	existing, _ := uc.modelRepo.GetModelByName(ctx, req.Name)
	if existing != nil {
		return nil, fmt.Errorf("model name already exists: %s", req.Name)
	}

	model := &Model{
		ProviderID:    req.ProviderId,
		Name:          req.Name,
		DisplayName:   req.DisplayName,
		Description:   req.Description,
		Version:       req.Version,
		Capabilities:  req.Capabilities,
		Limits:        req.Limits,
		Pricing:       req.Pricing,
		DefaultParams: req.DefaultParams,
		Status:        0, // 默认可用
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}

	return uc.modelRepo.CreateModel(ctx, model)
}

// UpdateModel 更新模型
func (uc *ModelUsecase) UpdateModel(ctx context.Context, req *pb.UpdateModelRequest) (*Model, error) {
	existing, err := uc.modelRepo.GetModel(ctx, req.Id)
	if err != nil {
		return nil, fmt.Errorf("model not found: %w", err)
	}

	// 更新字段
	if req.DisplayName != "" {
		existing.DisplayName = req.DisplayName
	}
	if req.Description != "" {
		existing.Description = req.Description
	}
	if req.Capabilities != nil {
		existing.Capabilities = req.Capabilities
	}
	if req.Limits != nil {
		existing.Limits = req.Limits
	}
	if req.Pricing != nil {
		existing.Pricing = req.Pricing
	}
	if req.DefaultParams != nil {
		existing.DefaultParams = req.DefaultParams
	}
	if req.Status >= 0 {
		existing.Status = req.Status
	}
	existing.UpdatedAt = time.Now()

	return uc.modelRepo.UpdateModel(ctx, existing)
}

// DeleteModel 删除模型
func (uc *ModelUsecase) DeleteModel(ctx context.Context, id int64) error {
	_, err := uc.modelRepo.GetModel(ctx, id)
	if err != nil {
		return fmt.Errorf("model not found: %w", err)
	}

	return uc.modelRepo.DeleteModel(ctx, id)
}

// ListModels 获取模型列表
func (uc *ModelUsecase) ListModels(ctx context.Context, providerID int64, status int32, capabilities []string, page, pageSize int32) ([]*Model, int64, error) {
	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 || pageSize > 100 {
		pageSize = 20
	}

	return uc.modelRepo.ListModels(ctx, providerID, status, capabilities, page, pageSize)
}

// GetModel 获取模型详情
func (uc *ModelUsecase) GetModel(ctx context.Context, id int64) (*Model, error) {
	return uc.modelRepo.GetModel(ctx, id)
}

// SwitchModel 切换用户默认模型
func (uc *ModelUsecase) SwitchModel(ctx context.Context, userID, modelID int64) error {
	// 检查模型是否存在且可用
	model, err := uc.modelRepo.GetModel(ctx, modelID)
	if err != nil {
		return fmt.Errorf("model not found: %w", err)
	}

	if model.Status != 0 {
		return fmt.Errorf("model is not available")
	}

	return uc.modelRepo.SwitchUserDefaultModel(ctx, userID, modelID)
}

// 配额管理相关方法

// GetUserQuota 获取用户配额
func (uc *ModelUsecase) GetUserQuota(ctx context.Context, userID, modelID int64) (*UserQuota, error) {
	return uc.quotaRepo.GetUserQuota(ctx, userID, modelID)
}

// UpdateUserQuota 更新用户配额
func (uc *ModelUsecase) UpdateUserQuota(ctx context.Context, quota *UserQuota) (*UserQuota, error) {
	return uc.quotaRepo.UpdateUserQuota(ctx, quota)
}

// GetUsageStats 获取使用统计
func (uc *ModelUsecase) GetUsageStats(ctx context.Context, userID, modelID int64, startDate, endDate, groupBy string) ([]*UsageStats, error) {
	return uc.quotaRepo.GetUsageStats(ctx, userID, modelID, startDate, endDate, groupBy)
}

// ResetUsage 重置使用量
func (uc *ModelUsecase) ResetUsage(ctx context.Context, userID, modelID int64, resetType string) (int64, error) {
	return uc.quotaRepo.ResetUsage(ctx, userID, modelID, resetType)
}

// CheckRateLimit 检查速率限制
func (uc *ModelUsecase) CheckRateLimit(ctx context.Context, userID, modelID int64, requestedTokens int32) (bool, string, int32, int32, int32, error) {
	return uc.quotaRepo.CheckRateLimit(ctx, userID, modelID, requestedTokens)
}

// 限流配置相关方法

// GetRateLimitConfig 获取限流配置
func (uc *ModelUsecase) GetRateLimitConfig(ctx context.Context, modelID int64, userLevel string) ([]*RateLimitConfig, error) {
	return uc.rateLimitRepo.GetRateLimitConfig(ctx, modelID, userLevel)
}

// UpdateRateLimitConfig 更新限流配置
func (uc *ModelUsecase) UpdateRateLimitConfig(ctx context.Context, config *RateLimitConfig) (*RateLimitConfig, error) {
	return uc.rateLimitRepo.UpdateRateLimitConfig(ctx, config)
}

// 健康检查相关方法

// HealthCheck 执行健康检查
func (uc *ModelUsecase) HealthCheck(ctx context.Context, modelID int64) (*ModelHealth, error) {
	// 如果 modelID 为 0，检查所有模型
	if modelID == 0 {
		// 这里可以扩展为检查所有模型的逻辑
		return nil, fmt.Errorf("health check for all models not implemented")
	}

	// 获取模型信息
	model, err := uc.modelRepo.GetModel(ctx, modelID)
	if err != nil {
		return nil, fmt.Errorf("model not found: %w", err)
	}

	// 获取提供商信息
	provider, err := uc.providerRepo.GetProvider(ctx, model.ProviderID)
	if err != nil {
		return nil, fmt.Errorf("provider not found: %w", err)
	}

	// 执行健康检查
	isHealthy, responseTime, errorMessage, _ := uc.TestProvider(ctx, provider.ID, "")

	health := &ModelHealth{
		ModelID:        modelID,
		IsHealthy:      isHealthy,
		ResponseTime:   responseTime,
		SuccessRate:    0.0, // 需要从统计数据计算
		TotalRequests:  0,   // 需要从统计数据计算
		FailedRequests: 0,   // 需要从统计数据计算
		ErrorMessage:   errorMessage,
		LastCheck:      time.Now(),
	}

	// 更新健康状态
	err = uc.healthRepo.UpdateModelHealth(ctx, health)
	if err != nil {
		uc.log.Warnf("failed to update model health: %v", err)
	}

	return health, nil
}

// GetModelMetrics 获取模型指标
func (uc *ModelUsecase) GetModelMetrics(ctx context.Context, modelID int64, startTime, endTime string) (*pb.GetModelMetricsReply, error) {
	return uc.healthRepo.GetModelMetrics(ctx, modelID, startTime, endTime)
}
