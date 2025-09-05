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

type knowledgeRepo struct {
	data   *Data
	logger *log.Helper
}

// NewKnowledgeRepo 创建知识库仓库实例
func NewKnowledgeRepo(data *Data, logger log.Logger) biz.KnowledgeRepo {
	return &knowledgeRepo{
		data:   data,
		logger: log.NewHelper(logger),
	}
}

// CreateKnowledgeBase 创建知识库
func (r *knowledgeRepo) CreateKnowledgeBase(ctx context.Context, kb *model.KnowledgeBase) (*model.KnowledgeBase, error) {
	if err := r.data.db.WithContext(ctx).Create(kb).Error; err != nil {
		return nil, err
	}
	return kb, nil
}

// GetKnowledgeBase 获取知识库
func (r *knowledgeRepo) GetKnowledgeBase(ctx context.Context, id int64) (*model.KnowledgeBase, error) {
	var kb model.KnowledgeBase
	err := r.data.db.WithContext(ctx).Where("id = ?", id).First(&kb).Error
	if err != nil {
		return nil, err
	}
	return &kb, nil
}

// UpdateKnowledgeBase 更新知识库
func (r *knowledgeRepo) UpdateKnowledgeBase(ctx context.Context, kb *model.KnowledgeBase) (*model.KnowledgeBase, error) {
	if err := r.data.db.WithContext(ctx).Save(kb).Error; err != nil {
		return nil, err
	}
	return kb, nil
}

// DeleteKnowledgeBase 删除知识库
func (r *knowledgeRepo) DeleteKnowledgeBase(ctx context.Context, id int64, hardDelete bool) error {
	if hardDelete {
		// 硬删除：彻底删除所有相关数据
		return r.data.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
			// 删除知识块
			if err := tx.Unscoped().Where("knowledge_base_id = ?", id).Delete(&model.KnowledgeChunk{}).Error; err != nil {
				return err
			}
			// 删除文档
			if err := tx.Unscoped().Where("knowledge_base_id = ?", id).Delete(&model.Document{}).Error; err != nil {
				return err
			}
			// 删除处理任务
			if err := tx.Unscoped().Where("knowledge_base_id = ?", id).Delete(&model.ProcessingJob{}).Error; err != nil {
				return err
			}
			// 删除知识库
			if err := tx.Unscoped().Delete(&model.KnowledgeBase{}, id).Error; err != nil {
				return err
			}
			return nil
		})
	} else {
		// 软删除：更新状态
		return r.data.db.WithContext(ctx).Model(&model.KnowledgeBase{}).Where("id = ?", id).Update("status", 5).Error
	}
}

