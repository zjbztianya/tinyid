package biz

import (
	"github.com/go-kratos/kratos/v2/log"
)

type Idgen struct {
	Hello string
}

type IdgenRepo interface {
	CreateIdgen(*Idgen) error
	UpdateIdgen(*Idgen) error
}

type IdgenUsecase struct {
	repo IdgenRepo
	log  *log.Helper
}

func NewIdgenUsecase(repo IdgenRepo, logger log.Logger) *IdgenUsecase {
	return &IdgenUsecase{repo: repo, log: log.NewHelper("usecase/Idgen", logger)}
}

func (uc *IdgenUsecase) Create(g *Idgen) error {
	return uc.repo.CreateIdgen(g)
}

func (uc *IdgenUsecase) Update(g *Idgen) error {
	return uc.repo.UpdateIdgen(g)
}
