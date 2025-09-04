package biz

import (
	"context"
	"errors"

	"github.com/go-kratos/kratos/v2/log"
)

// User is a User business model.
type User struct {
	ID       int64  `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Phone    string `json:"phone"`
	Nickname string `json:"nickname"`
	Avatar   string `json:"avatar"`
	Status   int32  `json:"status"`
}

// ListUserRequest 业务层列表查询请求
type ListUserRequest struct {
	Page     int32  `json:"page"`
	PageSize int32  `json:"page_size"`
	Keyword  string `json:"keyword"`
	Status   int32  `json:"status"`
}

// ListUserResponse 业务层列表查询响应
type ListUserResponse struct {
	Users    []*User `json:"users"`
	Total    int64   `json:"total"`
	Page     int32   `json:"page"`
	PageSize int32   `json:"page_size"`
}

// UserRepo is a User repo.
type UserRepo interface {
	// 基本CRUD操作
	Create(context.Context, *User, string) (*User, error) // 创建用户（带密码）
	GetByID(context.Context, int64) (*User, error)        // 根据ID获取用户
	Update(context.Context, *User) (*User, error)         // 更新用户
	Delete(context.Context, int64) error                  // 删除用户

	// 查询操作
	List(context.Context, *ListUserRequest) (*ListUserResponse, error) // 分页列表查询
	GetByUsername(context.Context, string) (*User, error)              // 根据用户名获取
	GetByEmail(context.Context, string) (*User, error)                 // 根据邮箱获取
	ExistsByUsername(context.Context, string) (bool, error)            // 检查用户名是否存在
	ExistsByEmail(context.Context, string) (bool, error)               // 检查邮箱是否存在
}

// UserUsecase is a User usecase.
type UserUsecase struct {
	repo UserRepo
	log  *log.Helper
}

// NewUserUsecase new a User usecase.
func NewUserUsecase(repo UserRepo, logger log.Logger) *UserUsecase {
	return &UserUsecase{repo: repo, log: log.NewHelper(logger)}
}

// CreateUser 创建用户
func (uc *UserUsecase) CreateUser(ctx context.Context, user *User, password string) (*User, error) {
	// 检查用户名是否已存在
	exists, err := uc.repo.ExistsByUsername(ctx, user.Username)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, errors.New("用户名已存在")
	}

	// 检查邮箱是否已存在
	exists, err = uc.repo.ExistsByEmail(ctx, user.Email)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, errors.New("邮箱已存在")
	}

	uc.log.WithContext(ctx).Infof("CreateUser: %v", user.Username)
	return uc.repo.Create(ctx, user, password)
}

// GetUser 根据ID获取用户
func (uc *UserUsecase) GetUser(ctx context.Context, id int64) (*User, error) {
	uc.log.WithContext(ctx).Infof("GetUser: %v", id)
	return uc.repo.GetByID(ctx, id)
}

// UpdateUser 更新用户
func (uc *UserUsecase) UpdateUser(ctx context.Context, user *User) (*User, error) {
	uc.log.WithContext(ctx).Infof("UpdateUser: %v", user.ID)
	return uc.repo.Update(ctx, user)
}

// DeleteUser 删除用户
func (uc *UserUsecase) DeleteUser(ctx context.Context, id int64) error {
	uc.log.WithContext(ctx).Infof("DeleteUser: %v", id)
	return uc.repo.Delete(ctx, id)
}

// ListUser 分页查询用户列表
func (uc *UserUsecase) ListUser(ctx context.Context, req *ListUserRequest) (*ListUserResponse, error) {
	uc.log.WithContext(ctx).Infof("ListUser: %+v", req)
	return uc.repo.List(ctx, req)
}
