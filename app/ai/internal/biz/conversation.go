package biz

import (
	"context"
	"time"

	"universal/app/ai/internal/data/model"

	"github.com/go-kratos/kratos/v2/log"
)

// ConversationUsecase 对话业务逻辑
type ConversationUsecase struct {
	repo   ConversationRepo
	logger *log.Helper
}

// ConversationRepo 对话仓库接口
type ConversationRepo interface {
	// 对话管理
	CreateConversation(ctx context.Context, conversation *model.Conversation) (*model.Conversation, error)
	GetConversation(ctx context.Context, id int64) (*model.Conversation, error)
	UpdateConversation(ctx context.Context, conversation *model.Conversation) (*model.Conversation, error)
	DeleteConversation(ctx context.Context, id int64, hardDelete bool) error
	ListConversations(ctx context.Context, userID int64, page, pageSize int32, filters ConversationFilter) ([]*model.Conversation, int64, error)
	ArchiveConversation(ctx context.Context, id int64) error
	RestoreConversation(ctx context.Context, id int64) error

	// 消息管理
	CreateMessage(ctx context.Context, message *model.Message) (*model.Message, error)
	GetMessage(ctx context.Context, id int64) (*model.Message, error)
	UpdateMessage(ctx context.Context, message *model.Message) (*model.Message, error)
	DeleteMessage(ctx context.Context, messageIDs []int64, hardDelete bool) error
	ListMessages(ctx context.Context, conversationID int64, page, pageSize int32, filters MessageFilter) ([]*model.Message, int64, error)

	// 工具调用
	CreateToolCall(ctx context.Context, toolCall *model.ToolCall) (*model.ToolCall, error)
	UpdateToolCall(ctx context.Context, toolCall *model.ToolCall) (*model.ToolCall, error)

	// 记忆管理
	GetConversationMemory(ctx context.Context, conversationID int64) (*model.ConversationMemory, error)
	SetConversationMemory(ctx context.Context, memory *model.ConversationMemory) error

	// 统计信息
	UpdateConversationStats(ctx context.Context, conversationID int64, inputTokens, outputTokens int64, duration time.Duration) error
	GetConversationStats(ctx context.Context, conversationID int64) (*ConversationStats, error)
}

// ConversationFilter 对话过滤条件
type ConversationFilter struct {
	Keyword  string
	Status   int32
	Tags     []string
	SortBy   string
	SortDesc bool
}

// MessageFilter 消息过滤条件
type MessageFilter struct {
	Role               string
	Status             int32
	IncludeToolCalls   bool
	IncludeAttachments bool
	IncludeMetrics     bool
}

// ConversationStats 对话统计信息
type ConversationStats struct {
	MessageCount          int64
	UserMessageCount      int64
	AssistantMessageCount int64
	TotalInputTokens      int64
	TotalOutputTokens     int64
	ToolCallCount         int64
	TotalDuration         time.Duration
	AverageResponseTime   float64
	LastMessageAt         *time.Time
}

// NewConversationUsecase 创建对话业务逻辑实例
func NewConversationUsecase(repo ConversationRepo, logger log.Logger) *ConversationUsecase {
	return &ConversationUsecase{
		repo:   repo,
		logger: log.NewHelper(logger),
	}
}

// CreateConversation 创建对话
func (uc *ConversationUsecase) CreateConversation(ctx context.Context, userID int64, title, modelName, systemPrompt string, config model.ConversationConfig, description string, tags []string, priority int32) (*model.Conversation, error) {
	now := time.Now()
	conversation := &model.Conversation{
		UserID:       userID,
		Title:        title,
		ModelName:    modelName,
		SystemPrompt: systemPrompt,
		Config:       config,
		Status:       1, // active
		Description:  description,
		Tags:         model.StringSlice(tags),
		Priority:     int(priority),
		CreatedAt:    now,
		UpdatedAt:    now,
		LastActiveAt: now,
	}

	return uc.repo.CreateConversation(ctx, conversation)
}

