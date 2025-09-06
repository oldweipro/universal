package service

import (
	"context"
	aiv1 "universal/api/ai/v1"
	gatewayv1 "universal/api/gateway/v1"
	"universal/app/gateway/internal/biz"
	"universal/app/gateway/internal/data"

	"github.com/go-kratos/kratos/v2/log"
)

// ConversationService 对话服务代理
type ConversationService struct {
	gatewayv1.UnimplementedConversationServer

	data *data.Data
	uc   *biz.UserUsecase
	log  *log.Helper
}

// NewConversationService 创建对话服务
func NewConversationService(data *data.Data, uc *biz.UserUsecase, logger log.Logger) *ConversationService {
	return &ConversationService{
		data: data,
		uc:   uc,
		log:  log.NewHelper(logger),
	}
}

// CreateConversation 创建新对话
func (s *ConversationService) CreateConversation(ctx context.Context, req *aiv1.CreateConversationRequest) (*aiv1.CreateConversationReply, error) {
	s.log.WithContext(ctx).Infof("CreateConversation called for user: %d", req.UserId)

	// 这里可以添加用户权限验证等逻辑
	// 然后转发到AI服务的对话模块
	return s.data.ConversationClient().CreateConversation(ctx, req)
}

// GetConversation 获取对话信息
func (s *ConversationService) GetConversation(ctx context.Context, req *aiv1.GetConversationRequest) (*aiv1.GetConversationReply, error) {
	s.log.WithContext(ctx).Infof("GetConversation called for conversation: %d", req.Id)
	return s.data.ConversationClient().GetConversation(ctx, req)
}

// UpdateConversation 更新对话信息
func (s *ConversationService) UpdateConversation(ctx context.Context, req *aiv1.UpdateConversationRequest) (*aiv1.UpdateConversationReply, error) {
	s.log.WithContext(ctx).Infof("UpdateConversation called for conversation: %d", req.Id)
	return s.data.ConversationClient().UpdateConversation(ctx, req)
}

// DeleteConversation 删除对话
func (s *ConversationService) DeleteConversation(ctx context.Context, req *aiv1.DeleteConversationRequest) (*aiv1.DeleteConversationReply, error) {
	s.log.WithContext(ctx).Infof("DeleteConversation called for conversation: %d", req.Id)
	return s.data.ConversationClient().DeleteConversation(ctx, req)
}

// ListConversations 列出对话
func (s *ConversationService) ListConversations(ctx context.Context, req *aiv1.ListConversationsRequest) (*aiv1.ListConversationsReply, error) {
	s.log.WithContext(ctx).Infof("ListConversations called for user: %d", req.UserId)
	return s.data.ConversationClient().ListConversations(ctx, req)
}

// ArchiveConversation 归档对话
func (s *ConversationService) ArchiveConversation(ctx context.Context, req *aiv1.ArchiveConversationRequest) (*aiv1.ArchiveConversationReply, error) {
	s.log.WithContext(ctx).Infof("ArchiveConversation called for conversation: %d", req.Id)
	return &aiv1.ArchiveConversationReply{}, nil
}

// RestoreConversation 恢复对话
func (s *ConversationService) RestoreConversation(ctx context.Context, req *aiv1.RestoreConversationRequest) (*aiv1.RestoreConversationReply, error) {
	s.log.WithContext(ctx).Infof("RestoreConversation called for conversation: %d", req.Id)
	return &aiv1.RestoreConversationReply{}, nil
}

// SendMessage 发送消息
func (s *ConversationService) SendMessage(ctx context.Context, req *aiv1.SendMessageRequest) (*aiv1.SendMessageReply, error) {
	s.log.WithContext(ctx).Infof("SendMessage called for conversation: %d", req.ConversationId)

	// 转发到AI服务的对话客户端
	return s.data.ConversationClient().SendMessage(ctx, req)
}

// SendStreamMessage 流式发送消息
func (s *ConversationService) SendStreamMessage(req *aiv1.SendMessageRequest, stream aiv1.Conversation_SendStreamMessageServer) error {
	s.log.Infof("SendStreamMessage called for conversation: %d", req.ConversationId)

	// 创建AI服务的流式客户端
	aiStream, err := s.data.ConversationClient().SendStreamMessage(context.Background(), req)
	if err != nil {
		return err
	}

	// 转发流式响应
	for {
		resp, err := aiStream.Recv()
		if err != nil {
			return err
		}

		if err := stream.Send(resp); err != nil {
			return err
		}

		if resp.IsComplete {
			break
		}
	}

	return nil
}

