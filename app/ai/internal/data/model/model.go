package model

import (
	"encoding/json"
	"time"

	pb "universal/api/ai/v1"
)

// Provider 提供商模型
type Provider struct {
	ID             int64     `gorm:"column:id;primaryKey;autoIncrement"`
	Name           string    `gorm:"column:name;type:varchar(100);not null;uniqueIndex"`
	DisplayName    string    `gorm:"column:display_name;type:varchar(200);not null"`
	Description    string    `gorm:"column:description;type:text"`
	APIBaseURL     string    `gorm:"column:api_base_url;type:varchar(500)"`
	DefaultAPIKey  string    `gorm:"column:default_api_key;type:varchar(500)"`
	DefaultHeaders string    `gorm:"column:default_headers;type:json"`
	Config         string    `gorm:"column:config;type:json"`
	Status         int32     `gorm:"column:status;type:tinyint;default:0;comment:'0:启用 1:禁用 2:维护中'"`
	CreatedAt      time.Time `gorm:"column:created_at;autoCreateTime"`
	UpdatedAt      time.Time `gorm:"column:updated_at;autoUpdateTime"`
}

func (Provider) TableName() string {
	return "ai_providers"
}

// GetDefaultHeaders 获取默认请求头
func (p *Provider) GetDefaultHeaders() map[string]string {
	if p.DefaultHeaders == "" {
		return make(map[string]string)
	}
	var headers map[string]string
	json.Unmarshal([]byte(p.DefaultHeaders), &headers)
	return headers
}

// SetDefaultHeaders 设置默认请求头
func (p *Provider) SetDefaultHeaders(headers map[string]string) {
	if headers == nil || len(headers) == 0 {
		p.DefaultHeaders = "null"
		return
	}
	data, _ := json.Marshal(headers)
	p.DefaultHeaders = string(data)
}

// GetConfig 获取配置
func (p *Provider) GetConfig() map[string]string {
	if p.Config == "" {
		return make(map[string]string)
	}
	var config map[string]string
	json.Unmarshal([]byte(p.Config), &config)
	return config
}

// SetConfig 设置配置
func (p *Provider) SetConfig(config map[string]string) {
	if config == nil || len(config) == 0 {
		p.Config = "null"
		return
	}
	data, _ := json.Marshal(config)
	p.Config = string(data)
}

// Model 模型
type Model struct {
	ID            int64     `gorm:"column:id;primaryKey;autoIncrement"`
	ProviderID    int64     `gorm:"column:provider_id;not null;index"`
	Name          string    `gorm:"column:name;type:varchar(100);not null;uniqueIndex"`
	DisplayName   string    `gorm:"column:display_name;type:varchar(200);not null"`
	Description   string    `gorm:"column:description;type:text"`
	Version       string    `gorm:"column:version;type:varchar(50)"`
	Capabilities  string    `gorm:"column:capabilities;type:json"`
	Limits        string    `gorm:"column:limits;type:json"`
	Pricing       string    `gorm:"column:pricing;type:json"`
	DefaultParams string    `gorm:"column:default_params;type:json"`
	Status        int32     `gorm:"column:status;type:tinyint;default:0;comment:'0:可用 1:不可用 2:维护中'"`
	CreatedAt     time.Time `gorm:"column:created_at;autoCreateTime"`
	UpdatedAt     time.Time `gorm:"column:updated_at;autoUpdateTime"`
}

func (Model) TableName() string {
	return "ai_models"
}

// GetCapabilities 获取能力配置
func (m *Model) GetCapabilities() *pb.ModelCapabilities {
	if m.Capabilities == "" {
		return &pb.ModelCapabilities{}
	}
	var capabilities pb.ModelCapabilities
	json.Unmarshal([]byte(m.Capabilities), &capabilities)
	return &capabilities
}

// SetCapabilities 设置能力配置
func (m *Model) SetCapabilities(capabilities *pb.ModelCapabilities) {
	if capabilities == nil {
		m.Capabilities = "null"
		return
	}
	data, _ := json.Marshal(capabilities)
	m.Capabilities = string(data)
}

// GetLimits 获取限制配置
func (m *Model) GetLimits() *pb.ModelLimits {
	if m.Limits == "" {
		return &pb.ModelLimits{}
	}
	var limits pb.ModelLimits
	json.Unmarshal([]byte(m.Limits), &limits)
	return &limits
}

// SetLimits 设置限制配置
func (m *Model) SetLimits(limits *pb.ModelLimits) {
	if limits == nil {
		m.Limits = "null"
		return
	}
	data, _ := json.Marshal(limits)
	m.Limits = string(data)
}

// GetPricing 获取定价配置
func (m *Model) GetPricing() *pb.ModelPricing {
	if m.Pricing == "" {
		return &pb.ModelPricing{}
	}
	var pricing pb.ModelPricing
	json.Unmarshal([]byte(m.Pricing), &pricing)
	return &pricing
}

// SetPricing 设置定价配置
func (m *Model) SetPricing(pricing *pb.ModelPricing) {
	if pricing == nil {
		m.Pricing = "null"
		return
	}
	data, _ := json.Marshal(pricing)
	m.Pricing = string(data)
}

// GetDefaultParams 获取默认参数
func (m *Model) GetDefaultParams() map[string]string {
	if m.DefaultParams == "" {
		return make(map[string]string)
	}
	var params map[string]string
	json.Unmarshal([]byte(m.DefaultParams), &params)
	return params
}

// SetDefaultParams 设置默认参数
func (m *Model) SetDefaultParams(params map[string]string) {
	if params == nil || len(params) == 0 {
		m.DefaultParams = "null"
		return
	}
	data, _ := json.Marshal(params)
	m.DefaultParams = string(data)
}

