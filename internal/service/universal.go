package service

import (
	"context"
	"universal/internal/biz"

	pb "universal/api/universal/v1"
)

type UniversalService struct {
	pb.UnimplementedUniversalServer

	uc *biz.UserUsecase
}

func NewUniversalService(uc *biz.UserUsecase) *UniversalService {
	return &UniversalService{uc: uc}
}

func (s *UniversalService) CreateUser(ctx context.Context, req *pb.CreateUserRequest) (*pb.CreateUserReply, error) {
	return &pb.CreateUserReply{}, nil
}
func (s *UniversalService) UpdateUser(ctx context.Context, req *pb.UpdateUserRequest) (*pb.UpdateUserReply, error) {
	return &pb.UpdateUserReply{}, nil
}
func (s *UniversalService) DeleteUser(ctx context.Context, req *pb.DeleteUserRequest) (*pb.OperationReply, error) {
	return &pb.OperationReply{}, nil
}
func (s *UniversalService) GetUser(ctx context.Context, req *pb.GetUserRequest) (*pb.GetUserReply, error) {
	return &pb.GetUserReply{}, nil
}
func (s *UniversalService) ListUser(ctx context.Context, req *pb.ListUserRequest) (*pb.ListUserReply, error) {
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
func (s *UniversalService) BatchDeleteUser(ctx context.Context, req *pb.BatchDeleteUserRequest) (*pb.BatchDeleteUserReply, error) {
	return &pb.BatchDeleteUserReply{}, nil
}
func (s *UniversalService) UpdateUserStatus(ctx context.Context, req *pb.UpdateUserStatusRequest) (*pb.OperationReply, error) {
	return &pb.OperationReply{}, nil
}
func (s *UniversalService) ChangePassword(ctx context.Context, req *pb.ChangePasswordRequest) (*pb.OperationReply, error) {
	return &pb.OperationReply{}, nil
}
func (s *UniversalService) GetUserStats(ctx context.Context, req *pb.GetUserStatsRequest) (*pb.GetUserStatsReply, error) {
	return &pb.GetUserStatsReply{}, nil
}