// GetMessages 获取消息列表
func (s *ConversationService) GetMessages(ctx context.Context, req *aiv1.GetMessagesRequest) (*aiv1.GetMessagesReply, error) {
	s.log.WithContext(ctx).Infof("GetMessages called for conversation: %d", req.ConversationId)
	return &aiv1.GetMessagesReply{}, nil
}

// DeleteMessage 删除消息
func (s *ConversationService) DeleteMessage(ctx context.Context, req *aiv1.DeleteMessageRequest) (*aiv1.DeleteMessageReply, error) {
	s.log.WithContext(ctx).Infof("DeleteMessage called for messages: %v", req.MessageIds)
	return &aiv1.DeleteMessageReply{}, nil
}

// RegenerateMessage 重新生成消息
func (s *ConversationService) RegenerateMessage(ctx context.Context, req *aiv1.RegenerateMessageRequest) (*aiv1.RegenerateMessageReply, error) {
	s.log.WithContext(ctx).Infof("RegenerateMessage called for message: %d", req.MessageId)
	return &aiv1.RegenerateMessageReply{}, nil
}

// GetConversationContext 获取对话上下文
func (s *ConversationService) GetConversationContext(ctx context.Context, req *aiv1.GetConversationContextRequest) (*aiv1.GetConversationContextReply, error) {
	s.log.WithContext(ctx).Infof("GetConversationContext called for conversation: %d", req.ConversationId)
	return &aiv1.GetConversationContextReply{}, nil
}

// UpdateConversationContext 更新对话上下文
func (s *ConversationService) UpdateConversationContext(ctx context.Context, req *aiv1.UpdateConversationContextRequest) (*aiv1.UpdateConversationContextReply, error) {
	s.log.WithContext(ctx).Infof("UpdateConversationContext called for conversation: %d", req.ConversationId)
	return &aiv1.UpdateConversationContextReply{}, nil
}

// SummarizeConversation 总结对话
func (s *ConversationService) SummarizeConversation(ctx context.Context, req *aiv1.SummarizeConversationRequest) (*aiv1.SummarizeConversationReply, error) {
	s.log.WithContext(ctx).Infof("SummarizeConversation called for conversation: %d", req.ConversationId)
	return &aiv1.SummarizeConversationReply{}, nil
}

// ClearConversationHistory 清空对话历史
func (s *ConversationService) ClearConversationHistory(ctx context.Context, req *aiv1.ClearConversationHistoryRequest) (*aiv1.ClearConversationHistoryReply, error) {
	s.log.WithContext(ctx).Infof("ClearConversationHistory called for conversation: %d", req.ConversationId)
	return &aiv1.ClearConversationHistoryReply{}, nil
}

// SetConversationMemory 设置对话记忆
func (s *ConversationService) SetConversationMemory(ctx context.Context, req *aiv1.SetConversationMemoryRequest) (*aiv1.SetConversationMemoryReply, error) {
	s.log.WithContext(ctx).Infof("SetConversationMemory called for conversation: %d", req.ConversationId)
	return &aiv1.SetConversationMemoryReply{}, nil
}

// GetConversationMemory 获取对话记忆
func (s *ConversationService) GetConversationMemory(ctx context.Context, req *aiv1.GetConversationMemoryRequest) (*aiv1.GetConversationMemoryReply, error) {
	s.log.WithContext(ctx).Infof("GetConversationMemory called for conversation: %d", req.ConversationId)
	return &aiv1.GetConversationMemoryReply{}, nil
}

// GetConversationStats 获取对话统计
func (s *ConversationService) GetConversationStats(ctx context.Context, req *aiv1.GetConversationStatsRequest) (*aiv1.GetConversationStatsReply, error) {
	s.log.WithContext(ctx).Infof("GetConversationStats called for conversation: %d", req.ConversationId)
	return &aiv1.GetConversationStatsReply{}, nil
}

// ExportConversation 导出对话
func (s *ConversationService) ExportConversation(ctx context.Context, req *aiv1.ExportConversationRequest) (*aiv1.ExportConversationReply, error) {
	s.log.WithContext(ctx).Infof("ExportConversation called for conversation: %d", req.ConversationId)
	return &aiv1.ExportConversationReply{}, nil
}

// ImportConversation 导入对话
func (s *ConversationService) ImportConversation(ctx context.Context, req *aiv1.ImportConversationRequest) (*aiv1.ImportConversationReply, error) {
	s.log.WithContext(ctx).Infof("ImportConversation called for user: %d", req.UserId)
	return &aiv1.ImportConversationReply{}, nil
}
