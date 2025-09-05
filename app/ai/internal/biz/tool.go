package biz

import (
	"context"
	"time"

	"universal/app/ai/internal/data/model"

	"github.com/go-kratos/kratos/v2/log"
)

// ToolUsecase 工具业务逻辑
type ToolUsecase struct {
	repo   ToolRepo
	logger *log.Helper
}

// ToolRepo 工具仓库接口
type ToolRepo interface {
	// MCP服务器管理
	CreateMcpServer(ctx context.Context, server *model.McpServer) (*model.McpServer, error)
	GetMcpServer(ctx context.Context, id string) (*model.McpServer, error)
	UpdateMcpServer(ctx context.Context, server *model.McpServer) (*model.McpServer, error)
	DeleteMcpServer(ctx context.Context, id string, forceDelete bool) error
	ListMcpServers(ctx context.Context, page, pageSize int32, filters McpServerFilter) ([]*model.McpServer, int64, error)
	TestMcpServer(ctx context.Context, id string, testCases []string) ([]TestResult, error)

	// 工具管理
	GetTool(ctx context.Context, name string) (*model.Tool, error)
	ListTools(ctx context.Context, page, pageSize int32, filters ToolFilter) ([]*model.Tool, int64, error)
	EnableTool(ctx context.Context, toolName string) error
	DisableTool(ctx context.Context, toolName string) error
	ConfigureTool(ctx context.Context, toolName string, config model.ToolConfig) error
	GetToolConfig(ctx context.Context, toolName string) (*model.ToolConfig, error)
	SyncServerTools(ctx context.Context, serverID string) ([]model.Tool, error)

	// 工具执行
	CreateToolExecution(ctx context.Context, execution *model.ToolExecution) (*model.ToolExecution, error)
	UpdateToolExecution(ctx context.Context, execution *model.ToolExecution) error
	GetToolExecution(ctx context.Context, id string) (*model.ToolExecution, error)
	ListToolExecutions(ctx context.Context, page, pageSize int32, filters ToolExecutionFilter) ([]*model.ToolExecution, int64, error)
	GetToolExecutionStats(ctx context.Context, toolName string, startTime, endTime time.Time, groupBy string) (*ToolExecutionStats, error)

	// 资源管理
	GetResource(ctx context.Context, uri string) (*model.Resource, error)
	ListResources(ctx context.Context, page, pageSize int32, filters ResourceFilter) ([]*model.Resource, int64, error)
	SearchResources(ctx context.Context, query string, filters map[string]string, limit int32) ([]*ResourceSearchResult, error)
	SyncServerResources(ctx context.Context, serverID string) ([]model.Resource, error)

	// 健康检查
	CreateHealthCheck(ctx context.Context, check *model.ServerHealthCheck) error
	GetServerHealthStatus(ctx context.Context, serverID string) (*ServerHealthStatus, error)
	ListHealthChecks(ctx context.Context, serverID string, limit int32) ([]*model.ServerHealthCheck, error)

	// 审计日志
	CreateAuditLog(ctx context.Context, log *model.ToolAuditLog) error
	ListAuditLogs(ctx context.Context, page, pageSize int32, filters AuditLogFilter) ([]*model.ToolAuditLog, int64, error)
}

// McpServerFilter MCP服务器过滤条件
type McpServerFilter struct {
	Status       int32
	Tags         []string
	IncludeStats bool
}

// ToolFilter 工具过滤条件
type ToolFilter struct {
	McpServer        string
	OnlyEnabled      bool
	Type             int32
	Category         int32
	Tags             []string
	MaxSecurityLevel int32
}

// ResourceFilter 资源过滤条件
type ResourceFilter struct {
	McpServer string
	MimeType  string
	Type      int32
	Tags      []string
}

// ToolExecutionFilter 工具执行过滤条件
type ToolExecutionFilter struct {
	ToolName       string
	UserID         int64
	ConversationID int64
	Status         int32
	StartTime      *time.Time
	EndTime        *time.Time
}

// AuditLogFilter 审计日志过滤条件
type AuditLogFilter struct {
	ToolID    int64
	UserID    int64
	Action    string
	Success   *bool
	StartTime *time.Time
	EndTime   *time.Time
}

// TestResult 测试结果
type TestResult struct {
	TestCase string
	Passed   bool
	Message  string
	Duration time.Duration
}

