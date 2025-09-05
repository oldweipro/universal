package model

import (
	"encoding/json"
	"time"

	"gorm.io/gorm"
)

// KnowledgeBase 知识库模型
type KnowledgeBase struct {
	ID                 int64               `gorm:"primarykey" json:"id"`
	UserID             int64               `gorm:"not null;index" json:"user_id"`
	Name               string              `gorm:"size:255;not null" json:"name"`
	Description        string              `gorm:"type:text" json:"description"`
	EmbeddingModel     string              `gorm:"size:100;not null" json:"embedding_model"`
	ChunkSize          int                 `gorm:"default:1000" json:"chunk_size"`
	ChunkOverlap       int                 `gorm:"default:200" json:"chunk_overlap"`
	Status             int                 `gorm:"default:1;index" json:"status"` // 1:active, 2:indexing, 3:error, 4:archived, 5:deleted
	Tags               StringSlice         `gorm:"type:json" json:"tags"`
	Version            int                 `gorm:"default:1" json:"version"`
	Language           string              `gorm:"size:10;default:'zh'" json:"language"`
	SupportedFileTypes StringSlice         `gorm:"type:json" json:"supported_file_types"`
	MaxFileSize        int64               `gorm:"default:104857600" json:"max_file_size"` // 100MB default
	AutoProcess        bool                `gorm:"default:true" json:"auto_process"`
	Config             KnowledgeBaseConfig `gorm:"type:json" json:"config"`
	CreatedAt          time.Time           `json:"created_at"`
	UpdatedAt          time.Time           `json:"updated_at"`
	LastIndexedAt      *time.Time          `json:"last_indexed_at"`
	DeletedAt          gorm.DeletedAt      `gorm:"index" json:"deleted_at"`

	// 统计信息
	DocumentCount   int64   `gorm:"default:0" json:"document_count"`
	ChunkCount      int64   `gorm:"default:0" json:"chunk_count"`
	TotalCharacters int64   `gorm:"default:0" json:"total_characters"`
	TotalTokens     int64   `gorm:"default:0" json:"total_tokens"`
	StorageSize     float64 `gorm:"default:0" json:"storage_size"` // MB
	IndexSize       float64 `gorm:"default:0" json:"index_size"`   // MB

	// 关联关系
	Documents []Document `gorm:"foreignKey:KnowledgeBaseID" json:"documents,omitempty"`
}

// KnowledgeBaseConfig 知识库配置
type KnowledgeBaseConfig struct {
	EmbeddingDimension   int                    `json:"embedding_dimension,omitempty"`
	ChunkingStrategy     string                 `json:"chunking_strategy,omitempty"` // fixed_size, recursive, semantic, paragraph, sentence
	SimilarityThreshold  float64                `json:"similarity_threshold,omitempty"`
	MaxChunksPerQuery    int                    `json:"max_chunks_per_query,omitempty"`
	EnableMetadataFilter bool                   `json:"enable_metadata_filter,omitempty"`
	StopWords            []string               `json:"stop_words,omitempty"`
	TextSplitter         string                 `json:"text_splitter,omitempty"`
	AdvancedOptions      map[string]interface{} `json:"advanced_options,omitempty"`
}

