package biz

import (
	"context"

	"github.com/go-kratos/kratos/v2/log"
)

// Ai is an Ai model.
type Ai struct {
	Hello string
}

// AiRepo is an Ai repo.
type AiRepo interface {
	CreateAi(context.Context, *Ai) (*Ai, error)
}

// AiUsecase is an Ai usecase.
type AiUsecase struct {
	repo AiRepo
	log  *log.Helper
}

// NewAiUsecase new an Ai usecase.
func NewAiUsecase(repo AiRepo, logger log.Logger) *AiUsecase {
	return &AiUsecase{repo: repo, log: log.NewHelper(logger)}
}

// CreateAi creates an Ai, and returns the new Ai.
func (uc *AiUsecase) CreateAi(ctx context.Context, g *Ai) (*Ai, error) {
	uc.log.WithContext(ctx).Infof("CreateAi: %v", g.Hello)
	return uc.repo.CreateAi(ctx, g)
}