// ToolExecutionStats 工具执行统计
type ToolExecutionStats struct {
	TotalExecutions    int64
	AverageDuration    float64
	SuccessRate        float64
	TotalCost          float64
	StatusDistribution map[string]int64
	TimeSeriesData     []TimeSeriesPoint
}

// TimeSeriesPoint 时序数据点
type TimeSeriesPoint struct {
	Timestamp time.Time
	Value     float64
	Label     string
}

// ResourceSearchResult 资源搜索结果
type ResourceSearchResult struct {
	Resource       *model.Resource
	RelevanceScore float64
	MatchedFields  []string
}

// ServerHealthStatus 服务器健康状态
type ServerHealthStatus struct {
	ServerID      string
	OverallStatus int // 1:healthy, 2:warning, 3:critical, 4:unknown
	Message       string
	LastCheckTime time.Time
	Checks        []HealthCheck
}

// HealthCheck 健康检查详情
type HealthCheck struct {
	Type     string
	Status   int
	Message  string
	Duration time.Duration
}

// ToolCallRequest 工具调用请求
type ToolCallRequest struct {
	Name           string
	Arguments      string
	ConversationID int64
	UserID         int64
	Context        model.ExecutionContext
	TimeoutSeconds int32
	Async          bool
	TraceID        string
}

// ToolCallResponse 工具调用响应
type ToolCallResponse struct {
	ExecutionID  string
	Result       string
	Status       int
	ErrorMessage string
	Metadata     map[string]string
	Metrics      model.ExecutionMetrics
	Warnings     []string
}

// BatchToolCallRequest 批量工具调用请求
type BatchToolCallRequest struct {
	ToolCalls      []BatchToolCall
	Parallel       bool
	MaxConcurrency int32
	StopOnError    bool
}

// BatchToolCall 批量工具调用
type BatchToolCall struct {
	ID             string
	ToolName       string
	Arguments      string
	DependsOn      []string
	TimeoutSeconds int32
}

// BatchToolCallResponse 批量工具调用响应
type BatchToolCallResponse struct {
	Results      []BatchToolResult
	AllSuccess   bool
	SuccessCount int32
	FailedCount  int32
}

// BatchToolResult 批量工具结果
type BatchToolResult struct {
	ID       string
	Response ToolCallResponse
}

// NewToolUsecase 创建工具业务逻辑实例
func NewToolUsecase(repo ToolRepo, logger log.Logger) *ToolUsecase {
	return &ToolUsecase{
		repo:   repo,
		logger: log.NewHelper(logger),
	}
}

