package model

import (
	"encoding/json"
	"time"

	"gorm.io/gorm"
)

// McpServer MCP服务器模型
type McpServer struct {
	ID                 string          `gorm:"primarykey;size:64" json:"id"`
	Name               string          `gorm:"size:255;not null" json:"name"`
	Description        string          `gorm:"type:text" json:"description"`
	Endpoint           string          `gorm:"size:500;not null" json:"endpoint"`
	Status             int             `gorm:"default:1;index" json:"status"` // 1:active, 2:inactive, 3:error, 4:maintenance, 5:deprecated
	Version            string          `gorm:"size:50" json:"version"`
	SupportedProtocols StringSlice     `gorm:"type:json" json:"supported_protocols"`
	Capabilities       StringSlice     `gorm:"type:json" json:"capabilities"`
	Owner              string          `gorm:"size:100" json:"owner"`
	Tags               StringSlice     `gorm:"type:json" json:"tags"`
	Config             McpServerConfig `gorm:"type:json" json:"config"`
	Metadata           KeyValueMap     `gorm:"type:json" json:"metadata"`
	SecurityPolicy     SecurityPolicy  `gorm:"type:json" json:"security_policy"`
	CreatedAt          time.Time       `json:"created_at"`
	UpdatedAt          time.Time       `json:"updated_at"`
	DeletedAt          gorm.DeletedAt  `gorm:"index" json:"deleted_at"`

	// 统计信息
	TotalRequests       int64      `gorm:"default:0" json:"total_requests"`
	SuccessfulRequests  int64      `gorm:"default:0" json:"successful_requests"`
	FailedRequests      int64      `gorm:"default:0" json:"failed_requests"`
	AverageResponseTime float64    `gorm:"default:0" json:"average_response_time"`
	LastRequestAt       *time.Time `json:"last_request_at"`
	ActiveConnections   int64      `gorm:"default:0" json:"active_connections"`
	UptimePercentage    float64    `gorm:"default:0" json:"uptime_percentage"`

	// 关联关系
	Tools     []Tool     `gorm:"foreignKey:McpServerID" json:"tools,omitempty"`
	Resources []Resource `gorm:"foreignKey:McpServerID" json:"resources,omitempty"`
}

// McpServerConfig MCP服务器配置
type McpServerConfig struct {
	TransportType          string                 `json:"transport_type,omitempty"`
	ConnectionParams       map[string]string      `json:"connection_params,omitempty"`
	ConnectionTimeout      int                    `json:"connection_timeout,omitempty"`
	RequestTimeout         int                    `json:"request_timeout,omitempty"`
	MaxRetries             int                    `json:"max_retries,omitempty"`
	SSLEnabled             bool                   `json:"ssl_enabled,omitempty"`
	SSLCertPath            string                 `json:"ssl_cert_path,omitempty"`
	AuthenticationRequired bool                   `json:"authentication_required,omitempty"`
	AuthConfig             map[string]string      `json:"auth_config,omitempty"`
	RateLimit              int                    `json:"rate_limit,omitempty"`
	CustomConfig           map[string]interface{} `json:"custom_config,omitempty"`
}

// SecurityPolicy 安全策略
type SecurityPolicy struct {
	AllowedOrigins      []string          `json:"allowed_origins,omitempty"`
	BlockedOrigins      []string          `json:"blocked_origins,omitempty"`
	RequiredPermissions []string          `json:"required_permissions,omitempty"`
	AuditEnabled        bool              `json:"audit_enabled,omitempty"`
	MaxRequestSize      int               `json:"max_request_size,omitempty"`
	RateLimitingEnabled bool              `json:"rate_limiting_enabled,omitempty"`
	SecurityHeaders     map[string]string `json:"security_headers,omitempty"`
}

