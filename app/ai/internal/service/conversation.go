package service

import (
	"context"

	pb "universal/api/ai/v1"
	"universal/app/ai/internal/biz"
	"universal/app/ai/internal/data/model"

	"github.com/go-kratos/kratos/v2/log"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type ConversationService struct {
	pb.UnimplementedConversationServer

	uc     *biz.ConversationUsecase
	logger *log.Helper
}

func NewConversationService(uc *biz.ConversationUsecase, logger log.Logger) *ConversationService {
	return &ConversationService{
		uc:     uc,
		logger: log.NewHelper(logger),
	}
}

func (s *ConversationService) CreateConversation(ctx context.Context, req *pb.CreateConversationRequest) (*pb.CreateConversationReply, error) {
	// 构建配置
	config := model.ConversationConfig{}
	if len(req.Config) > 0 {
		for k, v := range req.Config {
			if config.CustomParams == nil {
				config.CustomParams = make(map[string]interface{})
			}
			config.CustomParams[k] = v
		}
	}

	conversation, err := s.uc.CreateConversation(
		ctx,
		req.UserId,
		req.Title,
		req.ModelName,
		req.SystemPrompt,
		config,
		req.Description,
		req.Tags,
		req.Priority,
	)
	if err != nil {
		return nil, err
	}

	return &pb.CreateConversationReply{
		Conversation: s.convertConversationToProto(conversation),
	}, nil
}
func (s *ConversationService) GetConversation(ctx context.Context, req *pb.GetConversationRequest) (*pb.GetConversationReply, error) {
	conversation, err := s.uc.GetConversation(ctx, req.Id)
	if err != nil {
		return nil, err
	}

	return &pb.GetConversationReply{
		Conversation: s.convertConversationToProto(conversation),
	}, nil
}
func (s *ConversationService) UpdateConversation(ctx context.Context, req *pb.UpdateConversationRequest) (*pb.UpdateConversationReply, error) {
	// 构建配置
	config := model.ConversationConfig{}
	if len(req.Config) > 0 {
		config.CustomParams = make(map[string]interface{})
		for k, v := range req.Config {
			config.CustomParams[k] = v
		}
	}

	conversation, err := s.uc.UpdateConversation(
		ctx,
		req.Id,
		req.Title,
		req.SystemPrompt,
		req.Description,
		config,
		req.Tags,
		req.Priority,
	)
	if err != nil {
		return nil, err
	}

	return &pb.UpdateConversationReply{
		Conversation: s.convertConversationToProto(conversation),
	}, nil
}

func (s *ConversationService) DeleteConversation(ctx context.Context, req *pb.DeleteConversationRequest) (*pb.DeleteConversationReply, error) {
	err := s.uc.DeleteConversation(ctx, req.Id, req.HardDelete)
	if err != nil {
		return nil, err
	}

	return &pb.DeleteConversationReply{}, nil
}

func (s *ConversationService) ListConversations(ctx context.Context, req *pb.ListConversationsRequest) (*pb.ListConversationsReply, error) {
	conversations, total, err := s.uc.ListConversations(
		ctx,
		req.UserId,
		req.Page,
		req.PageSize,
		req.Keyword,
		int32(req.Status),
		req.Tags,
		req.SortBy,
		req.SortDesc,
	)
	if err != nil {
		return nil, err
	}

	// 转换为Proto消息
	protoConversations := make([]*pb.ConversationInfo, len(conversations))
	for i, conv := range conversations {
		protoConversations[i] = s.convertConversationToProto(conv)
	}

	return &pb.ListConversationsReply{
		Conversations: protoConversations,
		Total:         total,
		Page:          req.Page,
		PageSize:      req.PageSize,
	}, nil
}
func (s *ConversationService) ArchiveConversation(ctx context.Context, req *pb.ArchiveConversationRequest) (*pb.ArchiveConversationReply, error) {
	return &pb.ArchiveConversationReply{}, nil
}
func (s *ConversationService) RestoreConversation(ctx context.Context, req *pb.RestoreConversationRequest) (*pb.RestoreConversationReply, error) {
	return &pb.RestoreConversationReply{}, nil
}
func (s *ConversationService) SendMessage(ctx context.Context, req *pb.SendMessageRequest) (*pb.SendMessageReply, error) {
	// 转换附件信息
	var attachments []biz.MessageAttachmentInfo
	for _, att := range req.Attachments {
		attachments = append(attachments, biz.MessageAttachmentInfo{
			Name:     att.Name,
			URL:      att.Url,
			MimeType: att.MimeType,
			Size:     att.Size,
			Type:     int(att.Type),
			Metadata: att.Metadata,
		})
	}

	var parentID *int64
	if req.ParentMessageId != "" {
		// TODO: 解析parentMessageId字符串为int64
	}

	userMsg, assistantMsg, err := s.uc.SendMessage(
		ctx,
		req.ConversationId,
		req.Content,
		attachments,
		req.EnableTools,
		req.AllowedTools,
		req.Options,
		parentID,
	)
	if err != nil {
		return nil, err
	}

	reply := &pb.SendMessageReply{
		UserMessage: s.convertMessageToProto(userMsg),
	}

	if assistantMsg != nil {
		reply.AssistantMessage = s.convertMessageToProto(assistantMsg)
	}

	return reply, nil
}
func (s *ConversationService) SendStreamMessage(req *pb.SendMessageRequest, conn pb.Conversation_SendStreamMessageServer) error {
	for {
		err := conn.Send(&pb.SendMessageStreamReply{})
		if err != nil {
			return err
		}
	}
}
func (s *ConversationService) GetMessages(ctx context.Context, req *pb.GetMessagesRequest) (*pb.GetMessagesReply, error) {
	return &pb.GetMessagesReply{}, nil
}
func (s *ConversationService) DeleteMessage(ctx context.Context, req *pb.DeleteMessageRequest) (*pb.DeleteMessageReply, error) {
	return &pb.DeleteMessageReply{}, nil
}
func (s *ConversationService) RegenerateMessage(ctx context.Context, req *pb.RegenerateMessageRequest) (*pb.RegenerateMessageReply, error) {
	return &pb.RegenerateMessageReply{}, nil
}
func (s *ConversationService) GetConversationContext(ctx context.Context, req *pb.GetConversationContextRequest) (*pb.GetConversationContextReply, error) {
	return &pb.GetConversationContextReply{}, nil
}
func (s *ConversationService) UpdateConversationContext(ctx context.Context, req *pb.UpdateConversationContextRequest) (*pb.UpdateConversationContextReply, error) {
	return &pb.UpdateConversationContextReply{}, nil
}
func (s *ConversationService) SummarizeConversation(ctx context.Context, req *pb.SummarizeConversationRequest) (*pb.SummarizeConversationReply, error) {
	return &pb.SummarizeConversationReply{}, nil
}
func (s *ConversationService) ClearConversationHistory(ctx context.Context, req *pb.ClearConversationHistoryRequest) (*pb.ClearConversationHistoryReply, error) {
	return &pb.ClearConversationHistoryReply{}, nil
}
func (s *ConversationService) SetConversationMemory(ctx context.Context, req *pb.SetConversationMemoryRequest) (*pb.SetConversationMemoryReply, error) {
	return &pb.SetConversationMemoryReply{}, nil
}
func (s *ConversationService) GetConversationMemory(ctx context.Context, req *pb.GetConversationMemoryRequest) (*pb.GetConversationMemoryReply, error) {
	return &pb.GetConversationMemoryReply{}, nil
}
func (s *ConversationService) GetConversationStats(ctx context.Context, req *pb.GetConversationStatsRequest) (*pb.GetConversationStatsReply, error) {
	return &pb.GetConversationStatsReply{}, nil
}
func (s *ConversationService) ExportConversation(ctx context.Context, req *pb.ExportConversationRequest) (*pb.ExportConversationReply, error) {
	return &pb.ExportConversationReply{}, nil
}
func (s *ConversationService) ImportConversation(ctx context.Context, req *pb.ImportConversationRequest) (*pb.ImportConversationReply, error) {
	return &pb.ImportConversationReply{}, nil
}

// convertConversationToProto 将模型转换为Proto消息
func (s *ConversationService) convertConversationToProto(conv *model.Conversation) *pb.ConversationInfo {
	proto := &pb.ConversationInfo{
		Id:           conv.ID,
		UserId:       conv.UserID,
		Title:        conv.Title,
		ModelName:    conv.ModelName,
		SystemPrompt: conv.SystemPrompt,
		Config:       make(map[string]string),
		Status:       pb.ConversationStatus(conv.Status),
		CreatedAt:    timestamppb.New(conv.CreatedAt),
		UpdatedAt:    timestamppb.New(conv.UpdatedAt),
		LastActiveAt: timestamppb.New(conv.LastActiveAt),
		Description:  conv.Description,
		Tags:         []string(conv.Tags),
		Priority:     int32(conv.Priority),
	}

	// 转换配置
	if conv.Config.CustomParams != nil {
		for k, v := range conv.Config.CustomParams {
			if str, ok := v.(string); ok {
				proto.Config[k] = str
			}
		}
	}

	// 构建统计信息
	proto.Stats = &pb.ConversationStats{
		MessageCount:      conv.MessageCount,
		TotalInputTokens:  conv.TotalInputTokens,
		TotalOutputTokens: conv.TotalOutputTokens,
		ToolCallCount:     conv.ToolCallCount,
		TotalDuration:     nil, // TODO: 转换时间
	}

	return proto
}

// convertMessageToProto 将消息模型转换为Proto消息
func (s *ConversationService) convertMessageToProto(msg *model.Message) *pb.Message {
	proto := &pb.Message{
		Id:             msg.ID,
		ConversationId: msg.ConversationID,
		Role:           pb.MessageRole(s.roleStringToEnum(msg.Role)),
		Content:        msg.Content,
		CreatedAt:      timestamppb.New(msg.CreatedAt),
		Status:         pb.MessageStatus(msg.Status),
		IsEdited:       msg.IsEdited,
		EditReason:     msg.EditReason,
	}

	if msg.EditedAt != nil {
		proto.EditedAt = timestamppb.New(*msg.EditedAt)
	}

	// 构建指标信息
	proto.Metrics = &pb.MessageMetrics{
		InputTokens:  int32(msg.InputTokens),
		OutputTokens: int32(msg.OutputTokens),
		ResponseTime: msg.ResponseTime,
		Cost:         msg.Cost,
		ModelUsed:    msg.ModelUsed,
	}

	return proto
}

// roleStringToEnum 将角色字符串转换为枚举
func (s *ConversationService) roleStringToEnum(role string) int32 {
	switch role {
	case "user":
		return 1
	case "assistant":
		return 2
	case "system":
		return 3
	case "tool":
		return 4
	default:
		return 0
	}
}