// GetConversation 获取对话
func (uc *ConversationUsecase) GetConversation(ctx context.Context, id int64) (*model.Conversation, error) {
	return uc.repo.GetConversation(ctx, id)
}

// UpdateConversation 更新对话
func (uc *ConversationUsecase) UpdateConversation(ctx context.Context, id int64, title, systemPrompt, description string, config model.ConversationConfig, tags []string, priority int32) (*model.Conversation, error) {
	conversation, err := uc.repo.GetConversation(ctx, id)
	if err != nil {
		return nil, err
	}

	if title != "" {
		conversation.Title = title
	}
	if systemPrompt != "" {
		conversation.SystemPrompt = systemPrompt
	}
	if description != "" {
		conversation.Description = description
	}
	if len(config.CustomParams) > 0 || config.Temperature != nil {
		conversation.Config = config
	}
	if len(tags) > 0 {
		conversation.Tags = model.StringSlice(tags)
	}
	if priority > 0 {
		conversation.Priority = int(priority)
	}

	conversation.UpdatedAt = time.Now()

	return uc.repo.UpdateConversation(ctx, conversation)
}

// DeleteConversation 删除对话
func (uc *ConversationUsecase) DeleteConversation(ctx context.Context, id int64, hardDelete bool) error {
	return uc.repo.DeleteConversation(ctx, id, hardDelete)
}

// ListConversations 获取对话列表
func (uc *ConversationUsecase) ListConversations(ctx context.Context, userID int64, page, pageSize int32, keyword string, status int32, tags []string, sortBy string, sortDesc bool) ([]*model.Conversation, int64, error) {
	filter := ConversationFilter{
		Keyword:  keyword,
		Status:   status,
		Tags:     tags,
		SortBy:   sortBy,
		SortDesc: sortDesc,
	}

	return uc.repo.ListConversations(ctx, userID, page, pageSize, filter)
}

// ArchiveConversation 归档对话
func (uc *ConversationUsecase) ArchiveConversation(ctx context.Context, id int64) error {
	return uc.repo.ArchiveConversation(ctx, id)
}

// RestoreConversation 恢复对话
func (uc *ConversationUsecase) RestoreConversation(ctx context.Context, id int64) error {
	return uc.repo.RestoreConversation(ctx, id)
}

// SendMessage 发送消息
func (uc *ConversationUsecase) SendMessage(ctx context.Context, conversationID int64, content string, attachments []MessageAttachmentInfo, enableTools bool, allowedTools []string, options map[string]string, parentMessageID *int64) (*model.Message, *model.Message, error) {
	// 创建用户消息
	userMessage := &model.Message{
		ConversationID: conversationID,
		Role:           "user",
		Content:        content,
		Status:         1, // pending
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}

	if parentMessageID != nil {
		userMessage.ParentMessageID = parentMessageID
	}

	userMessage, err := uc.repo.CreateMessage(ctx, userMessage)
	if err != nil {
		return nil, nil, err
	}

	// TODO: 实现AI回复逻辑
	// 这里应该调用AI服务生成回复
	assistantMessage := &model.Message{
		ConversationID: conversationID,
		Role:           "assistant",
		Content:        "这是一个模拟回复", // 临时内容
		Status:         3,          // completed
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
		ResponseTime:   1.5, // 模拟响应时间
		InputTokens:    100,
		OutputTokens:   50,
	}

	assistantMessage, err = uc.repo.CreateMessage(ctx, assistantMessage)
	if err != nil {
		return userMessage, nil, err
	}

	// 更新对话统计信息
	err = uc.repo.UpdateConversationStats(ctx, conversationID, int64(assistantMessage.InputTokens), int64(assistantMessage.OutputTokens), time.Duration(assistantMessage.ResponseTime*float64(time.Second)))
	if err != nil {
		uc.logger.Warnw("failed to update conversation stats", "error", err)
	}

	return userMessage, assistantMessage, nil
}

