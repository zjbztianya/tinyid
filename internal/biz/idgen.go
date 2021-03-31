package biz

import (
	"context"
	"github.com/go-kratos/kratos/v2/log"
)

type IdgenRepo interface {
	CreateSnowflakeId() (int64, error)
	CreateSegmentId(tag string) (int64, error)
}

type IdgenUsecase struct {
	repo IdgenRepo
	log  *log.Helper
}

func NewIdgenUsecase(repo IdgenRepo, logger log.Logger) *IdgenUsecase {
	return &IdgenUsecase{repo: repo,
		log: log.NewHelper("usecase/Idgen", logger),
	}
}

func (uc *IdgenUsecase) Get(ctx context.Context, tag string) (int64, error) {
	// tag 非空为segment模式的biz_tag,否则为snowflake
	if tag == "" {
		return uc.repo.CreateSnowflakeId()
	}
	return 0, nil
}