// RegisterMcpServer 注册MCP服务器
func (uc *ToolUsecase) RegisterMcpServer(ctx context.Context, name, description, endpoint string, config model.McpServerConfig, metadata map[string]string, tags []string) (*model.McpServer, error) {
	now := time.Now()

	server := &model.McpServer{
		ID:          uc.generateServerID(),
		Name:        name,
		Description: description,
		Endpoint:    endpoint,
		Status:      1, // active
		Config:      config,
		Metadata:    model.KeyValueMap(metadata),
		Tags:        model.StringSlice(tags),
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	createdServer, err := uc.repo.CreateMcpServer(ctx, server)
	if err != nil {
		return nil, err
	}

	// 同步服务器的工具和资源
	go func() {
		if _, err := uc.repo.SyncServerTools(context.Background(), server.ID); err != nil {
			uc.logger.Warnw("failed to sync server tools", "server_id", server.ID, "error", err)
		}
		if _, err := uc.repo.SyncServerResources(context.Background(), server.ID); err != nil {
			uc.logger.Warnw("failed to sync server resources", "server_id", server.ID, "error", err)
		}
	}()

	return createdServer, nil
}

// GetMcpServer 获取MCP服务器
func (uc *ToolUsecase) GetMcpServer(ctx context.Context, id string, includeStats, includeHealth bool) (*model.McpServer, error) {
	server, err := uc.repo.GetMcpServer(ctx, id)
	if err != nil {
		return nil, err
	}

	// TODO: 根据includeStats和includeHealth参数加载额外信息

	return server, nil
}

// ListMcpServers 列出MCP服务器
func (uc *ToolUsecase) ListMcpServers(ctx context.Context, page, pageSize int32, statusFilter int32, tags []string, includeStats bool) ([]*model.McpServer, int64, error) {
	filter := McpServerFilter{
		Status:       statusFilter,
		Tags:         tags,
		IncludeStats: includeStats,
	}

	return uc.repo.ListMcpServers(ctx, page, pageSize, filter)
}

// TestMcpServer 测试MCP服务器
func (uc *ToolUsecase) TestMcpServer(ctx context.Context, id string, testCases []string) ([]TestResult, *ServerHealthStatus, error) {
	results, err := uc.repo.TestMcpServer(ctx, id, testCases)
	if err != nil {
		return nil, nil, err
	}

	// 创建健康检查记录
	overallStatus := 1 // healthy
	for _, result := range results {
		if !result.Passed {
			overallStatus = 3 // critical
			break
		}
	}

	healthCheck := &model.ServerHealthCheck{
		McpServerID: id,
		CheckType:   "api_test",
		Status:      overallStatus,
		Message:     "Automated server test",
		Duration:    0, // TODO: 计算总耗时
		CreatedAt:   time.Now(),
	}

	err = uc.repo.CreateHealthCheck(ctx, healthCheck)
	if err != nil {
		uc.logger.Warnw("failed to create health check record", "server_id", id, "error", err)
	}

	// 获取健康状态
	health, err := uc.repo.GetServerHealthStatus(ctx, id)
	if err != nil {
		uc.logger.Warnw("failed to get server health status", "server_id", id, "error", err)
	}

	return results, health, nil
}

// ListTools 列出工具
func (uc *ToolUsecase) ListTools(ctx context.Context, page, pageSize int32, mcpServer string, onlyEnabled bool, typeFilter, categoryFilter int32, tags []string, maxSecurityLevel int32) ([]*model.Tool, int64, error) {
	filter := ToolFilter{
		McpServer:        mcpServer,
		OnlyEnabled:      onlyEnabled,
		Type:             typeFilter,
		Category:         categoryFilter,
		Tags:             tags,
		MaxSecurityLevel: maxSecurityLevel,
	}

	return uc.repo.ListTools(ctx, page, pageSize, filter)
}

// GetTool 获取工具详情
func (uc *ToolUsecase) GetTool(ctx context.Context, name string, includeStats, includeConfig bool) (*model.Tool, error) {
	tool, err := uc.repo.GetTool(ctx, name)
	if err != nil {
		return nil, err
	}

	// TODO: 根据参数加载额外信息

	return tool, nil
}

// CallTool 调用工具
func (uc *ToolUsecase) CallTool(ctx context.Context, req ToolCallRequest) (*ToolCallResponse, error) {
	// 获取工具信息
	tool, err := uc.repo.GetTool(ctx, req.Name)
	if err != nil {
		return nil, err
	}

	if !tool.Enabled {
		return &ToolCallResponse{
			Status:       4, // failed
			ErrorMessage: "tool is disabled",
		}, nil
	}

	// 创建执行记录
	execution := &model.ToolExecution{
		ID:             uc.generateExecutionID(),
		ToolID:         tool.ID,
		ToolName:       req.Name,
		UserID:         req.UserID,
		ConversationID: req.ConversationID,
		Arguments:      req.Arguments,
		Status:         1, // pending
		Context:        req.Context,
		TraceID:        req.TraceID,
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}

	execution, err = uc.repo.CreateToolExecution(ctx, execution)
	if err != nil {
		return nil, err
	}

	// 如果是异步执行
	if req.Async {
		// 启动异步执行
		go uc.executeToolAsync(context.Background(), execution)

		return &ToolCallResponse{
			ExecutionID: execution.ID,
			Status:      2, // running
		}, nil
	}

	// 同步执行
	return uc.executeTool(ctx, execution)
}

// BatchCallTools 批量调用工具
func (uc *ToolUsecase) BatchCallTools(ctx context.Context, req BatchToolCallRequest) (*BatchToolCallResponse, error) {
	results := make([]BatchToolResult, len(req.ToolCalls))
	successCount := int32(0)
	failedCount := int32(0)

	// TODO: 实现批量调用逻辑
	// 1. 解析依赖关系
	// 2. 根据并行设置和依赖关系安排执行顺序
	// 3. 执行工具调用

	for i, toolCall := range req.ToolCalls {
		// 简化实现：顺序执行
		callReq := ToolCallRequest{
			Name:           toolCall.ToolName,
			Arguments:      toolCall.Arguments,
			TimeoutSeconds: toolCall.TimeoutSeconds,
		}

		response, err := uc.CallTool(ctx, callReq)
		if err != nil {
			response = &ToolCallResponse{
				Status:       4, // failed
				ErrorMessage: err.Error(),
			}
		}

		results[i] = BatchToolResult{
			ID:       toolCall.ID,
			Response: *response,
		}

		if response.Status == 3 { // success
			successCount++
		} else {
			failedCount++
			if req.StopOnError {
				break
			}
		}
	}

	return &BatchToolCallResponse{
		Results:      results,
		AllSuccess:   failedCount == 0,
		SuccessCount: successCount,
		FailedCount:  failedCount,
	}, nil
}

// GetToolExecutionHistory 获取工具执行历史
func (uc *ToolUsecase) GetToolExecutionHistory(ctx context.Context, page, pageSize int32, toolName string, userID, conversationID int64, status int32, startTime, endTime *time.Time) ([]*model.ToolExecution, int64, error) {
	filter := ToolExecutionFilter{
		ToolName:       toolName,
		UserID:         userID,
		ConversationID: conversationID,
		Status:         status,
		StartTime:      startTime,
		EndTime:        endTime,
	}

	return uc.repo.ListToolExecutions(ctx, page, pageSize, filter)
}

// GetToolExecutionStats 获取工具执行统计
func (uc *ToolUsecase) GetToolExecutionStats(ctx context.Context, toolName string, startTime, endTime time.Time, groupBy string) (*ToolExecutionStats, error) {
	return uc.repo.GetToolExecutionStats(ctx, toolName, startTime, endTime, groupBy)
}

// 私有方法
func (uc *ToolUsecase) executeTool(ctx context.Context, execution *model.ToolExecution) (*ToolCallResponse, error) {
	startTime := time.Now()

	// 更新状态为运行中
	execution.Status = 2 // running
	execution.StartedAt = &startTime
	execution.UpdatedAt = time.Now()

	err := uc.repo.UpdateToolExecution(ctx, execution)
	if err != nil {
		uc.logger.Warnw("failed to update execution status", "execution_id", execution.ID, "error", err)
	}

	// TODO: 实现实际的工具调用逻辑
	// 这里应该：
	// 1. 根据工具配置连接到MCP服务器
	// 2. 发送工具调用请求
	// 3. 获取结果
	// 4. 处理错误和超时

	// 模拟执行
	time.Sleep(100 * time.Millisecond)

	// 模拟成功结果
	endTime := time.Now()
	duration := endTime.Sub(startTime)

	execution.Status = 3 // success
	execution.Result = `{"success": true, "message": "Tool executed successfully"}`
	execution.CompletedAt = &endTime
	execution.ExecutionTime = duration.Milliseconds()
	execution.UpdatedAt = time.Now()

	// 模拟指标
	execution.Metrics = model.ExecutionMetrics{
		MemoryUsed: 1024 * 1024, // 1MB
		CPUUsage:   0.1,         // 10%
		Cost:       0.001,       // $0.001
	}

	err = uc.repo.UpdateToolExecution(ctx, execution)
	if err != nil {
		uc.logger.Warnw("failed to update execution result", "execution_id", execution.ID, "error", err)
	}

	return &ToolCallResponse{
		ExecutionID: execution.ID,
		Result:      execution.Result,
		Status:      execution.Status,
		Metrics:     execution.Metrics,
	}, nil
}

func (uc *ToolUsecase) executeToolAsync(ctx context.Context, execution *model.ToolExecution) {
	_, err := uc.executeTool(ctx, execution)
	if err != nil {
		uc.logger.Errorw("async tool execution failed", "execution_id", execution.ID, "error", err)
	}
}

func (uc *ToolUsecase) generateServerID() string {
	// TODO: 生成唯一的服务器ID
	return "server_" + time.Now().Format("20060102150405")
}

func (uc *ToolUsecase) generateExecutionID() string {
	// TODO: 生成唯一的执行ID
	return "exec_" + time.Now().Format("20060102150405000")
}