// ListKnowledgeBases 获取知识库列表
func (r *knowledgeRepo) ListKnowledgeBases(ctx context.Context, userID int64, page, pageSize int32, filters biz.KnowledgeBaseFilter) ([]*model.KnowledgeBase, int64, error) {
	var kbs []*model.KnowledgeBase
	var total int64

	query := r.data.db.WithContext(ctx).Model(&model.KnowledgeBase{}).Where("user_id = ?", userID)

	// 应用过滤条件
	if filters.Keyword != "" {
		query = query.Where("name LIKE ? OR description LIKE ?",
			fmt.Sprintf("%%%s%%", filters.Keyword),
			fmt.Sprintf("%%%s%%", filters.Keyword))
	}

	if filters.Status > 0 {
		query = query.Where("status = ?", filters.Status)
	}

	if len(filters.Tags) > 0 {
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
	if err := query.Offset(int(offset)).Limit(int(pageSize)).Find(&kbs).Error; err != nil {
		return nil, 0, err
	}

	return kbs, total, nil
}

// CreateDocument 创建文档
func (r *knowledgeRepo) CreateDocument(ctx context.Context, doc *model.Document) (*model.Document, error) {
	return r.createDocumentWithTx(ctx, r.data.db, doc)
}

func (r *knowledgeRepo) createDocumentWithTx(ctx context.Context, tx *gorm.DB, doc *model.Document) (*model.Document, error) {
	if err := tx.WithContext(ctx).Create(doc).Error; err != nil {
		return nil, err
	}

	// 更新知识库的文档数量统计
	updates := map[string]interface{}{
		"document_count":   gorm.Expr("document_count + 1"),
		"total_characters": gorm.Expr("total_characters + ?", len(doc.Content)),
		"storage_size":     gorm.Expr("storage_size + ?", float64(doc.FileSize)/1024/1024),
		"updated_at":       time.Now(),
	}

	if err := tx.WithContext(ctx).Model(&model.KnowledgeBase{}).Where("id = ?", doc.KnowledgeBaseID).Updates(updates).Error; err != nil {
		r.logger.Warnw("failed to update knowledge base stats", "error", err)
	}

	return doc, nil
}

// GetDocument 获取文档
func (r *knowledgeRepo) GetDocument(ctx context.Context, id int64) (*model.Document, error) {
	var doc model.Document
	err := r.data.db.WithContext(ctx).Where("id = ?", id).First(&doc).Error
	if err != nil {
		return nil, err
	}
	return &doc, nil
}

// UpdateDocument 更新文档
func (r *knowledgeRepo) UpdateDocument(ctx context.Context, doc *model.Document) (*model.Document, error) {
	if err := r.data.db.WithContext(ctx).Save(doc).Error; err != nil {
		return nil, err
	}
	return doc, nil
}

// DeleteDocument 删除文档
func (r *knowledgeRepo) DeleteDocument(ctx context.Context, documentIDs []int64, hardDelete bool) error {
	if len(documentIDs) == 0 {
		return nil
	}

	if hardDelete {
		return r.data.db.WithContext(ctx).Unscoped().Delete(&model.Document{}, documentIDs).Error
	} else {
		return r.data.db.WithContext(ctx).Model(&model.Document{}).Where("id IN ?", documentIDs).Update("status", 6).Error
	}
}

// ListDocuments 获取文档列表
func (r *knowledgeRepo) ListDocuments(ctx context.Context, kbID int64, page, pageSize int32, filters biz.DocumentFilter) ([]*model.Document, int64, error) {
	var docs []*model.Document
	var total int64

	query := r.data.db.WithContext(ctx).Model(&model.Document{}).Where("knowledge_base_id = ?", kbID)

	// 应用过滤条件
	if filters.Keyword != "" {
		query = query.Where("name LIKE ? OR content LIKE ?",
			fmt.Sprintf("%%%s%%", filters.Keyword),
			fmt.Sprintf("%%%s%%", filters.Keyword))
	}

	if filters.Status > 0 {
		query = query.Where("status = ?", filters.Status)
	}

	if filters.MimeType != "" {
		query = query.Where("mime_type = ?", filters.MimeType)
	}

	if len(filters.Tags) > 0 {
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
	if err := query.Offset(int(offset)).Limit(int(pageSize)).Find(&docs).Error; err != nil {
		return nil, 0, err
	}

	return docs, total, nil
}

// CreateKnowledgeChunks 创建知识块
func (r *knowledgeRepo) CreateKnowledgeChunks(ctx context.Context, chunks []*model.KnowledgeChunk) error {
	if len(chunks) == 0 {
		return nil
	}

	return r.data.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// 批量创建知识块
		if err := tx.CreateInBatches(chunks, 100).Error; err != nil {
			return err
		}

		// 更新文档的块数量
		if len(chunks) > 0 {
			documentID := chunks[0].DocumentID
			chunkCount := len(chunks)

			if err := tx.Model(&model.Document{}).Where("id = ?", documentID).Updates(map[string]interface{}{
				"chunk_count": gorm.Expr("chunk_count + ?", chunkCount),
				"updated_at":  time.Now(),
			}).Error; err != nil {
				return err
			}

			// 更新知识库的块数量统计
			kbID := chunks[0].KnowledgeBaseID
			totalTokens := int64(0)
			for _, chunk := range chunks {
				totalTokens += int64(chunk.TokenCount)
			}

			if err := tx.Model(&model.KnowledgeBase{}).Where("id = ?", kbID).Updates(map[string]interface{}{
				"chunk_count":  gorm.Expr("chunk_count + ?", chunkCount),
				"total_tokens": gorm.Expr("total_tokens + ?", totalTokens),
				"updated_at":   time.Now(),
			}).Error; err != nil {
				r.logger.Warnw("failed to update knowledge base chunk stats", "error", err)
			}
		}

		return nil
	})
}

// GetKnowledgeChunk 获取知识块
func (r *knowledgeRepo) GetKnowledgeChunk(ctx context.Context, id int64) (*model.KnowledgeChunk, error) {
	var chunk model.KnowledgeChunk
	err := r.data.db.WithContext(ctx).Where("id = ?", id).First(&chunk).Error
	if err != nil {
		return nil, err
	}
	return &chunk, nil
}

// UpdateKnowledgeChunk 更新知识块
func (r *knowledgeRepo) UpdateKnowledgeChunk(ctx context.Context, chunk *model.KnowledgeChunk) (*model.KnowledgeChunk, error) {
	if err := r.data.db.WithContext(ctx).Save(chunk).Error; err != nil {
		return nil, err
	}
	return chunk, nil
}

// ListKnowledgeChunks 获取知识块列表
func (r *knowledgeRepo) ListKnowledgeChunks(ctx context.Context, documentID int64, page, pageSize int32) ([]*model.KnowledgeChunk, int64, error) {
	var chunks []*model.KnowledgeChunk
	var total int64

	query := r.data.db.WithContext(ctx).Model(&model.KnowledgeChunk{}).Where("document_id = ?", documentID)

	// 获取总数
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 分页
	offset := (page - 1) * pageSize
	if err := query.Order("chunk_index ASC").Offset(int(offset)).Limit(int(pageSize)).Find(&chunks).Error; err != nil {
		return nil, 0, err
	}

	return chunks, total, nil
}

