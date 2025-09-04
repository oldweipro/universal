package biz

import (
	"context"

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
	List(context.Context, *ListUserRequest) (*ListUserResponse, error) // 分页列表查询
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

// ListUser 分页查询用户列表
func (uc *UserUsecase) ListUser(ctx context.Context, req *ListUserRequest) (*ListUserResponse, error) {
	uc.log.WithContext(ctx).Infof("ListUser: %+v", req)
	return uc.repo.List(ctx, req)
}
