package data

import (
	"context"
	userv1 "universal/api/user/v1"
	"universal/internal/biz"

	"github.com/go-kratos/kratos/v2/log"
)

type userRepo struct {
	data *Data
	log  *log.Helper
}

// NewUserRepo .
func NewUserRepo(data *Data, logger log.Logger) biz.UserRepo {
	return &userRepo{
		data: data,
		log:  log.NewHelper(log.With(logger, "module", "data/universal/user")),
	}
}

func (r *userRepo) List(ctx context.Context, req *biz.ListUserRequest) (*biz.ListUserResponse, error) {
	listUser, err := r.data.uc.ListUser(ctx, &userv1.ListUserRequest{Page: 1, PageSize: 10})
	if err != nil {
		return nil, err
	}
	users := make([]*biz.User, len(listUser.Users))
	for i, u := range listUser.Users {
		users[i] = &biz.User{
			ID:       u.Id,
			Username: u.Username,
			Email:    u.Email,
			Phone:    u.Phone,
			Nickname: u.Nickname,
			Avatar:   u.Avatar,
			Status:   u.Status,
		}
	}
	return &biz.ListUserResponse{
		Users:    users,
		Total:    listUser.Total,
		Page:     listUser.Page,
		PageSize: listUser.PageSize,
	}, nil
}
