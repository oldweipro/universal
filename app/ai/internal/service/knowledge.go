package service

import (
	"context"
	"io"

	pb "universal/api/ai/v1"
	"universal/app/ai/internal/biz"
	"universal/app/ai/internal/data/model"

	"github.com/go-kratos/kratos/v2/log"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type KnowledgeService struct {
	pb.UnimplementedKnowledgeServer

	uc     *biz.KnowledgeUsecase
	logger *log.Helper
}

func NewKnowledgeService(uc *biz.KnowledgeUsecase, logger log.Logger) *KnowledgeService {
	return &KnowledgeService{
		uc:     uc,
		logger: log.NewHelper(logger),
	}
}

func (s *KnowledgeService) CreateKnowledgeBase(ctx context.Context, req *pb.CreateKnowledgeBaseRequest) (*pb.CreateKnowledgeBaseReply, error) {
	// 构建配置
	config := model.KnowledgeBaseConfig{
		EmbeddingDimension:   int(req.Config.EmbeddingDimension),
		ChunkingStrategy:     s.convertChunkingStrategy(req.Config.ChunkingStrategy),
		SimilarityThreshold:  req.Config.SimilarityThreshold,
		MaxChunksPerQuery:    int(req.Config.MaxChunksPerQuery),
		EnableMetadataFilter: req.Config.EnableMetadataFilter,
		StopWords:            req.Config.StopWords,
		TextSplitter:         req.Config.TextSplitter,
	}

	if req.Config.AdvancedOptions != nil {
		config.AdvancedOptions = make(map[string]interface{})
		//for k, v := range req.Config.AdvancedOptions {
		//config.AdvancedOptions[k] = v.String()
		//}
	}

	kb, err := s.uc.CreateKnowledgeBase(
		ctx,
		req.UserId,
		req.Name,
		req.Description,
		req.EmbeddingModel,
		req.ChunkSize,
		req.ChunkOverlap,
		config,
		req.Tags,
		req.Language,
	)
	if err != nil {
		return nil, err
	}

	return &pb.CreateKnowledgeBaseReply{
		KnowledgeBase: s.convertKnowledgeBaseToProto(kb),
	}, nil
}
func (s *KnowledgeService) UpdateKnowledgeBase(ctx context.Context, req *pb.UpdateKnowledgeBaseRequest) (*pb.UpdateKnowledgeBaseReply, error) {
	return &pb.UpdateKnowledgeBaseReply{}, nil
}
func (s *KnowledgeService) DeleteKnowledgeBase(ctx context.Context, req *pb.DeleteKnowledgeBaseRequest) (*pb.DeleteKnowledgeBaseReply, error) {
	return &pb.DeleteKnowledgeBaseReply{}, nil
}
func (s *KnowledgeService) ListKnowledgeBases(ctx context.Context, req *pb.ListKnowledgeBasesRequest) (*pb.ListKnowledgeBasesReply, error) {
	return &pb.ListKnowledgeBasesReply{}, nil
}
func (s *KnowledgeService) GetKnowledgeBase(ctx context.Context, req *pb.GetKnowledgeBaseRequest) (*pb.GetKnowledgeBaseReply, error) {
	return &pb.GetKnowledgeBaseReply{}, nil
}
func (s *KnowledgeService) UploadDocument(ctx context.Context, req *pb.UploadDocumentRequest) (*pb.UploadDocumentReply, error) {
	return &pb.UploadDocumentReply{}, nil
}
func (s *KnowledgeService) BatchUploadDocuments(conn pb.Knowledge_BatchUploadDocumentsServer) error {
	for {
		_, err := conn.Recv()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			return err
		}

		err = conn.Send(&pb.BatchUploadDocumentsReply{})
		if err != nil {
			return err
		}
	}
}
func (s *KnowledgeService) UpdateDocument(ctx context.Context, req *pb.UpdateDocumentRequest) (*pb.UpdateDocumentReply, error) {
	return &pb.UpdateDocumentReply{}, nil
}
func (s *KnowledgeService) DeleteDocument(ctx context.Context, req *pb.DeleteDocumentRequest) (*pb.DeleteDocumentReply, error) {
	return &pb.DeleteDocumentReply{}, nil
}
func (s *KnowledgeService) ListDocuments(ctx context.Context, req *pb.ListDocumentsRequest) (*pb.ListDocumentsReply, error) {
	return &pb.ListDocumentsReply{}, nil
}
func (s *KnowledgeService) GetDocument(ctx context.Context, req *pb.GetDocumentRequest) (*pb.GetDocumentReply, error) {
	return &pb.GetDocumentReply{}, nil
}
func (s *KnowledgeService) ProcessDocument(ctx context.Context, req *pb.ProcessDocumentRequest) (*pb.ProcessDocumentReply, error) {
	return &pb.ProcessDocumentReply{}, nil
}
func (s *KnowledgeService) SearchKnowledge(ctx context.Context, req *pb.SearchKnowledgeRequest) (*pb.SearchKnowledgeReply, error) {
	return &pb.SearchKnowledgeReply{}, nil
}
func (s *KnowledgeService) HybridSearch(ctx context.Context, req *pb.HybridSearchRequest) (*pb.HybridSearchReply, error) {
	return &pb.HybridSearchReply{}, nil
}
func (s *KnowledgeService) AdvancedSearch(ctx context.Context, req *pb.AdvancedSearchRequest) (*pb.AdvancedSearchReply, error) {
	return &pb.AdvancedSearchReply{}, nil
}
func (s *KnowledgeService) GetKnowledgeChunk(ctx context.Context, req *pb.GetKnowledgeChunkRequest) (*pb.GetKnowledgeChunkReply, error) {
	return &pb.GetKnowledgeChunkReply{}, nil
}
func (s *KnowledgeService) UpdateKnowledgeChunk(ctx context.Context, req *pb.UpdateKnowledgeChunkRequest) (*pb.UpdateKnowledgeChunkReply, error) {
	return &pb.UpdateKnowledgeChunkReply{}, nil
}
func (s *KnowledgeService) ListKnowledgeChunks(ctx context.Context, req *pb.ListKnowledgeChunksRequest) (*pb.ListKnowledgeChunksReply, error) {
	return &pb.ListKnowledgeChunksReply{}, nil
}
func (s *KnowledgeService) ReindexKnowledgeBase(ctx context.Context, req *pb.ReindexKnowledgeBaseRequest) (*pb.ReindexKnowledgeBaseReply, error) {
	return &pb.ReindexKnowledgeBaseReply{}, nil
}
func (s *KnowledgeService) GetKnowledgeBaseStats(ctx context.Context, req *pb.GetKnowledgeBaseStatsRequest) (*pb.GetKnowledgeBaseStatsReply, error) {
	return &pb.GetKnowledgeBaseStatsReply{}, nil
}
func (s *KnowledgeService) AnalyzeKnowledgeBase(ctx context.Context, req *pb.AnalyzeKnowledgeBaseRequest) (*pb.AnalyzeKnowledgeBaseReply, error) {
	return &pb.AnalyzeKnowledgeBaseReply{}, nil
}