// DeleteKnowledgeChunks 删除知识块
func (r *knowledgeRepo) DeleteKnowledgeChunks(ctx context.Context, documentID int64) error {
	return r.data.db.WithContext(ctx).Where("document_id = ?", documentID).Delete(&model.KnowledgeChunk{}).Error
}

// SearchKnowledge 知识搜索
func (r *knowledgeRepo) SearchKnowledge(ctx context.Context, kbID int64, query string, limit int32, threshold float64, filters map[string]string) ([]*model.KnowledgeChunk, error) {
	// TODO: 实现向量搜索
	// 这里需要集成向量数据库（如Milvus、Pinecone等）进行语义搜索
	// 暂时使用关键词搜索作为占位符

	var chunks []*model.KnowledgeChunk

	dbQuery := r.data.db.WithContext(ctx).Model(&model.KnowledgeChunk{}).
		Where("knowledge_base_id = ?", kbID).
		Where("content LIKE ?", fmt.Sprintf("%%%s%%", query)).
		Order("created_at DESC").
		Limit(int(limit))

	// 应用过滤器
	for key, value := range filters {
		switch key {
		case "language":
			dbQuery = dbQuery.Where("language = ?", value)
		case "chunk_type":
			dbQuery = dbQuery.Where("chunk_type = ?", value)
		}
	}

	if err := dbQuery.Find(&chunks).Error; err != nil {
		return nil, err
	}

	// 模拟相似度分数
	for i, _ := range chunks {
		chunks[i].Score = 0.8 - float64(i)*0.1
	}

	return chunks, nil
}

// HybridSearch 混合搜索
func (r *knowledgeRepo) HybridSearch(ctx context.Context, kbID int64, query string, limit int32, semanticWeight, keywordWeight, threshold float64, filters map[string]string) ([]*biz.HybridSearchResult, error) {
	// TODO: 实现混合搜索（语义搜索 + 关键词搜索）
	// 这里使用简化的实现

	// 先进行语义搜索
	semanticChunks, err := r.SearchKnowledge(ctx, kbID, query, limit, threshold, filters)
	if err != nil {
		return nil, err
	}

	// 进行关键词搜索
	var keywordChunks []*model.KnowledgeChunk
	keywordQuery := r.data.db.WithContext(ctx).Model(&model.KnowledgeChunk{}).
		Where("knowledge_base_id = ?", kbID).
		Where("content LIKE ? OR JSON_EXTRACT(keywords, '$') LIKE ?",
			fmt.Sprintf("%%%s%%", query),
			fmt.Sprintf("%%%s%%", query)).
		Order("created_at DESC").
		Limit(int(limit))

	if err := keywordQuery.Find(&keywordChunks).Error; err != nil {
		return nil, err
	}

	// 合并结果并计算综合分数
	resultMap := make(map[int64]*biz.HybridSearchResult)

	// 处理语义搜索结果
	for _, chunk := range semanticChunks {
		result := &biz.HybridSearchResult{
			Chunk:         chunk,
			SemanticScore: chunk.Score,
			KeywordScore:  0,
		}
		resultMap[chunk.ID] = result
	}

	// 处理关键词搜索结果
	for i, chunk := range keywordChunks {
		keywordScore := 0.9 - float64(i)*0.1
		if result, exists := resultMap[chunk.ID]; exists {
			result.KeywordScore = keywordScore
		} else {
			result := &biz.HybridSearchResult{
				Chunk:         chunk,
				SemanticScore: 0,
				KeywordScore:  keywordScore,
			}
			resultMap[chunk.ID] = result
		}
	}

	// 计算综合分数并转换为切片
	var results []*biz.HybridSearchResult
	for _, result := range resultMap {
		result.CombinedScore = result.SemanticScore*semanticWeight + result.KeywordScore*keywordWeight
		if result.CombinedScore >= threshold {
			results = append(results, result)
		}
	}

	// 按综合分数排序
	// TODO: 实现排序逻辑

	// 限制结果数量
	if len(results) > int(limit) {
		results = results[:limit]
	}

	return results, nil
}

