package biz

import (
	"context"
	"github.com/go-kratos/kratos/v2/log"
	"tinyid/pkg/snowflake"
)

type IdgenRepo interface {
	WorkerID() int64
}

type IdgenUsecase struct {
	repo      IdgenRepo
	log       *log.Helper
	snowflake *snowflake.SnowFlake
}

func NewIdgenUsecase(repo IdgenRepo, logger log.Logger) *IdgenUsecase {
	return &IdgenUsecase{
		repo:      repo,
		log:       log.NewHelper("usecase/Idgen", logger),
		snowflake: snowflake.NewSnowFlake(repo.WorkerID()),
	}
}

func (uc *IdgenUsecase) GetSnowflakeID(ctx context.Context) (int64, error) {
	return uc.snowflake.GenID()
}

func (uc *IdgenUsecase) GetSegmentID(ctx context.Context, tag string) (int64, error) {
	return 0, nil
}
