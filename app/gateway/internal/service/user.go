package service

import (
	"context"
	pb "universal/api/gateway/v1"
	"universal/app/gateway/internal/biz"
)

type UserService struct {
	pb.UnimplementedUserServer

	uc *biz.UserUsecase
}

func NewUserService(uc *biz.UserUsecase) *UserService {
	return &UserService{uc: uc}
}

func (s *UserService) CreateUser(ctx context.Context, req *pb.CreateUserRequest) (*pb.CreateUserReply, error) {
	return &pb.CreateUserReply{}, nil
}
func (s *UserService) UpdateUser(ctx context.Context, req *pb.UpdateUserRequest) (*pb.UpdateUserReply, error) {
	return &pb.UpdateUserReply{}, nil
}
func (s *UserService) DeleteUser(ctx context.Context, req *pb.DeleteUserRequest) (*pb.OperationReply, error) {
	return &pb.OperationReply{}, nil
}
func (s *UserService) GetUser(ctx context.Context, req *pb.GetUserRequest) (*pb.GetUserReply, error) {
	return &pb.GetUserReply{}, nil
}
func (s *UserService) ListUser(ctx context.Context, req *pb.ListUserRequest) (*pb.ListUserReply, error) {
	user, err := s.uc.ListUser(ctx, nil)
	if err != nil {
		return nil, err
	}
	users := make([]*pb.UserInfo, len(user.Users))
	for i, u := range user.Users {
		users[i] = &pb.UserInfo{
			Id:       u.ID,
			Username: u.Username,
			Email:    u.Email,
			Phone:    u.Phone,
			Nickname: u.Nickname,
			Avatar:   u.Avatar,
			Status:   u.Status,
		}
	}
	return &pb.ListUserReply{
		Page:     user.Page,
		PageSize: user.PageSize,
		Total:    user.Total,
		Users:    users,
	}, nil
}
func (s *UserService) BatchDeleteUser(ctx context.Context, req *pb.BatchDeleteUserRequest) (*pb.BatchDeleteUserReply, error) {
	return &pb.BatchDeleteUserReply{}, nil
}
func (s *UserService) UpdateUserStatus(ctx context.Context, req *pb.UpdateUserStatusRequest) (*pb.OperationReply, error) {
	return &pb.OperationReply{}, nil
}
func (s *UserService) ChangePassword(ctx context.Context, req *pb.ChangePasswordRequest) (*pb.OperationReply, error) {
	return &pb.OperationReply{}, nil
}
func (s *UserService) GetUserStats(ctx context.Context, req *pb.GetUserStatsRequest) (*pb.GetUserStatsReply, error) {
	return &pb.GetUserStatsReply{}, nil
}