// Tool 工具模型
type Tool struct {
	ID            int64            `gorm:"primarykey" json:"id"`
	Name          string           `gorm:"size:255;not null;uniqueIndex:idx_tool_server" json:"name"`
	Description   string           `gorm:"type:text" json:"description"`
	Schema        string           `gorm:"type:longtext" json:"schema"` // JSON schema
	McpServerID   string           `gorm:"size:64;not null;index;uniqueIndex:idx_tool_server" json:"mcp_server_id"`
	Enabled       bool             `gorm:"default:true" json:"enabled"`
	Type          int              `gorm:"default:1" json:"type"`           // 1:function, 2:api, 3:script, 4:service, 5:database, 6:file, 7:network
	Category      int              `gorm:"default:1" json:"category"`       // 1:productivity, 2:development, 3:analysis, 4:communication, 5:content, 6:system, 7:utility
	SecurityLevel int              `gorm:"default:1" json:"security_level"` // 1:public, 2:restricted, 3:private, 4:admin
	Tags          StringSlice      `gorm:"type:json" json:"tags"`
	Config        ToolConfig       `gorm:"type:json" json:"config"`
	Dependencies  ToolDependencies `gorm:"type:json" json:"dependencies"`
	Version       ToolVersion      `gorm:"type:json" json:"version"`
	Metadata      KeyValueMap      `gorm:"type:json" json:"metadata"`
	CreatedAt     time.Time        `json:"created_at"`
	UpdatedAt     time.Time        `json:"updated_at"`
	DeletedAt     gorm.DeletedAt   `gorm:"index" json:"deleted_at"`

	// 使用统计
	TotalCalls      int64      `gorm:"default:0" json:"total_calls"`
	SuccessfulCalls int64      `gorm:"default:0" json:"successful_calls"`
	FailedCalls     int64      `gorm:"default:0" json:"failed_calls"`
	AverageDuration float64    `gorm:"default:0" json:"average_duration"`
	LastCalledAt    *time.Time `json:"last_called_at"`
	SuccessRate     float64    `gorm:"default:0" json:"success_rate"`
	TotalCost       float64    `gorm:"default:0" json:"total_cost"`

	// 关联关系
	McpServer  McpServer       `gorm:"foreignKey:McpServerID" json:"mcp_server,omitempty"`
	Executions []ToolExecution `gorm:"foreignKey:ToolID" json:"executions,omitempty"`
}

// ToolConfig 工具配置
type ToolConfig struct {
	TimeoutSeconds      int                    `json:"timeout_seconds,omitempty"`
	RetryCount          int                    `json:"retry_count,omitempty"`
	CacheEnabled        bool                   `json:"cache_enabled,omitempty"`
	CacheTTLSeconds     int                    `json:"cache_ttl_seconds,omitempty"`
	RateLimitPerMinute  int                    `json:"rate_limit_per_minute,omitempty"`
	AsyncExecution      bool                   `json:"async_execution,omitempty"`
	EnvironmentVars     map[string]string      `json:"environment_vars,omitempty"`
	RequiredPermissions []string               `json:"required_permissions,omitempty"`
	ExecutionContext    string                 `json:"execution_context,omitempty"`
	CustomConfig        map[string]interface{} `json:"custom_config,omitempty"`
}

// ToolDependencies 工具依赖关系
type ToolDependencies []ToolDependency

type ToolDependency struct {
	ToolName          string `json:"tool_name"`
	DependencyType    string `json:"dependency_type"` // prerequisite, optional, conflict
	Required          bool   `json:"required"`
	VersionConstraint string `json:"version_constraint,omitempty"`
}

// ToolVersion 工具版本信息
type ToolVersion struct {
	Version            string    `json:"version"`
	Changelog          string    `json:"changelog,omitempty"`
	Deprecated         bool      `json:"deprecated,omitempty"`
	DeprecationMessage string    `json:"deprecation_message,omitempty"`
	ReleaseDate        time.Time `json:"release_date,omitempty"`
}

// Resource 资源模型
type Resource struct {
	ID           int64               `gorm:"primarykey" json:"id"`
	URI          string              `gorm:"size:500;not null;index" json:"uri"`
	Name         string              `gorm:"size:255;not null" json:"name"`
	Description  string              `gorm:"type:text" json:"description"`
	MimeType     string              `gorm:"size:100" json:"mime_type"`
	McpServerID  string              `gorm:"size:64;not null;index" json:"mcp_server_id"`
	Type         int                 `gorm:"default:1" json:"type"` // 1:file, 2:database, 3:api, 4:stream, 5:memory, 6:cache
	Size         int64               `gorm:"default:0" json:"size"`
	Hash         string              `gorm:"size:64" json:"hash"`
	Cached       bool                `gorm:"default:false" json:"cached"`
	Tags         StringSlice         `gorm:"type:json" json:"tags"`
	Version      string              `gorm:"size:50" json:"version"`
	Permissions  ResourcePermissions `gorm:"type:json" json:"permissions"`
	Metadata     KeyValueMap         `gorm:"type:json" json:"metadata"`
	CreatedAt    time.Time           `json:"created_at"`
	UpdatedAt    time.Time           `json:"updated_at"`
	LastModified *time.Time          `json:"last_modified"`
	DeletedAt    gorm.DeletedAt      `gorm:"index" json:"deleted_at"`

	// 访问统计
	AccessCount           int64      `gorm:"default:0" json:"access_count"`
	LastAccessedAt        *time.Time `json:"last_accessed_at"`
	DownloadCount         int64      `gorm:"default:0" json:"download_count"`
	AverageAccessDuration float64    `gorm:"default:0" json:"average_access_duration"`

	// 关联关系
	McpServer McpServer `gorm:"foreignKey:McpServerID" json:"mcp_server,omitempty"`
}

