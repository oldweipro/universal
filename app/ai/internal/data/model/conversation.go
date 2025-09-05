package model

import (
	"encoding/json"
	"time"

	"gorm.io/gorm"
)

// Conversation 对话记录模型
type Conversation struct {
	ID               int64              `gorm:"primarykey" json:"id"`
	UserID           int64              `gorm:"not null;index" json:"user_id"`
	Title            string             `gorm:"size:255;not null" json:"title"`
	ModelName        string             `gorm:"size:100;not null" json:"model_name"`
	SystemPrompt     string             `gorm:"type:text" json:"system_prompt"`
	Config           ConversationConfig `gorm:"type:json" json:"config"`
	Status           int                `gorm:"default:1;index" json:"status"` // 1:active, 2:archived, 3:deleted, 4:paused
	Description      string             `gorm:"type:text" json:"description"`
	Tags             StringSlice        `gorm:"type:json" json:"tags"`
	Priority         int                `gorm:"default:0" json:"priority"`
	AutoArchiveAfter int64              `gorm:"default:0" json:"auto_archive_after"` // seconds
	CreatedAt        time.Time          `json:"created_at"`
	UpdatedAt        time.Time          `json:"updated_at"`
	LastActiveAt     time.Time          `gorm:"index" json:"last_active_at"`
	DeletedAt        gorm.DeletedAt     `gorm:"index" json:"deleted_at"`

	// 统计信息
	MessageCount      int64 `gorm:"default:0" json:"message_count"`
	TotalInputTokens  int64 `gorm:"default:0" json:"total_input_tokens"`
	TotalOutputTokens int64 `gorm:"default:0" json:"total_output_tokens"`
	ToolCallCount     int64 `gorm:"default:0" json:"tool_call_count"`
	TotalDuration     int64 `gorm:"default:0" json:"total_duration"` // milliseconds

	// 关联关系
	Messages []Message `gorm:"foreignKey:ConversationID" json:"messages,omitempty"`
}

// ConversationConfig 对话配置
type ConversationConfig struct {
	Temperature      *float64               `json:"temperature,omitempty"`
	MaxTokens        *int                   `json:"max_tokens,omitempty"`
	TopP             *float64               `json:"top_p,omitempty"`
	FrequencyPenalty *float64               `json:"frequency_penalty,omitempty"`
	PresencePenalty  *float64               `json:"presence_penalty,omitempty"`
	StopSequences    []string               `json:"stop_sequences,omitempty"`
	CustomParams     map[string]interface{} `json:"custom_params,omitempty"`
}

// ConversationMemory 对话记忆
type ConversationMemory struct {
	ID              int64       `gorm:"primarykey" json:"id"`
	ConversationID  int64       `gorm:"not null;uniqueIndex" json:"conversation_id"`
	Summary         string      `gorm:"type:text" json:"summary"`
	KeyPoints       StringSlice `gorm:"type:json" json:"key_points"`
	UserPreferences KeyValueMap `gorm:"type:json" json:"user_preferences"`
	ImportantFacts  StringSlice `gorm:"type:json" json:"important_facts"`
	MemoryVersion   int         `gorm:"default:1" json:"memory_version"`
	CreatedAt       time.Time   `json:"created_at"`
	UpdatedAt       time.Time   `json:"updated_at"`

	// 关联关系
	Conversation Conversation `gorm:"foreignKey:ConversationID" json:"conversation,omitempty"`
}