// Document 文档模型
type Document struct {
	ID                 int64            `gorm:"primarykey" json:"id"`
	KnowledgeBaseID    int64            `gorm:"not null;index" json:"knowledge_base_id"`
	Name               string           `gorm:"size:255;not null" json:"name"`
	Content            string           `gorm:"type:longtext" json:"content"`
	FilePath           string           `gorm:"size:500" json:"file_path"`
	MimeType           string           `gorm:"size:100" json:"mime_type"`
	FileSize           int64            `gorm:"default:0" json:"file_size"`
	ChunkCount         int              `gorm:"default:0" json:"chunk_count"`
	Status             int              `gorm:"default:1;index" json:"status"` // 1:uploaded, 2:processing, 3:processed, 4:failed, 5:archived, 6:deleted
	Tags               StringSlice      `gorm:"type:json" json:"tags"`
	Language           string           `gorm:"size:10" json:"language"`
	SourceURL          string           `gorm:"size:500" json:"source_url"`
	Hash               string           `gorm:"size:64;index" json:"hash"`
	Version            int              `gorm:"default:1" json:"version"`
	ProcessingProgress float64          `gorm:"default:0" json:"processing_progress"`
	ProcessingError    string           `gorm:"type:text" json:"processing_error"`
	Metadata           DocumentMetadata `gorm:"type:json" json:"metadata"`
	CreatedAt          time.Time        `json:"created_at"`
	UpdatedAt          time.Time        `json:"updated_at"`
	LastProcessedAt    *time.Time       `json:"last_processed_at"`
	DeletedAt          gorm.DeletedAt   `gorm:"index" json:"deleted_at"`

	// 关联关系
	KnowledgeBase KnowledgeBase    `gorm:"foreignKey:KnowledgeBaseID" json:"knowledge_base,omitempty"`
	Chunks        []KnowledgeChunk `gorm:"foreignKey:DocumentID" json:"chunks,omitempty"`
}

// DocumentMetadata 文档元数据
type DocumentMetadata struct {
	Title        string            `json:"title,omitempty"`
	Author       string            `json:"author,omitempty"`
	Subject      string            `json:"subject,omitempty"`
	Keywords     []string          `json:"keywords,omitempty"`
	CreatedDate  *time.Time        `json:"created_date,omitempty"`
	ModifiedDate *time.Time        `json:"modified_date,omitempty"`
	Category     string            `json:"category,omitempty"`
	PageCount    int               `json:"page_count,omitempty"`
	Encoding     string            `json:"encoding,omitempty"`
	CustomFields map[string]string `json:"custom_fields,omitempty"`
}

// KnowledgeChunk 知识块模型
type KnowledgeChunk struct {
	ID              int64           `gorm:"primarykey" json:"id"`
	DocumentID      int64           `gorm:"not null;index" json:"document_id"`
	KnowledgeBaseID int64           `gorm:"not null;index" json:"knowledge_base_id"`
	Content         string          `gorm:"type:text;not null" json:"content"`
	ChunkIndex      int             `gorm:"not null" json:"chunk_index"`
	StartPosition   int             `gorm:"default:0" json:"start_position"`
	EndPosition     int             `gorm:"default:0" json:"end_position"`
	Language        string          `gorm:"size:10" json:"language"`
	Keywords        StringSlice     `gorm:"type:json" json:"keywords"`
	ChunkType       int             `gorm:"default:1" json:"chunk_type"` // 1:text, 2:title, 3:paragraph, 4:list, 5:table, 6:code
	CharacterCount  int             `gorm:"default:0" json:"character_count"`
	TokenCount      int             `gorm:"default:0" json:"token_count"`
	Embedding       EmbeddingVector `gorm:"type:json" json:"embedding"`
	Metadata        ChunkMetadata   `gorm:"type:json" json:"metadata"`
	CreatedAt       time.Time       `json:"created_at"`
	UpdatedAt       time.Time       `json:"updated_at"`

	// 搜索相关字段（运行时使用）
	Score float64 `gorm:"-" json:"score,omitempty"`

	// 关联关系
	Document      Document      `gorm:"foreignKey:DocumentID" json:"document,omitempty"`
	KnowledgeBase KnowledgeBase `gorm:"foreignKey:KnowledgeBaseID" json:"knowledge_base,omitempty"`
}

// ChunkMetadata 块元数据
type ChunkMetadata struct {
	Section    string            `json:"section,omitempty"`
	PageNumber int               `json:"page_number,omitempty"`
	LineNumber int               `json:"line_number,omitempty"`
	Confidence float64           `json:"confidence,omitempty"`
	Source     string            `json:"source,omitempty"`
	Custom     map[string]string `json:"custom,omitempty"`
}

// EmbeddingVector 向量嵌入类型
type EmbeddingVector []float64

func (e EmbeddingVector) Value() (interface{}, error) {
	if len(e) == 0 {
		return "[]", nil
	}
	b, err := json.Marshal(e)
	return string(b), err
}

