package data

import (
	"context"
	"universal/app/ai/internal/biz"

	"github.com/go-kratos/kratos/v2/log"
)

type aiRepo struct {
	data *Data
	log  *log.Helper
}

// NewAiRepo .
func NewAiRepo(data *Data, logger log.Logger) biz.AiRepo {
	return &aiRepo{
		data: data,
		log:  log.NewHelper(logger),
	}
}

func (r *aiRepo) CreateAi(ctx context.Context, g *biz.Ai) (*biz.Ai, error) {
	return g, nil
}
