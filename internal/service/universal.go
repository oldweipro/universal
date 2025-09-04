package service

import (
	"context"

	pb "universal/api/universal/v1"
)

type UniversalService struct {
	pb.UnimplementedUniversalServer
}

func NewUniversalService() *UniversalService {
	return &UniversalService{}
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
	return &pb.ListUserReply{Page: 1, PageSize: 10}, nil
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
