package biz

import (
	"context"
	pb "universal/api/system/v1"

	"github.com/go-kratos/kratos/v2/log"
)

// User is a User model.
type User struct {
	Hello string
}

// UserRepo is a User repo.
type UserRepo interface {
	Save(context.Context, *User) (*User, error)
	Update(context.Context, *User) (*User, error)
	FindByID(context.Context, int64) (*User, error)
	ListByHello(context.Context, string) ([]*User, error)
	ListAll(context.Context) ([]*pb.UserInfo, error)
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

// CreateUser creates a User, and returns the new User.
func (uc *UserUsecase) CreateUser(ctx context.Context, g *User) (*User, error) {
	uc.log.WithContext(ctx).Infof("CreateUser: %v", g.Hello)
	return uc.repo.Save(ctx, g)
}

// ListUser list User, and returns the User list.
func (uc *UserUsecase) ListUser(ctx context.Context, req *pb.ListUserRequest) ([]*pb.UserInfo, error) {
	uc.log.WithContext(ctx).Infof("ListUser: %v", req)
	return uc.repo.ListAll(ctx)
}
