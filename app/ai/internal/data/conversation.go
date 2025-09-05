package data

import (
	"context"
	"fmt"
	"time"

	"universal/app/ai/internal/biz"
	"universal/app/ai/internal/data/model"

	"github.com/go-kratos/kratos/v2/log"
	"gorm.io/gorm"
)

type conversationRepo struct {
	data   *Data
	logger *log.Helper
}

// NewConversationRepo 创建对话仓库实例
func NewConversationRepo(data *Data, logger log.Logger) biz.ConversationRepo {
	return &conversationRepo{
		data:   data,
		logger: log.NewHelper(logger),
	}
}

// CreateConversation 创建对话
func (r *conversationRepo) CreateConversation(ctx context.Context, conversation *model.Conversation) (*model.Conversation, error) {
	if err := r.data.db.WithContext(ctx).Create(conversation).Error; err != nil {
		return nil, err
	}
	return conversation, nil
}

// GetConversation 获取对话
func (r *conversationRepo) GetConversation(ctx context.Context, id int64) (*model.Conversation, error) {
	var conversation model.Conversation
	err := r.data.db.WithContext(ctx).Where("id = ?", id).First(&conversation).Error
	if err != nil {
		return nil, err
	}
	return &conversation, nil
}

// UpdateConversation 更新对话
func (r *conversationRepo) UpdateConversation(ctx context.Context, conversation *model.Conversation) (*model.Conversation, error) {
	if err := r.data.db.WithContext(ctx).Save(conversation).Error; err != nil {
		return nil, err
	}
	return conversation, nil
}

// DeleteConversation 删除对话
func (r *conversationRepo) DeleteConversation(ctx context.Context, id int64, hardDelete bool) error {
	if hardDelete {
		// 硬删除：彻底删除记录
		return r.data.db.WithContext(ctx).Unscoped().Delete(&model.Conversation{}, id).Error
	} else {
		// 软删除：更新状态
		return r.data.db.WithContext(ctx).Model(&model.Conversation{}).Where("id = ?", id).Update("status", 3).Error
	}
}

// ListConversations 获取对话列表
func (r *conversationRepo) ListConversations(ctx context.Context, userID int64, page, pageSize int32, filters biz.ConversationFilter) ([]*model.Conversation, int64, error) {
	var conversations []*model.Conversation
	var total int64

	query := r.data.db.WithContext(ctx).Model(&model.Conversation{}).Where("user_id = ?", userID)

	// 应用过滤条件
	if filters.Keyword != "" {
		query = query.Where("title LIKE ? OR description LIKE ?",
			fmt.Sprintf("%%%s%%", filters.Keyword),
			fmt.Sprintf("%%%s%%", filters.Keyword))
	}

	if filters.Status > 0 {
		query = query.Where("status = ?", filters.Status)
	}

	if len(filters.Tags) > 0 {
		// 标签过滤 - 使用JSON查询
		for _, tag := range filters.Tags {
			query = query.Where("JSON_CONTAINS(tags, ?)", fmt.Sprintf(`"%s"`, tag))
		}
	}

	// 获取总数
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 排序
	orderBy := "created_at DESC"
	if filters.SortBy != "" {
		direction := "ASC"
		if filters.SortDesc {
			direction = "DESC"
		}
		orderBy = fmt.Sprintf("%s %s", filters.SortBy, direction)
	}
	query = query.Order(orderBy)

	// 分页
	offset := (page - 1) * pageSize
	if err := query.Offset(int(offset)).Limit(int(pageSize)).Find(&conversations).Error; err != nil {
		return nil, 0, err
	}

	return conversations, total, nil
}

// ArchiveConversation 归档对话
func (r *conversationRepo) ArchiveConversation(ctx context.Context, id int64) error {
	return r.data.db.WithContext(ctx).Model(&model.Conversation{}).Where("id = ?", id).Update("status", 2).Error
}

// RestoreConversation 恢复对话
func (r *conversationRepo) RestoreConversation(ctx context.Context, id int64) error {
	return r.data.db.WithContext(ctx).Model(&model.Conversation{}).Where("id = ?", id).Update("status", 1).Error
}

// CreateMessage 创建消息
func (r *conversationRepo) CreateMessage(ctx context.Context, message *model.Message) (*model.Message, error) {
	return r.createMessageWithTx(ctx, r.data.db, message)
}

func (r *conversationRepo) createMessageWithTx(ctx context.Context, tx *gorm.DB, message *model.Message) (*model.Message, error) {
	if err := tx.WithContext(ctx).Create(message).Error; err != nil {
		return nil, err
	}

	// 更新对话的最后活跃时间和消息数量
	updates := map[string]interface{}{
		"last_active_at": time.Now(),
		"message_count":  gorm.Expr("message_count + 1"),
		"updated_at":     time.Now(),
	}

	if err := tx.WithContext(ctx).Model(&model.Conversation{}).Where("id = ?", message.ConversationID).Updates(updates).Error; err != nil {
		r.logger.Warnw("failed to update conversation stats", "error", err)
	}

	return message, nil
}

// GetMessage 获取消息
func (r *conversationRepo) GetMessage(ctx context.Context, id int64) (*model.Message, error) {
	var message model.Message
	err := r.data.db.WithContext(ctx).Where("id = ?", id).First(&message).Error
	if err != nil {
		return nil, err
	}
	return &message, nil
}

// UpdateMessage 更新消息
func (r *conversationRepo) UpdateMessage(ctx context.Context, message *model.Message) (*model.Message, error) {
	if err := r.data.db.WithContext(ctx).Save(message).Error; err != nil {
		return nil, err
	}
	return message, nil
}

