package model

// User 用户数据库模型
type User struct {
	BaseModel
	Username string `gorm:"uniqueIndex;not null;size:50" json:"username"`
	Email    string `gorm:"uniqueIndex;not null;size:100" json:"email"`
	Phone    string `gorm:"size:20" json:"phone"`
	Password string `gorm:"not null;size:255" json:"-"` // 密码不返回给前端
	Nickname string `gorm:"size:50" json:"nickname"`
	Avatar   string `gorm:"size:255" json:"avatar"`
	Status   int32  `gorm:"default:1;comment:状态 1:正常 0:禁用" json:"status"`
}

// TableName 自定义表名
func (User) TableName() string {
	return "users"
}