func (e *EmbeddingVector) Scan(value interface{}) error {
	if value == nil {
		*e = EmbeddingVector{}
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

	return json.Unmarshal(bytes, e)
}

// ProcessingJob 处理任务模型
type ProcessingJob struct {
	ID              string                 `gorm:"primarykey" json:"id"`
	DocumentID      int64                  `gorm:"not null;index" json:"document_id"`
	KnowledgeBaseID int64                  `gorm:"not null;index" json:"knowledge_base_id"`
	JobType         string                 `gorm:"size:50;not null" json:"job_type"` // process, reindex, analyze
	Status          int                    `gorm:"default:1;index" json:"status"`    // 1:pending, 2:running, 3:completed, 4:failed, 5:cancelled
	Progress        float64                `gorm:"default:0" json:"progress"`
	ErrorMessage    string                 `gorm:"type:text" json:"error_message"`
	Result          JobResult              `gorm:"type:json" json:"result"`
	Options         map[string]interface{} `gorm:"type:json" json:"options"`
	CreatedAt       time.Time              `json:"created_at"`
	UpdatedAt       time.Time              `json:"updated_at"`
	StartedAt       *time.Time             `json:"started_at"`
	CompletedAt     *time.Time             `json:"completed_at"`

	// 关联关系
	Document      Document      `gorm:"foreignKey:DocumentID" json:"document,omitempty"`
	KnowledgeBase KnowledgeBase `gorm:"foreignKey:KnowledgeBaseID" json:"knowledge_base,omitempty"`
}

// JobResult 任务结果
type JobResult struct {
	ChunksCreated   int                    `json:"chunks_created,omitempty"`
	TokensProcessed int64                  `json:"tokens_processed,omitempty"`
	ProcessingTime  float64                `json:"processing_time,omitempty"`
	Warnings        []string               `json:"warnings,omitempty"`
	Metrics         map[string]interface{} `json:"metrics,omitempty"`
}

// SearchHistory 搜索历史模型
type SearchHistory struct {
	ID              int64       `gorm:"primarykey" json:"id"`
	UserID          int64       `gorm:"not null;index" json:"user_id"`
	KnowledgeBaseID int64       `gorm:"not null;index" json:"knowledge_base_id"`
	Query           string      `gorm:"type:text;not null" json:"query"`
	SearchType      string      `gorm:"size:50;not null" json:"search_type"` // semantic, keyword, hybrid, advanced
	ResultCount     int         `gorm:"default:0" json:"result_count"`
	MaxScore        float64     `gorm:"default:0" json:"max_score"`
	ResponseTime    float64     `gorm:"default:0" json:"response_time"`
	Filters         KeyValueMap `gorm:"type:json" json:"filters"`
	CreatedAt       time.Time   `json:"created_at"`

	// 关联关系
	KnowledgeBase KnowledgeBase `gorm:"foreignKey:KnowledgeBaseID" json:"knowledge_base,omitempty"`
}

// 自定义GORM类型转换器实现
func (c KnowledgeBaseConfig) Value() (interface{}, error) {
	b, err := json.Marshal(c)
	return string(b), err
}

func (c *KnowledgeBaseConfig) Scan(value interface{}) error {
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

func (m DocumentMetadata) Value() (interface{}, error) {
	b, err := json.Marshal(m)
	return string(b), err
}

func (m *DocumentMetadata) Scan(value interface{}) error {
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

func (m ChunkMetadata) Value() (interface{}, error) {
	b, err := json.Marshal(m)
	return string(b), err
}

func (m *ChunkMetadata) Scan(value interface{}) error {
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

func (r JobResult) Value() (interface{}, error) {
	b, err := json.Marshal(r)
	return string(b), err
}

func (r *JobResult) Scan(value interface{}) error {
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

	return json.Unmarshal(bytes, r)
}

// TableName 返回表名
func (KnowledgeBase) TableName() string {
	return "knowledge_bases"
}

func (Document) TableName() string {
	return "documents"
}

func (KnowledgeChunk) TableName() string {
	return "knowledge_chunks"
}

func (ProcessingJob) TableName() string {
	return "processing_jobs"
}

func (SearchHistory) TableName() string {
	return "search_histories"
}