// DeleteMessage 删除消息
func (r *conversationRepo) DeleteMessage(ctx context.Context, messageIDs []int64, hardDelete bool) error {
	if len(messageIDs) == 0 {
		return nil
	}

	if hardDelete {
		return r.data.db.WithContext(ctx).Unscoped().Delete(&model.Message{}, messageIDs).Error
	} else {
		return r.data.db.WithContext(ctx).Model(&model.Message{}).Where("id IN ?", messageIDs).Update("status", 5).Error
	}
}

// ListMessages 获取消息列表
func (r *conversationRepo) ListMessages(ctx context.Context, conversationID int64, page, pageSize int32, filters biz.MessageFilter) ([]*model.Message, int64, error) {
	var messages []*model.Message
	var total int64

	query := r.data.db.WithContext(ctx).Model(&model.Message{}).Where("conversation_id = ?", conversationID)

	// 应用过滤条件
	if filters.Role != "" {
		query = query.Where("role = ?", filters.Role)
	}

	if filters.Status > 0 {
		query = query.Where("status = ?", filters.Status)
	}

	// 预加载关联数据
	if filters.IncludeToolCalls {
		query = query.Preload("ToolCalls")
	}

	if filters.IncludeAttachments {
		query = query.Preload("Attachments")
	}

	// 获取总数
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 分页和排序
	offset := (page - 1) * pageSize
	if err := query.Order("created_at ASC").Offset(int(offset)).Limit(int(pageSize)).Find(&messages).Error; err != nil {
		return nil, 0, err
	}

	return messages, total, nil
}

// CreateToolCall 创建工具调用
func (r *conversationRepo) CreateToolCall(ctx context.Context, toolCall *model.ToolCall) (*model.ToolCall, error) {
	if err := r.data.db.WithContext(ctx).Create(toolCall).Error; err != nil {
		return nil, err
	}
	return toolCall, nil
}

// UpdateToolCall 更新工具调用
func (r *conversationRepo) UpdateToolCall(ctx context.Context, toolCall *model.ToolCall) (*model.ToolCall, error) {
	if err := r.data.db.WithContext(ctx).Save(toolCall).Error; err != nil {
		return nil, err
	}
	return toolCall, nil
}

// GetConversationMemory 获取对话记忆
func (r *conversationRepo) GetConversationMemory(ctx context.Context, conversationID int64) (*model.ConversationMemory, error) {
	var memory model.ConversationMemory
	err := r.data.db.WithContext(ctx).Where("conversation_id = ?", conversationID).First(&memory).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil // 没有记忆记录
		}
		return nil, err
	}
	return &memory, nil
}

// SetConversationMemory 设置对话记忆
func (r *conversationRepo) SetConversationMemory(ctx context.Context, memory *model.ConversationMemory) error {
	// 使用UPSERT操作
	return r.data.db.WithContext(ctx).Save(memory).Error
}

// UpdateConversationStats 更新对话统计信息
func (r *conversationRepo) UpdateConversationStats(ctx context.Context, conversationID int64, inputTokens, outputTokens int64, duration time.Duration) error {
	updates := map[string]interface{}{
		"total_input_tokens":  gorm.Expr("total_input_tokens + ?", inputTokens),
		"total_output_tokens": gorm.Expr("total_output_tokens + ?", outputTokens),
		"total_duration":      gorm.Expr("total_duration + ?", duration.Milliseconds()),
		"updated_at":          time.Now(),
	}

	return r.data.db.WithContext(ctx).Model(&model.Conversation{}).Where("id = ?", conversationID).Updates(updates).Error
}

// GetConversationStats 获取对话统计信息
func (r *conversationRepo) GetConversationStats(ctx context.Context, conversationID int64) (*biz.ConversationStats, error) {
	var conversation model.Conversation
	if err := r.data.db.WithContext(ctx).Where("id = ?", conversationID).First(&conversation).Error; err != nil {
		return nil, err
	}

	// 计算各角色消息数量
	var userMessageCount, assistantMessageCount int64
	r.data.db.WithContext(ctx).Model(&model.Message{}).
		Where("conversation_id = ? AND role = ?", conversationID, "user").
		Count(&userMessageCount)
	r.data.db.WithContext(ctx).Model(&model.Message{}).
		Where("conversation_id = ? AND role = ?", conversationID, "assistant").
		Count(&assistantMessageCount)

	// 计算平均响应时间
	var avgResponseTime float64
	r.data.db.WithContext(ctx).Model(&model.Message{}).
		Where("conversation_id = ? AND role = ? AND response_time > 0", conversationID, "assistant").
		Select("AVG(response_time)").Scan(&avgResponseTime)

	// 获取最后消息时间
	var lastMessage model.Message
	var lastMessageAt *time.Time
	if err := r.data.db.WithContext(ctx).Where("conversation_id = ?", conversationID).
		Order("created_at DESC").First(&lastMessage).Error; err == nil {
		lastMessageAt = &lastMessage.CreatedAt
	}

	stats := &biz.ConversationStats{
		MessageCount:          conversation.MessageCount,
		UserMessageCount:      userMessageCount,
		AssistantMessageCount: assistantMessageCount,
		TotalInputTokens:      conversation.TotalInputTokens,
		TotalOutputTokens:     conversation.TotalOutputTokens,
		ToolCallCount:         conversation.ToolCallCount,
		TotalDuration:         time.Duration(conversation.TotalDuration) * time.Millisecond,
		AverageResponseTime:   avgResponseTime,
		LastMessageAt:         lastMessageAt,
	}

	return stats, nil
}