// ResourcePermissions 资源权限
type ResourcePermissions struct {
	Readable      bool     `json:"readable"`
	Writable      bool     `json:"writable"`
	Executable    bool     `json:"executable"`
	Deletable     bool     `json:"deletable"`
	RequiredRoles []string `json:"required_roles,omitempty"`
	AllowedUsers  []string `json:"allowed_users,omitempty"`
}

// ToolExecution 工具执行记录
type ToolExecution struct {
	ID             string           `gorm:"primarykey;size:64" json:"id"`
	ToolID         int64            `gorm:"not null;index" json:"tool_id"`
	ToolName       string           `gorm:"size:255;not null;index" json:"tool_name"`
	UserID         int64            `gorm:"index" json:"user_id"`
	ConversationID int64            `gorm:"index" json:"conversation_id"`
	Arguments      string           `gorm:"type:longtext" json:"arguments"`
	Result         string           `gorm:"type:longtext" json:"result"`
	Status         int              `gorm:"default:1;index" json:"status"` // 1:pending, 2:running, 3:success, 4:failed, 5:timeout, 6:cancelled
	ErrorMessage   string           `gorm:"type:text" json:"error_message"`
	ExecutionTime  int64            `gorm:"default:0" json:"execution_time"` // milliseconds
	Context        ExecutionContext `gorm:"type:json" json:"context"`
	Metrics        ExecutionMetrics `gorm:"type:json" json:"metrics"`
	Warnings       StringSlice      `gorm:"type:json" json:"warnings"`
	TraceID        string           `gorm:"size:64;index" json:"trace_id"`
	CreatedAt      time.Time        `json:"created_at"`
	UpdatedAt      time.Time        `json:"updated_at"`
	StartedAt      *time.Time       `json:"started_at"`
	CompletedAt    *time.Time       `json:"completed_at"`

	// 关联关系
	Tool Tool `gorm:"foreignKey:ToolID" json:"tool,omitempty"`
}

// ExecutionContext 执行上下文
type ExecutionContext struct {
	Environment      map[string]string      `json:"environment,omitempty"`
	WorkingDirectory string                 `json:"working_directory,omitempty"`
	Permissions      []string               `json:"permissions,omitempty"`
	Variables        map[string]interface{} `json:"variables,omitempty"`
	SessionID        string                 `json:"session_id,omitempty"`
}

// ExecutionMetrics 执行指标
type ExecutionMetrics struct {
	MemoryUsed           int64   `json:"memory_used,omitempty"` // bytes
	CPUUsage             float64 `json:"cpu_usage,omitempty"`   // percentage
	NetworkBytesSent     int64   `json:"network_bytes_sent,omitempty"`
	NetworkBytesReceived int64   `json:"network_bytes_received,omitempty"`
	FileOperations       int     `json:"file_operations,omitempty"`
	Cost                 float64 `json:"cost,omitempty"`
}

// ServerHealthCheck 服务器健康检查
type ServerHealthCheck struct {
	ID          int64       `gorm:"primarykey" json:"id"`
	McpServerID string      `gorm:"size:64;not null;index" json:"mcp_server_id"`
	CheckType   string      `gorm:"size:50;not null" json:"check_type"` // ping, connect, api_test, tool_test
	Status      int         `gorm:"not null;index" json:"status"`       // 1:healthy, 2:warning, 3:critical, 4:unknown
	Message     string      `gorm:"type:text" json:"message"`
	Duration    int64       `gorm:"default:0" json:"duration"` // milliseconds
	Details     KeyValueMap `gorm:"type:json" json:"details"`
	CreatedAt   time.Time   `json:"created_at"`

	// 关联关系
	McpServer McpServer `gorm:"foreignKey:McpServerID" json:"mcp_server,omitempty"`
}

