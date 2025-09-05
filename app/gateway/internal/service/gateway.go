package service

import (
	"context"

	pb "universal/api/gateway/v1"
)

type GatewayService struct {
	pb.UnimplementedGatewayServer
}

func NewGatewayService() *GatewayService {
	return &GatewayService{}
}

func (s *GatewayService) GetGatewayInfo(ctx context.Context, req *pb.GetGatewayInfoRequest) (*pb.GetGatewayInfoReply, error) {
	return &pb.GetGatewayInfoReply{}, nil
}
func (s *GatewayService) GetGatewayHealth(ctx context.Context, req *pb.GetGatewayHealthRequest) (*pb.GetGatewayHealthReply, error) {
	return &pb.GetGatewayHealthReply{}, nil
}