// 转换函数
func (s *KnowledgeService) convertKnowledgeBaseToProto(kb *model.KnowledgeBase) *pb.KnowledgeBase {
	proto := &pb.KnowledgeBase{
		Id:                 kb.ID,
		UserId:             kb.UserID,
		Name:               kb.Name,
		Description:        kb.Description,
		EmbeddingModel:     kb.EmbeddingModel,
		ChunkSize:          int32(kb.ChunkSize),
		ChunkOverlap:       int32(kb.ChunkOverlap),
		Status:             pb.KnowledgeBaseStatus(kb.Status),
		Tags:               []string(kb.Tags),
		Version:            int32(kb.Version),
		Language:           kb.Language,
		SupportedFileTypes: []string(kb.SupportedFileTypes),
		MaxFileSize:        kb.MaxFileSize,
		AutoProcess:        kb.AutoProcess,
		CreatedAt:          timestamppb.New(kb.CreatedAt),
		UpdatedAt:          timestamppb.New(kb.UpdatedAt),
	}

	if kb.LastIndexedAt != nil {
		proto.LastIndexedAt = timestamppb.New(*kb.LastIndexedAt)
	}

	// 构建配置
	proto.Config = &pb.KnowledgeBaseConfig{
		//EmbeddingModel:       kb.Config.EmbeddingModel,
		EmbeddingDimension: int32(kb.Config.EmbeddingDimension),
		//ChunkSize:           int32(kb.Config.ChunkingStrategy != ""),
		ChunkOverlap:         int32(kb.ChunkOverlap),
		ChunkingStrategy:     s.convertChunkingStrategyToProto(kb.Config.ChunkingStrategy),
		SimilarityThreshold:  kb.Config.SimilarityThreshold,
		MaxChunksPerQuery:    int32(kb.Config.MaxChunksPerQuery),
		EnableMetadataFilter: kb.Config.EnableMetadataFilter,
		StopWords:            kb.Config.StopWords,
		TextSplitter:         kb.Config.TextSplitter,
	}

	// 构建统计信息
	proto.Stats = &pb.KnowledgeBaseStats{
		DocumentCount:   kb.DocumentCount,
		ChunkCount:      kb.ChunkCount,
		TotalCharacters: kb.TotalCharacters,
		TotalTokens:     kb.TotalTokens,
		StorageSizeMb:   kb.StorageSize,
		IndexSizeMb:     kb.IndexSize,
		LastUpdated:     timestamppb.New(kb.UpdatedAt),
	}

	return proto
}