// Message 消息记录模型
type Message struct {
	ID              int64      `gorm:"primarykey" json:"id"`
	ConversationID  int64      `gorm:"not null;index" json:"conversation_id"`
	Role            string     `gorm:"size:20;not null;index" json:"role"` // user, assistant, system, tool
	Content         string     `gorm:"type:longtext;not null" json:"content"`
	Status          int        `gorm:"default:1;index" json:"status"` // 1:pending, 2:processing, 3:completed, 4:failed, 5:deleted
	ParentMessageID *int64     `gorm:"index" json:"parent_message_id"`
	IsEdited        bool       `gorm:"default:false" json:"is_edited"`
	EditReason      string     `gorm:"size:255" json:"edit_reason"`
	EditedAt        *time.Time `json:"edited_at"`

	// Token统计
	InputTokens  int     `gorm:"default:0" json:"input_tokens"`
	OutputTokens int     `gorm:"default:0" json:"output_tokens"`
	ResponseTime float64 `gorm:"default:0" json:"response_time"` // seconds
	Cost         float64 `gorm:"default:0" json:"cost"`
	ModelUsed    string  `gorm:"size:100" json:"model_used"`

	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at"`

	// 关联关系
	Conversation Conversation        `gorm:"foreignKey:ConversationID" json:"conversation,omitempty"`
	ToolCalls    []ToolCall          `gorm:"foreignKey:MessageID" json:"tool_calls,omitempty"`
	Attachments  []MessageAttachment `gorm:"foreignKey:MessageID" json:"attachments,omitempty"`
}

// ToolCall 工具调用记录
type ToolCall struct {
	ID            string      `gorm:"primarykey" json:"id"`
	MessageID     int64       `gorm:"not null;index" json:"message_id"`
	Name          string      `gorm:"size:100;not null" json:"name"`
	Arguments     string      `gorm:"type:text" json:"arguments"`
	Result        string      `gorm:"type:text" json:"result"`
	Status        int         `gorm:"default:1;index" json:"status"` // 1:pending, 2:running, 3:success, 4:failed, 5:timeout, 6:cancelled
	ErrorMessage  string      `gorm:"type:text" json:"error_message"`
	ExecutionTime int64       `gorm:"default:0" json:"execution_time"` // milliseconds
	RetryCount    int         `gorm:"default:0" json:"retry_count"`
	Metadata      KeyValueMap `gorm:"type:json" json:"metadata"`
	CreatedAt     time.Time   `json:"created_at"`
	UpdatedAt     time.Time   `json:"updated_at"`

	// 关联关系
	Message Message `gorm:"foreignKey:MessageID" json:"message,omitempty"`
}

// MessageAttachment 消息附件
type MessageAttachment struct {
	ID        int64       `gorm:"primarykey" json:"id"`
	MessageID int64       `gorm:"not null;index" json:"message_id"`
	Name      string      `gorm:"size:255;not null" json:"name"`
	URL       string      `gorm:"size:500;not null" json:"url"`
	MimeType  string      `gorm:"size:100" json:"mime_type"`
	Size      int64       `gorm:"default:0" json:"size"`
	Type      int         `gorm:"default:1" json:"type"` // 1:image, 2:document, 3:audio, 4:video, 5:code
	Metadata  KeyValueMap `gorm:"type:json" json:"metadata"`
	CreatedAt time.Time   `json:"created_at"`
	UpdatedAt time.Time   `json:"updated_at"`

	// 关联关系
	Message Message `gorm:"foreignKey:MessageID" json:"message,omitempty"`
}

// StringSlice 字符串切片的自定义类型，用于JSON序列化
type StringSlice []string

func (s StringSlice) Value() (interface{}, error) {
	if len(s) == 0 {
		return "[]", nil
	}
	b, err := json.Marshal(s)
	return string(b), err
}

func (s *StringSlice) Scan(value interface{}) error {
	if value == nil {
		*s = StringSlice{}
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

	return json.Unmarshal(bytes, s)
}

// KeyValueMap 键值对映射的自定义类型
type KeyValueMap map[string]string

func (kv KeyValueMap) Value() (interface{}, error) {
	if len(kv) == 0 {
		return "{}", nil
	}
	b, err := json.Marshal(kv)
	return string(b), err
}

func (kv *KeyValueMap) Scan(value interface{}) error {
	if value == nil {
		*kv = KeyValueMap{}
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

	return json.Unmarshal(bytes, kv)
}

// TableName 返回表名
func (Conversation) TableName() string {
	return "conversations"
}

func (ConversationMemory) TableName() string {
	return "conversation_memories"
}

func (Message) TableName() string {
	return "messages"
}

func (ToolCall) TableName() string {
	return "tool_calls"
}

func (MessageAttachment) TableName() string {
	return "message_attachments"
}
