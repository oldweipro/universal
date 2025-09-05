package biz

import (
	"context"
	"time"

	"universal/app/ai/internal/data/model"

	"github.com/go-kratos/kratos/v2/log"
)

// KnowledgeUsecase 知识库业务逻辑
type KnowledgeUsecase struct {
	repo   KnowledgeRepo
	logger *log.Helper
}

// KnowledgeRepo 知识库仓库接口
type KnowledgeRepo interface {
	// 知识库管理
	CreateKnowledgeBase(ctx context.Context, kb *model.KnowledgeBase) (*model.KnowledgeBase, error)
	GetKnowledgeBase(ctx context.Context, id int64) (*model.KnowledgeBase, error)
	UpdateKnowledgeBase(ctx context.Context, kb *model.KnowledgeBase) (*model.KnowledgeBase, error)
	DeleteKnowledgeBase(ctx context.Context, id int64, hardDelete bool) error
	ListKnowledgeBases(ctx context.Context, userID int64, page, pageSize int32, filters KnowledgeBaseFilter) ([]*model.KnowledgeBase, int64, error)

	// 文档管理
	CreateDocument(ctx context.Context, doc *model.Document) (*model.Document, error)
	GetDocument(ctx context.Context, id int64) (*model.Document, error)
	UpdateDocument(ctx context.Context, doc *model.Document) (*model.Document, error)
	DeleteDocument(ctx context.Context, documentIDs []int64, hardDelete bool) error
	ListDocuments(ctx context.Context, kbID int64, page, pageSize int32, filters DocumentFilter) ([]*model.Document, int64, error)

	// 知识块管理
	CreateKnowledgeChunks(ctx context.Context, chunks []*model.KnowledgeChunk) error
	GetKnowledgeChunk(ctx context.Context, id int64) (*model.KnowledgeChunk, error)
	UpdateKnowledgeChunk(ctx context.Context, chunk *model.KnowledgeChunk) (*model.KnowledgeChunk, error)
	ListKnowledgeChunks(ctx context.Context, documentID int64, page, pageSize int32) ([]*model.KnowledgeChunk, int64, error)
	DeleteKnowledgeChunks(ctx context.Context, documentID int64) error

	// 知识搜索
	SearchKnowledge(ctx context.Context, kbID int64, query string, limit int32, threshold float64, filters map[string]string) ([]*model.KnowledgeChunk, error)
	HybridSearch(ctx context.Context, kbID int64, query string, limit int32, semanticWeight, keywordWeight, threshold float64, filters map[string]string) ([]*HybridSearchResult, error)

	// 统计和分析
	GetKnowledgeBaseStats(ctx context.Context, kbID int64) (*KnowledgeBaseStats, error)
	UpdateKnowledgeBaseStats(ctx context.Context, kbID int64, stats KnowledgeBaseStatsUpdate) error

	// 处理任务
	CreateProcessingJob(ctx context.Context, job *model.ProcessingJob) (*model.ProcessingJob, error)
	UpdateProcessingJob(ctx context.Context, job *model.ProcessingJob) error
	GetProcessingJob(ctx context.Context, id string) (*model.ProcessingJob, error)
	ListProcessingJobs(ctx context.Context, kbID int64, status []int, limit int32) ([]*model.ProcessingJob, error)
}

// KnowledgeBaseFilter 知识库过滤条件
type KnowledgeBaseFilter struct {
	Keyword  string
	Status   int32
	Tags     []string
	SortBy   string
	SortDesc bool
}

// DocumentFilter 文档过滤条件
type DocumentFilter struct {
	Keyword  string
	Status   int32
	Tags     []string
	MimeType string
	SortBy   string
	SortDesc bool
}

// KnowledgeBaseStats 知识库统计信息
type KnowledgeBaseStats struct {
	DocumentCount        int64
	ChunkCount           int64
	TotalCharacters      int64
	TotalTokens          int64
	AverageChunkSize     float64
	FileTypeDistribution map[string]int64
	LanguageDistribution map[string]int64
	StorageSize          float64
	IndexSize            float64
	LastUpdated          time.Time
}

// KnowledgeBaseStatsUpdate 知识库统计更新
type KnowledgeBaseStatsUpdate struct {
	DocumentCountDelta int64
	ChunkCountDelta    int64
	CharactersDelta    int64
	TokensDelta        int64
	StorageSizeDelta   float64
	IndexSizeDelta     float64
}

// HybridSearchResult 混合搜索结果
type HybridSearchResult struct {
	Chunk         *model.KnowledgeChunk
	SemanticScore float64
	KeywordScore  float64
	CombinedScore float64
}

// DocumentUploadInfo 文档上传信息
type DocumentUploadInfo struct {
	Name        string
	Content     []byte
	MimeType    string
	Metadata    model.DocumentMetadata
	Tags        []string
	AutoProcess bool
}