func (s *KnowledgeService) convertDocumentToProto(doc *model.Document) *pb.Document {
	proto := &pb.Document{
		Id:                 doc.ID,
		KnowledgeBaseId:    doc.KnowledgeBaseID,
		Name:               doc.Name,
		Content:            doc.Content,
		FilePath:           doc.FilePath,
		MimeType:           doc.MimeType,
		FileSize:           doc.FileSize,
		ChunkCount:         int32(doc.ChunkCount),
		Status:             pb.DocumentStatus(doc.Status),
		Tags:               []string(doc.Tags),
		Language:           doc.Language,
		SourceUrl:          doc.SourceURL,
		Hash:               doc.Hash,
		Version:            int32(doc.Version),
		ProcessingProgress: doc.ProcessingProgress,
		ProcessingError:    doc.ProcessingError,
		CreatedAt:          timestamppb.New(doc.CreatedAt),
		UpdatedAt:          timestamppb.New(doc.UpdatedAt),
	}

	if doc.LastProcessedAt != nil {
		proto.LastProcessedAt = timestamppb.New(*doc.LastProcessedAt)
	}

	// 构建元数据
	proto.Metadata = &pb.DocumentMetadata{
		Title:     doc.Metadata.Title,
		Author:    doc.Metadata.Author,
		Subject:   doc.Metadata.Subject,
		Keywords:  doc.Metadata.Keywords,
		Category:  doc.Metadata.Category,
		PageCount: int32(doc.Metadata.PageCount),
		Encoding:  doc.Metadata.Encoding,
	}

	if doc.Metadata.CreatedDate != nil {
		proto.Metadata.CreatedDate = timestamppb.New(*doc.Metadata.CreatedDate)
	}

	if doc.Metadata.ModifiedDate != nil {
		proto.Metadata.ModifiedDate = timestamppb.New(*doc.Metadata.ModifiedDate)
	}

	if doc.Metadata.CustomFields != nil {
		proto.Metadata.CustomFields = doc.Metadata.CustomFields
	}

	return proto
}

func (s *KnowledgeService) convertKnowledgeChunkToProto(chunk *model.KnowledgeChunk) *pb.KnowledgeChunk {
	proto := &pb.KnowledgeChunk{
		Id:              chunk.ID,
		DocumentId:      chunk.DocumentID,
		KnowledgeBaseId: chunk.KnowledgeBaseID,
		Content:         chunk.Content,
		Score:           chunk.Score,
		ChunkIndex:      int32(chunk.ChunkIndex),
		StartPosition:   int32(chunk.StartPosition),
		EndPosition:     int32(chunk.EndPosition),
		Language:        chunk.Language,
		Keywords:        []string(chunk.Keywords),
		ChunkType:       pb.ChunkType(chunk.ChunkType),
		CharacterCount:  int32(chunk.CharacterCount),
		TokenCount:      int32(chunk.TokenCount),
		CreatedAt:       timestamppb.New(chunk.CreatedAt),
	}

	// 构建元数据
	if chunk.Metadata.Custom != nil {
		proto.Metadata = chunk.Metadata.Custom
	}

	return proto
}

func (s *KnowledgeService) convertChunkingStrategy(strategy pb.ChunkingStrategy) string {
	switch strategy {
	case pb.ChunkingStrategy_CHUNKING_STRATEGY_FIXED_SIZE:
		return "fixed_size"
	case pb.ChunkingStrategy_CHUNKING_STRATEGY_RECURSIVE:
		return "recursive"
	case pb.ChunkingStrategy_CHUNKING_STRATEGY_SEMANTIC:
		return "semantic"
	case pb.ChunkingStrategy_CHUNKING_STRATEGY_PARAGRAPH:
		return "paragraph"
	case pb.ChunkingStrategy_CHUNKING_STRATEGY_SENTENCE:
		return "sentence"
	default:
		return "fixed_size"
	}
}

func (s *KnowledgeService) convertChunkingStrategyToProto(strategy string) pb.ChunkingStrategy {
	switch strategy {
	case "fixed_size":
		return pb.ChunkingStrategy_CHUNKING_STRATEGY_FIXED_SIZE
	case "recursive":
		return pb.ChunkingStrategy_CHUNKING_STRATEGY_RECURSIVE
	case "semantic":
		return pb.ChunkingStrategy_CHUNKING_STRATEGY_SEMANTIC
	case "paragraph":
		return pb.ChunkingStrategy_CHUNKING_STRATEGY_PARAGRAPH
	case "sentence":
		return pb.ChunkingStrategy_CHUNKING_STRATEGY_SENTENCE
	default:
		return pb.ChunkingStrategy_CHUNKING_STRATEGY_FIXED_SIZE
	}
}