// UserQuota 用户配额
type UserQuota struct {
	ID                  int64     `gorm:"column:id;primaryKey;autoIncrement"`
	UserID              int64     `gorm:"column:user_id;not null;index"`
	ModelID             int64     `gorm:"column:model_id;not null;index"`
	DailyTokenLimit     int64     `gorm:"column:daily_token_limit;default:0"`
	MonthlyTokenLimit   int64     `gorm:"column:monthly_token_limit;default:0"`
	DailyRequestLimit   int64     `gorm:"column:daily_request_limit;default:0"`
	MonthlyRequestLimit int64     `gorm:"column:monthly_request_limit;default:0"`
	DailyTokensUsed     int64     `gorm:"column:daily_tokens_used;default:0"`
	MonthlyTokensUsed   int64     `gorm:"column:monthly_tokens_used;default:0"`
	DailyRequestsUsed   int64     `gorm:"column:daily_requests_used;default:0"`
	MonthlyRequestsUsed int64     `gorm:"column:monthly_requests_used;default:0"`
	ResetDailyAt        time.Time `gorm:"column:reset_daily_at"`
	ResetMonthlyAt      time.Time `gorm:"column:reset_monthly_at"`
	CreatedAt           time.Time `gorm:"column:created_at;autoCreateTime"`
	UpdatedAt           time.Time `gorm:"column:updated_at;autoUpdateTime"`
}

func (UserQuota) TableName() string {
	return "ai_user_quotas"
}

// UsageStats 使用统计
type UsageStats struct {
	ID                 int64     `gorm:"column:id;primaryKey;autoIncrement"`
	UserID             int64     `gorm:"column:user_id;not null;index"`
	ModelID            int64     `gorm:"column:model_id;not null;index"`
	Date               string    `gorm:"column:date;type:date;not null;index"`
	TotalRequests      int64     `gorm:"column:total_requests;default:0"`
	SuccessfulRequests int64     `gorm:"column:successful_requests;default:0"`
	FailedRequests     int64     `gorm:"column:failed_requests;default:0"`
	TotalTokens        int64     `gorm:"column:total_tokens;default:0"`
	InputTokens        int64     `gorm:"column:input_tokens;default:0"`
	OutputTokens       int64     `gorm:"column:output_tokens;default:0"`
	TotalCost          float64   `gorm:"column:total_cost;type:decimal(10,4);default:0"`
	AvgResponseTime    float64   `gorm:"column:avg_response_time;type:decimal(10,2);default:0"`
	CreatedAt          time.Time `gorm:"column:created_at;autoCreateTime"`
	UpdatedAt          time.Time `gorm:"column:updated_at;autoUpdateTime"`
}

func (UsageStats) TableName() string {
	return "ai_usage_stats"
}

// RateLimitConfig 限流配置
type RateLimitConfig struct {
	ID                 int64     `gorm:"column:id;primaryKey;autoIncrement"`
	ModelID            int64     `gorm:"column:model_id;not null;index"`
	UserLevel          string    `gorm:"column:user_level;type:varchar(50);not null;index;comment:'free/pro/enterprise'"`
	RequestsPerMinute  int32     `gorm:"column:requests_per_minute;default:0"`
	RequestsPerHour    int32     `gorm:"column:requests_per_hour;default:0"`
	RequestsPerDay     int32     `gorm:"column:requests_per_day;default:0"`
	TokensPerMinute    int32     `gorm:"column:tokens_per_minute;default:0"`
	ConcurrentRequests int32     `gorm:"column:concurrent_requests;default:0"`
	BurstLimit         int32     `gorm:"column:burst_limit;default:0"`
	CreatedAt          time.Time `gorm:"column:created_at;autoCreateTime"`
	UpdatedAt          time.Time `gorm:"column:updated_at;autoUpdateTime"`
}

func (RateLimitConfig) TableName() string {
	return "ai_rate_limit_configs"
}

// ModelHealth 模型健康状态
type ModelHealth struct {
	ID             int64     `gorm:"column:id;primaryKey;autoIncrement"`
	ModelID        int64     `gorm:"column:model_id;not null;uniqueIndex"`
	IsHealthy      bool      `gorm:"column:is_healthy;default:false"`
	ResponseTime   float64   `gorm:"column:response_time;type:decimal(10,2);default:0"`
	SuccessRate    float64   `gorm:"column:success_rate;type:decimal(5,4);default:0"`
	TotalRequests  int64     `gorm:"column:total_requests;default:0"`
	FailedRequests int64     `gorm:"column:failed_requests;default:0"`
	ErrorMessage   string    `gorm:"column:error_message;type:text"`
	LastCheck      time.Time `gorm:"column:last_check"`
	CreatedAt      time.Time `gorm:"column:created_at;autoCreateTime"`
	UpdatedAt      time.Time `gorm:"column:updated_at;autoUpdateTime"`
}

func (ModelHealth) TableName() string {
	return "ai_model_health"
}

// UserDefaultModel 用户默认模型
type UserDefaultModel struct {
	ID        int64     `gorm:"column:id;primaryKey;autoIncrement"`
	UserID    int64     `gorm:"column:user_id;not null;uniqueIndex"`
	ModelID   int64     `gorm:"column:model_id;not null"`
	CreatedAt time.Time `gorm:"column:created_at;autoCreateTime"`
	UpdatedAt time.Time `gorm:"column:updated_at;autoUpdateTime"`
}

func (UserDefaultModel) TableName() string {
	return "ai_user_default_models"
}
