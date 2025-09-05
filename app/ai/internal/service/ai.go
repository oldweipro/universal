package service

import (
	"universal/app/ai/internal/biz"

	pb "universal/api/ai/v1"
)

type AiService struct {
	pb.UnimplementedAiServer
	uc *biz.AiUsecase
}

func NewAiService(uc *biz.AiUsecase) *AiService {
	return &AiService{uc: uc}
}