// ToolAuditLog 工具审计日志
type ToolAuditLog struct {
	ID           int64     `gorm:"primarykey" json:"id"`
	ToolID       int64     `gorm:"not null;index" json:"tool_id"`
	UserID       int64     `gorm:"not null;index" json:"user_id"`
	Action       string    `gorm:"size:50;not null" json:"action"` // call, enable, disable, configure
	Details      string    `gorm:"type:text" json:"details"`
	IPAddress    string    `gorm:"size:45" json:"ip_address"`
	UserAgent    string    `gorm:"size:500" json:"user_agent"`
	Success      bool      `gorm:"not null" json:"success"`
	ErrorMessage string    `gorm:"type:text" json:"error_message"`
	CreatedAt    time.Time `json:"created_at"`

	// 关联关系
	Tool Tool `gorm:"foreignKey:ToolID" json:"tool,omitempty"`
}

// 自定义GORM类型转换器实现
func (c McpServerConfig) Value() (interface{}, error) {
	b, err := json.Marshal(c)
	return string(b), err
}

func (c *McpServerConfig) Scan(value interface{}) error {
	if value == nil {
		return nil
	}

	var bytes []byte
	switch v := value.(type) {
	case []byte:
		bytes = v
	case string:
		bytes = []byte(v)
	default:
		return nil
	}

	return json.Unmarshal(bytes, c)
}

func (s SecurityPolicy) Value() (interface{}, error) {
	b, err := json.Marshal(s)
	return string(b), err
}

func (s *SecurityPolicy) Scan(value interface{}) error {
	if value == nil {
		return nil
	}

	var bytes []byte
	switch v := value.(type) {
	case []byte:
		bytes = v
	case string:
		bytes = []byte(v)
	default:
		return nil
	}

	return json.Unmarshal(bytes, s)
}

func (c ToolConfig) Value() (interface{}, error) {
	b, err := json.Marshal(c)
	return string(b), err
}

func (c *ToolConfig) Scan(value interface{}) error {
	if value == nil {
		return nil
	}

	var bytes []byte
	switch v := value.(type) {
	case []byte:
		bytes = v
	case string:
		bytes = []byte(v)
	default:
		return nil
	}

	return json.Unmarshal(bytes, c)
}

func (d ToolDependencies) Value() (interface{}, error) {
	b, err := json.Marshal(d)
	return string(b), err
}

func (d *ToolDependencies) Scan(value interface{}) error {
	if value == nil {
		*d = ToolDependencies{}
		return nil
	}

	var bytes []byte
	switch v := value.(type) {
	case []byte:
		bytes = v
	case string:
		bytes = []byte(v)
	default:
		return nil
	}

	return json.Unmarshal(bytes, d)
}

func (v ToolVersion) Value() (interface{}, error) {
	b, err := json.Marshal(v)
	return string(b), err
}

func (v *ToolVersion) Scan(value interface{}) error {
	if value == nil {
		return nil
	}

	var bytes []byte
	switch val := value.(type) {
	case []byte:
		bytes = val
	case string:
		bytes = []byte(val)
	default:
		return nil
	}

	return json.Unmarshal(bytes, v)
}

func (p ResourcePermissions) Value() (interface{}, error) {
	b, err := json.Marshal(p)
	return string(b), err
}

func (p *ResourcePermissions) Scan(value interface{}) error {
	if value == nil {
		return nil
	}

	var bytes []byte
	switch v := value.(type) {
	case []byte:
		bytes = v
	case string:
		bytes = []byte(v)
	default:
		return nil
	}

	return json.Unmarshal(bytes, p)
}

func (c ExecutionContext) Value() (interface{}, error) {
	b, err := json.Marshal(c)
	return string(b), err
}

func (c *ExecutionContext) Scan(value interface{}) error {
	if value == nil {
		return nil
	}

	var bytes []byte
	switch v := value.(type) {
	case []byte:
		bytes = v
	case string:
		bytes = []byte(v)
	default:
		return nil
	}

	return json.Unmarshal(bytes, c)
}

func (m ExecutionMetrics) Value() (interface{}, error) {
	b, err := json.Marshal(m)
	return string(b), err
}

func (m *ExecutionMetrics) Scan(value interface{}) error {
	if value == nil {
		return nil
	}

	var bytes []byte
	switch v := value.(type) {
	case []byte:
		bytes = v
	case string:
		bytes = []byte(v)
	default:
		return nil
	}

	return json.Unmarshal(bytes, m)
}

// TableName 返回表名
func (McpServer) TableName() string {
	return "mcp_servers"
}

func (Tool) TableName() string {
	return "tools"
}

func (Resource) TableName() string {
	return "resources"
}

func (ToolExecution) TableName() string {
	return "tool_executions"
}

func (ServerHealthCheck) TableName() string {
	return "server_health_checks"
}

func (ToolAuditLog) TableName() string {
	return "tool_audit_logs"
}
