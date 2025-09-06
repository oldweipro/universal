package service

import (
	"context"
	pb "universal/api/gateway/v1"
	"universal/app/gateway/internal/biz"

	"github.com/go-kratos/kratos/v2/log"
)

type GatewayService struct {
	pb.UnimplementedGatewayServer

	uc  *biz.GatewayUsecase
	log *log.Helper
}

func NewGatewayService(uc *biz.GatewayUsecase, logger log.Logger) *GatewayService {
	return &GatewayService{
		uc:  uc,
		log: log.NewHelper(logger),
	}
}

func (s *GatewayService) GetGatewayInfo(ctx context.Context, req *pb.GetGatewayInfoRequest) (*pb.GetGatewayInfoReply, error) {
	s.log.WithContext(ctx).Infof("GetGatewayInfo called")

	info, err := s.uc.GetGatewayInfo(ctx)
	if err != nil {
		s.log.WithContext(ctx).Errorf("Failed to get gateway info: %v", err)
		return nil, err
	}

	return &pb.GetGatewayInfoReply{
		Info: info,
	}, nil
}

func (s *GatewayService) GetGatewayHealth(ctx context.Context, req *pb.GetGatewayHealthRequest) (*pb.GetGatewayHealthReply, error) {
	s.log.WithContext(ctx).Infof("GetGatewayHealth called")

	health, err := s.uc.GetGatewayHealth(ctx)
	if err != nil {
		s.log.WithContext(ctx).Errorf("Failed to get gateway health: %v", err)
		return nil, err
	}

	return health, nil
}
