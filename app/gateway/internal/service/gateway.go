package service

import (
	"context"
	pb "universal/api/gateway/v1"
	"universal/app/gateway/internal/biz"
)

type GatewayService struct {
	pb.UnimplementedGatewayServer

	uc *biz.UserUsecase
}

func NewGatewayService(uc *biz.UserUsecase) *GatewayService {
	return &GatewayService{uc: uc}
}

func (s *GatewayService) CreateUser(ctx context.Context, req *pb.CreateUserRequest) (*pb.CreateUserReply, error) {
	return &pb.CreateUserReply{}, nil
}
func (s *GatewayService) UpdateUser(ctx context.Context, req *pb.UpdateUserRequest) (*pb.UpdateUserReply, error) {
	return &pb.UpdateUserReply{}, nil
}
func (s *GatewayService) DeleteUser(ctx context.Context, req *pb.DeleteUserRequest) (*pb.OperationReply, error) {
	return &pb.OperationReply{}, nil
}
func (s *GatewayService) GetUser(ctx context.Context, req *pb.GetUserRequest) (*pb.GetUserReply, error) {
	return &pb.GetUserReply{}, nil
}
func (s *GatewayService) ListUser(ctx context.Context, req *pb.ListUserRequest) (*pb.ListUserReply, error) {
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
func (s *GatewayService) BatchDeleteUser(ctx context.Context, req *pb.BatchDeleteUserRequest) (*pb.BatchDeleteUserReply, error) {
	return &pb.BatchDeleteUserReply{}, nil
}
func (s *GatewayService) UpdateUserStatus(ctx context.Context, req *pb.UpdateUserStatusRequest) (*pb.OperationReply, error) {
	return &pb.OperationReply{}, nil
}
func (s *GatewayService) ChangePassword(ctx context.Context, req *pb.ChangePasswordRequest) (*pb.OperationReply, error) {
	return &pb.OperationReply{}, nil
}
func (s *GatewayService) GetUserStats(ctx context.Context, req *pb.GetUserStatsRequest) (*pb.GetUserStatsReply, error) {
	return &pb.GetUserStatsReply{}, nil
}