// NewKnowledgeUsecase 创建知识库业务逻辑实例
func NewKnowledgeUsecase(repo KnowledgeRepo, logger log.Logger) *KnowledgeUsecase {
	return &KnowledgeUsecase{
		repo:   repo,
		logger: log.NewHelper(logger),
	}
}

// CreateKnowledgeBase 创建知识库
func (uc *KnowledgeUsecase) CreateKnowledgeBase(ctx context.Context, userID int64, name, description, embeddingModel string, chunkSize, chunkOverlap int32, config model.KnowledgeBaseConfig, tags []string, language string) (*model.KnowledgeBase, error) {
	now := time.Now()

	kb := &model.KnowledgeBase{
		UserID:         userID,
		Name:           name,
		Description:    description,
		EmbeddingModel: embeddingModel,
		ChunkSize:      int(chunkSize),
		ChunkOverlap:   int(chunkOverlap),
		Status:         1, // active
		Tags:           model.StringSlice(tags),
		Version:        1,
		Language:       language,
		Config:         config,
		AutoProcess:    true,
		CreatedAt:      now,
		UpdatedAt:      now,
	}

	// 设置默认支持的文件类型
	if len(kb.SupportedFileTypes) == 0 {
		kb.SupportedFileTypes = model.StringSlice{"txt", "pdf", "docx", "md", "html"}
	}

	return uc.repo.CreateKnowledgeBase(ctx, kb)
}

// GetKnowledgeBase 获取知识库
func (uc *KnowledgeUsecase) GetKnowledgeBase(ctx context.Context, id int64) (*model.KnowledgeBase, error) {
	return uc.repo.GetKnowledgeBase(ctx, id)
}

// UpdateKnowledgeBase 更新知识库
func (uc *KnowledgeUsecase) UpdateKnowledgeBase(ctx context.Context, id int64, name, description, embeddingModel string, chunkSize, chunkOverlap int32, config model.KnowledgeBaseConfig, tags []string, reindexAfterUpdate bool) (*model.KnowledgeBase, error) {
	kb, err := uc.repo.GetKnowledgeBase(ctx, id)
	if err != nil {
		return nil, err
	}

	// 检查是否需要重新索引
	needReindex := false
	if embeddingModel != "" && embeddingModel != kb.EmbeddingModel {
		kb.EmbeddingModel = embeddingModel
		needReindex = true
	}
	if chunkSize > 0 && int(chunkSize) != kb.ChunkSize {
		kb.ChunkSize = int(chunkSize)
		needReindex = true
	}
	if chunkOverlap > 0 && int(chunkOverlap) != kb.ChunkOverlap {
		kb.ChunkOverlap = int(chunkOverlap)
		needReindex = true
	}

	// 更新其他字段
	if name != "" {
		kb.Name = name
	}
	if description != "" {
		kb.Description = description
	}
	if len(tags) > 0 {
		kb.Tags = model.StringSlice(tags)
	}
	if config.EmbeddingDimension > 0 {
		kb.Config = config
		needReindex = true
	}

	kb.UpdatedAt = time.Now()

	// 如果需要重新索引且用户请求重新索引
	if needReindex && reindexAfterUpdate {
		kb.Status = 2 // indexing
		// TODO: 启动重新索引任务
		uc.logger.Infow("knowledge base configuration changed, reindexing required", "kb_id", id)
	}

	return uc.repo.UpdateKnowledgeBase(ctx, kb)
}

// DeleteKnowledgeBase 删除知识库
func (uc *KnowledgeUsecase) DeleteKnowledgeBase(ctx context.Context, id int64, hardDelete bool) error {
	return uc.repo.DeleteKnowledgeBase(ctx, id, hardDelete)
}

// ListKnowledgeBases 获取知识库列表
func (uc *KnowledgeUsecase) ListKnowledgeBases(ctx context.Context, userID int64, page, pageSize int32, keyword string, status int32, tags []string, sortBy string, sortDesc bool) ([]*model.KnowledgeBase, int64, error) {
	filter := KnowledgeBaseFilter{
		Keyword:  keyword,
		Status:   status,
		Tags:     tags,
		SortBy:   sortBy,
		SortDesc: sortDesc,
	}

	return uc.repo.ListKnowledgeBases(ctx, userID, page, pageSize, filter)
}

// UploadDocument 上传文档
func (uc *KnowledgeUsecase) UploadDocument(ctx context.Context, kbID int64, uploadInfo DocumentUploadInfo) (*model.Document, error) {
	// 生成文档哈希
	hash := uc.generateContentHash(uploadInfo.Content)

	now := time.Now()
	doc := &model.Document{
		KnowledgeBaseID: kbID,
		Name:            uploadInfo.Name,
		Content:         string(uploadInfo.Content),
		MimeType:        uploadInfo.MimeType,
		FileSize:        int64(len(uploadInfo.Content)),
		Status:          1, // uploaded
		Tags:            model.StringSlice(uploadInfo.Tags),
		Hash:            hash,
		Version:         1,
		Metadata:        uploadInfo.Metadata,
		CreatedAt:       now,
		UpdatedAt:       now,
	}

	// 检测文档语言
	doc.Language = uc.detectLanguage(string(uploadInfo.Content))

	document, err := uc.repo.CreateDocument(ctx, doc)
	if err != nil {
		return nil, err
	}

	// 如果启用自动处理，开始处理文档
	if uploadInfo.AutoProcess {
		err = uc.processDocument(ctx, document.ID, false)
		if err != nil {
			uc.logger.Warnw("failed to start document processing", "doc_id", document.ID, "error", err)
		}
	}

	return document, nil
}

