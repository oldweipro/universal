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
}

// NewAiUsecase new an Ai usecase.
func NewAiUsecase(repo AiRepo) *AiUsecase {
	return &AiUsecase{repo: repo}
}

// CreateAi creates an Ai, and returns the new Ai.
func (uc *AiUsecase) CreateAi(ctx context.Context, g *Ai) (*Ai, error) {
	log.Infof("CreateAi: %v", g.Hello)
	return uc.repo.CreateAi(ctx, g)
}