// GetKnowledgeBaseStats 获取知识库统计信息
func (r *knowledgeRepo) GetKnowledgeBaseStats(ctx context.Context, kbID int64) (*biz.KnowledgeBaseStats, error) {
	var kb model.KnowledgeBase
	if err := r.data.db.WithContext(ctx).Where("id = ?", kbID).First(&kb).Error; err != nil {
		return nil, err
	}

	// 计算文件类型分布
	var fileTypeStats []struct {
		MimeType string `json:"mime_type"`
		Count    int64  `json:"count"`
	}
	r.data.db.WithContext(ctx).Model(&model.Document{}).
		Where("knowledge_base_id = ?", kbID).
		Select("mime_type, COUNT(*) as count").
		Group("mime_type").
		Scan(&fileTypeStats)

	fileTypeDistribution := make(map[string]int64)
	for _, stat := range fileTypeStats {
		fileTypeDistribution[stat.MimeType] = stat.Count
	}

	// 计算语言分布
	var languageStats []struct {
		Language string `json:"language"`
		Count    int64  `json:"count"`
	}
	r.data.db.WithContext(ctx).Model(&model.Document{}).
		Where("knowledge_base_id = ?", kbID).
		Select("language, COUNT(*) as count").
		Group("language").
		Scan(&languageStats)

	languageDistribution := make(map[string]int64)
	for _, stat := range languageStats {
		languageDistribution[stat.Language] = stat.Count
	}

	// 计算平均块大小
	var avgChunkSize float64
	r.data.db.WithContext(ctx).Model(&model.KnowledgeChunk{}).
		Where("knowledge_base_id = ?", kbID).
		Select("AVG(character_count)").
		Scan(&avgChunkSize)

	stats := &biz.KnowledgeBaseStats{
		DocumentCount:        kb.DocumentCount,
		ChunkCount:           kb.ChunkCount,
		TotalCharacters:      kb.TotalCharacters,
		TotalTokens:          kb.TotalTokens,
		AverageChunkSize:     avgChunkSize,
		FileTypeDistribution: fileTypeDistribution,
		LanguageDistribution: languageDistribution,
		StorageSize:          kb.StorageSize,
		IndexSize:            kb.IndexSize,
		LastUpdated:          kb.UpdatedAt,
	}

	return stats, nil
}

// UpdateKnowledgeBaseStats 更新知识库统计信息
func (r *knowledgeRepo) UpdateKnowledgeBaseStats(ctx context.Context, kbID int64, stats biz.KnowledgeBaseStatsUpdate) error {
	updates := map[string]interface{}{
		"updated_at": time.Now(),
	}

	if stats.DocumentCountDelta != 0 {
		updates["document_count"] = gorm.Expr("document_count + ?", stats.DocumentCountDelta)
	}
	if stats.ChunkCountDelta != 0 {
		updates["chunk_count"] = gorm.Expr("chunk_count + ?", stats.ChunkCountDelta)
	}
	if stats.CharactersDelta != 0 {
		updates["total_characters"] = gorm.Expr("total_characters + ?", stats.CharactersDelta)
	}
	if stats.TokensDelta != 0 {
		updates["total_tokens"] = gorm.Expr("total_tokens + ?", stats.TokensDelta)
	}
	if stats.StorageSizeDelta != 0 {
		updates["storage_size"] = gorm.Expr("storage_size + ?", stats.StorageSizeDelta)
	}
	if stats.IndexSizeDelta != 0 {
		updates["index_size"] = gorm.Expr("index_size + ?", stats.IndexSizeDelta)
	}

	return r.data.db.WithContext(ctx).Model(&model.KnowledgeBase{}).Where("id = ?", kbID).Updates(updates).Error
}

// CreateProcessingJob 创建处理任务
func (r *knowledgeRepo) CreateProcessingJob(ctx context.Context, job *model.ProcessingJob) (*model.ProcessingJob, error) {
	if err := r.data.db.WithContext(ctx).Create(job).Error; err != nil {
		return nil, err
	}
	return job, nil
}

// UpdateProcessingJob 更新处理任务
func (r *knowledgeRepo) UpdateProcessingJob(ctx context.Context, job *model.ProcessingJob) error {
	return r.data.db.WithContext(ctx).Save(job).Error
}

// GetProcessingJob 获取处理任务
func (r *knowledgeRepo) GetProcessingJob(ctx context.Context, id string) (*model.ProcessingJob, error) {
	var job model.ProcessingJob
	err := r.data.db.WithContext(ctx).Where("id = ?", id).First(&job).Error
	if err != nil {
		return nil, err
	}
	return &job, nil
}

// ListProcessingJobs 获取处理任务列表
func (r *knowledgeRepo) ListProcessingJobs(ctx context.Context, kbID int64, status []int, limit int32) ([]*model.ProcessingJob, error) {
	var jobs []*model.ProcessingJob

	query := r.data.db.WithContext(ctx).Model(&model.ProcessingJob{}).Where("knowledge_base_id = ?", kbID)

	if len(status) > 0 {
		query = query.Where("status IN ?", status)
	}

	if err := query.Order("created_at DESC").Limit(int(limit)).Find(&jobs).Error; err != nil {
		return nil, err
	}

	return jobs, nil
}