// ProcessDocument 处理文档
func (uc *KnowledgeUsecase) ProcessDocument(ctx context.Context, documentID int64, forceReprocess bool) error {
	return uc.processDocument(ctx, documentID, forceReprocess)
}

func (uc *KnowledgeUsecase) processDocument(ctx context.Context, documentID int64, forceReprocess bool) error {
	doc, err := uc.repo.GetDocument(ctx, documentID)
	if err != nil {
		return err
	}

	// 检查是否需要处理
	if !forceReprocess && doc.Status == 3 { // already processed
		return nil
	}

	// 创建处理任务
	job := &model.ProcessingJob{
		ID:              uc.generateJobID(),
		DocumentID:      documentID,
		KnowledgeBaseID: doc.KnowledgeBaseID,
		JobType:         "process",
		Status:          1, // pending
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
	}

	_, err = uc.repo.CreateProcessingJob(ctx, job)
	if err != nil {
		return err
	}

	// TODO: 异步处理文档
	// 这里应该启动一个异步任务来处理文档，包括：
	// 1. 文档分块
	// 2. 生成嵌入向量
	// 3. 存储到向量数据库

	uc.logger.Infow("document processing job created", "job_id", job.ID, "doc_id", documentID)

	return nil
}

// SearchKnowledge 在知识库中搜索
func (uc *KnowledgeUsecase) SearchKnowledge(ctx context.Context, kbID int64, query string, limit int32, threshold float64, filters map[string]string, includeMetadata bool) ([]*model.KnowledgeChunk, error) {
	chunks, err := uc.repo.SearchKnowledge(ctx, kbID, query, limit, threshold, filters)
	if err != nil {
		return nil, err
	}

	// 记录搜索历史
	// TODO: 实现搜索历史记录

	return chunks, nil
}

// HybridSearch 混合搜索
func (uc *KnowledgeUsecase) HybridSearch(ctx context.Context, kbID int64, query string, limit int32, semanticWeight, keywordWeight, threshold float64, filters map[string]string) ([]*HybridSearchResult, error) {
	return uc.repo.HybridSearch(ctx, kbID, query, limit, semanticWeight, keywordWeight, threshold, filters)
}

// GetKnowledgeBaseStats 获取知识库统计信息
func (uc *KnowledgeUsecase) GetKnowledgeBaseStats(ctx context.Context, kbID int64, includeDetailed bool) (*KnowledgeBaseStats, error) {
	return uc.repo.GetKnowledgeBaseStats(ctx, kbID)
}

// ReindexKnowledgeBase 重新索引知识库
func (uc *KnowledgeUsecase) ReindexKnowledgeBase(ctx context.Context, kbID int64, forceReindex bool) (string, error) {
	// 创建重新索引任务
	job := &model.ProcessingJob{
		ID:              uc.generateJobID(),
		KnowledgeBaseID: kbID,
		JobType:         "reindex",
		Status:          1, // pending
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
		Options: map[string]interface{}{
			"force_reindex": forceReindex,
		},
	}

	_, err := uc.repo.CreateProcessingJob(ctx, job)
	if err != nil {
		return "", err
	}

	// TODO: 启动异步重新索引任务

	return job.ID, nil
}

// DeleteDocument 删除文档
func (uc *KnowledgeUsecase) DeleteDocument(ctx context.Context, documentIDs []int64, hardDelete bool) error {
	// 如果是硬删除，需要先删除相关的知识块
	if hardDelete {
		for _, docID := range documentIDs {
			err := uc.repo.DeleteKnowledgeChunks(ctx, docID)
			if err != nil {
				uc.logger.Warnw("failed to delete knowledge chunks", "doc_id", docID, "error", err)
			}
		}
	}

	return uc.repo.DeleteDocument(ctx, documentIDs, hardDelete)
}

// 辅助方法
func (uc *KnowledgeUsecase) generateContentHash(content []byte) string {
	// TODO: 实现内容哈希生成
	return "mock_hash"
}

func (uc *KnowledgeUsecase) detectLanguage(content string) string {
	// TODO: 实现语言检测
	// 可以使用第三方库如 lingua-go
	return "zh"
}

func (uc *KnowledgeUsecase) generateJobID() string {
	// TODO: 生成唯一的任务ID
	return "job_" + time.Now().Format("20060102150405")
}
