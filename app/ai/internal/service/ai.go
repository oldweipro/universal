package service

import (
	"context"
	"universal/app/ai/internal/biz"

	pb "universal/api/ai/v1"
)

type AiService struct {
	pb.UnimplementedAiServer
	uc *biz.AiUsecase
}

func NewAiService(uc *biz.AiUsecase) *AiService {
	return &AiService{uc: uc}
}

func (s *AiService) CreateConversation(ctx context.Context, req *pb.CreateConversationRequest) (*pb.CreateConversationReply, error) {
	return &pb.CreateConversationReply{}, nil
}
func (s *AiService) GetConversation(ctx context.Context, req *pb.GetConversationRequest) (*pb.GetConversationReply, error) {
	return &pb.GetConversationReply{}, nil
}
func (s *AiService) UpdateConversation(ctx context.Context, req *pb.UpdateConversationRequest) (*pb.UpdateConversationReply, error) {
	return &pb.UpdateConversationReply{}, nil
}
func (s *AiService) DeleteConversation(ctx context.Context, req *pb.DeleteConversationRequest) (*pb.DeleteConversationReply, error) {
	return &pb.DeleteConversationReply{}, nil
}
func (s *AiService) ListConversations(ctx context.Context, req *pb.ListConversationsRequest) (*pb.ListConversationsReply, error) {
	return &pb.ListConversationsReply{}, nil
}
func (s *AiService) SendMessage(ctx context.Context, req *pb.SendMessageRequest) (*pb.SendMessageReply, error) {
	return &pb.SendMessageReply{}, nil
}
func (s *AiService) GetMessages(ctx context.Context, req *pb.GetMessagesRequest) (*pb.GetMessagesReply, error) {
	return &pb.GetMessagesReply{}, nil
}
func (s *AiService) ListModels(ctx context.Context, req *pb.ListModelsRequest) (*pb.ListModelsReply, error) {
	return &pb.ListModelsReply{}, nil
}
func (s *AiService) GetModelConfig(ctx context.Context, req *pb.GetModelConfigRequest) (*pb.GetModelConfigReply, error) {
	return &pb.GetModelConfigReply{}, nil
}
func (s *AiService) UpdateModelConfig(ctx context.Context, req *pb.UpdateModelConfigRequest) (*pb.UpdateModelConfigReply, error) {
	return &pb.UpdateModelConfigReply{}, nil
}
func (s *AiService) ListTools(ctx context.Context, req *pb.ListToolsRequest) (*pb.ListToolsReply, error) {
	return &pb.ListToolsReply{}, nil
}
func (s *AiService) CallTool(ctx context.Context, req *pb.CallToolRequest) (*pb.CallToolResponse, error) {
	return &pb.CallToolResponse{}, nil
}
func (s *AiService) GetToolSchema(ctx context.Context, req *pb.GetToolSchemaRequest) (*pb.GetToolSchemaReply, error) {
	return &pb.GetToolSchemaReply{}, nil
}
func (s *AiService) ListResources(ctx context.Context, req *pb.ListResourcesRequest) (*pb.ListResourcesReply, error) {
	return &pb.ListResourcesReply{}, nil
}
func (s *AiService) GetResource(ctx context.Context, req *pb.GetResourceRequest) (*pb.GetResourceReply, error) {
	return &pb.GetResourceReply{}, nil
}
func (s *AiService) CreateKnowledgeBase(ctx context.Context, req *pb.CreateKnowledgeBaseRequest) (*pb.CreateKnowledgeBaseReply, error) {
	return &pb.CreateKnowledgeBaseReply{}, nil
}
func (s *AiService) UpdateKnowledgeBase(ctx context.Context, req *pb.UpdateKnowledgeBaseRequest) (*pb.UpdateKnowledgeBaseReply, error) {
	return &pb.UpdateKnowledgeBaseReply{}, nil
}
func (s *AiService) DeleteKnowledgeBase(ctx context.Context, req *pb.DeleteKnowledgeBaseRequest) (*pb.DeleteKnowledgeBaseReply, error) {
	return &pb.DeleteKnowledgeBaseReply{}, nil
}
func (s *AiService) ListKnowledgeBases(ctx context.Context, req *pb.ListKnowledgeBasesRequest) (*pb.ListKnowledgeBasesReply, error) {
	return &pb.ListKnowledgeBasesReply{}, nil
}
func (s *AiService) UploadDocument(ctx context.Context, req *pb.UploadDocumentRequest) (*pb.UploadDocumentReply, error) {
	return &pb.UploadDocumentReply{}, nil
}
func (s *AiService) SearchKnowledge(ctx context.Context, req *pb.SearchKnowledgeRequest) (*pb.SearchKnowledgeReply, error) {
	return &pb.SearchKnowledgeReply{}, nil
}