// GetMessages 获取消息列表
func (uc *ConversationUsecase) GetMessages(ctx context.Context, conversationID int64, page, pageSize int32, includeToolCalls, includeAttachments, includeMetrics bool, roleFilter, statusFilter string) ([]*model.Message, int64, error) {
	filter := MessageFilter{
		Role:               roleFilter,
		IncludeToolCalls:   includeToolCalls,
		IncludeAttachments: includeAttachments,
		IncludeMetrics:     includeMetrics,
	}

	// 转换状态过滤器
	switch statusFilter {
	case "pending":
		filter.Status = 1
	case "processing":
		filter.Status = 2
	case "completed":
		filter.Status = 3
	case "failed":
		filter.Status = 4
	case "deleted":
		filter.Status = 5
	}

	return uc.repo.ListMessages(ctx, conversationID, page, pageSize, filter)
}

// DeleteMessage 删除消息
func (uc *ConversationUsecase) DeleteMessage(ctx context.Context, messageIDs []int64, hardDelete bool) error {
	return uc.repo.DeleteMessage(ctx, messageIDs, hardDelete)
}

// RegenerateMessage 重新生成消息
func (uc *ConversationUsecase) RegenerateMessage(ctx context.Context, messageID int64, options map[string]string) (*model.Message, error) {
	originalMessage, err := uc.repo.GetMessage(ctx, messageID)
	if err != nil {
		return nil, err
	}

	// TODO: 实现重新生成逻辑
	// 这里应该调用AI服务重新生成回复
	newMessage := &model.Message{
		ConversationID:  originalMessage.ConversationID,
		Role:            originalMessage.Role,
		Content:         "这是重新生成的回复", // 临时内容
		Status:          3,           // completed
		ParentMessageID: originalMessage.ParentMessageID,
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
		ResponseTime:    1.2,
		InputTokens:     120,
		OutputTokens:    60,
	}

	return uc.repo.CreateMessage(ctx, newMessage)
}

// GetConversationMemory 获取对话记忆
func (uc *ConversationUsecase) GetConversationMemory(ctx context.Context, conversationID int64) (*model.ConversationMemory, error) {
	return uc.repo.GetConversationMemory(ctx, conversationID)
}

// SetConversationMemory 设置对话记忆
func (uc *ConversationUsecase) SetConversationMemory(ctx context.Context, conversationID int64, summary string, keyPoints []string, userPreferences map[string]string, importantFacts []string) error {
	memory := &model.ConversationMemory{
		ConversationID:  conversationID,
		Summary:         summary,
		KeyPoints:       model.StringSlice(keyPoints),
		UserPreferences: model.KeyValueMap(userPreferences),
		ImportantFacts:  model.StringSlice(importantFacts),
		MemoryVersion:   1,
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
	}

	return uc.repo.SetConversationMemory(ctx, memory)
}

// GetConversationStats 获取对话统计信息
func (uc *ConversationUsecase) GetConversationStats(ctx context.Context, conversationID int64) (*ConversationStats, error) {
	return uc.repo.GetConversationStats(ctx, conversationID)
}

// ClearConversationHistory 清空对话历史
func (uc *ConversationUsecase) ClearConversationHistory(ctx context.Context, conversationID int64, keepSystemMessages bool) error {
	// TODO: 实现清空历史逻辑
	// 需要删除所有非系统消息（如果keepSystemMessages为true）
	return nil
}

// SummarizeConversation 总结对话
func (uc *ConversationUsecase) SummarizeConversation(ctx context.Context, conversationID int64, maxMessages int32, summaryStyle string) (string, []string, error) {
	// TODO: 实现对话总结逻辑
	// 这里应该调用AI服务生成对话总结
	summary := "这是对话的总结"
	keyPoints := []string{"要点1", "要点2", "要点3"}

	return summary, keyPoints, nil
}

// MessageAttachmentInfo 消息附件信息
type MessageAttachmentInfo struct {
	Name     string
	URL      string
	MimeType string
	Size     int64
	Type     int
	Metadata map[string]string
}
