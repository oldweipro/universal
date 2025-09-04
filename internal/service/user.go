package service

import (
	"context"
	pb "universal/api/system/v1"
	"universal/internal/biz"
)

type UserService struct {
	pb.UnimplementedUserServer
	uc *biz.UserUsecase
}

func NewUserService(uc *biz.UserUsecase) *UserService {
	return &UserService{uc: uc}
}

func (s *UserService) CreateUser(ctx context.Context, req *pb.CreateUserRequest) (*pb.CreateUserReply, error) {
	// 将protobuf请求转换为业务层User
	user := &biz.User{
		Username: req.Username,
		Email:    req.Email,
		Phone:    req.Phone,
		Nickname: req.Nickname,
		Avatar:   req.Avatar,
		Status:   req.Status,
	}

	// 调用业务层创建用户，传入密码
	createdUser, err := s.uc.CreateUser(ctx, user, req.Password)
	if err != nil {
		return nil, err
	}

	// 将业务层User转换为protobuf响应
	userInfo := &pb.UserInfo{
		Id:       createdUser.ID,
		Username: createdUser.Username,
		Email:    createdUser.Email,
		Phone:    createdUser.Phone,
		Nickname: createdUser.Nickname,
		Avatar:   createdUser.Avatar,
		Status:   createdUser.Status,
	}

	return &pb.CreateUserReply{User: userInfo}, nil
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
	// 转换protobuf请求为业务层请求
	bizReq := &biz.ListUserRequest{
		Page:     req.Page,
		PageSize: req.PageSize,
		Keyword:  req.Keyword,
		Status:   req.Status,
	}

	bizResp, err := s.uc.ListUser(ctx, bizReq)
	if err != nil {
		return nil, err
	}

	// 转换为protobuf模型
	var pbUsers []*pb.UserInfo
	for _, bizUser := range bizResp.Users {
		pbUsers = append(pbUsers, &pb.UserInfo{
			Id:       bizUser.ID,
			Username: bizUser.Username,
			Email:    bizUser.Email,
			Phone:    bizUser.Phone,
			Nickname: bizUser.Nickname,
			Avatar:   bizUser.Avatar,
			Status:   bizUser.Status,
		})
	}

	return &pb.ListUserReply{
		Users:    pbUsers,
		Total:    bizResp.Total,
		Page:     bizResp.Page,
		PageSize: bizResp.PageSize,
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
