package service

import (
	"context"
	"github.com/go-kratos/kratos/v2/log"
	"tinyid/internal/biz"

	pb "tinyid/api/idgen/v1"
)

type IdgenService struct {
	pb.UnimplementedIdgenServer
	uc  *biz.IdgenUsecase
	log *log.Helper
}

func NewIdgenService(uc *biz.IdgenUsecase, logger log.Logger) *IdgenService {
	return &IdgenService{uc: uc, log: log.NewHelper("service/idgen", logger)}
}

func (s *IdgenService) SegmentId(ctx context.Context, req *pb.SegmentRequest) (*pb.IdReply, error) {
	return &pb.IdReply{}, nil
}
func (s *IdgenService) SnowflakeId(ctx context.Context, req *pb.SnowflakeRequest) (*pb.IdReply, error) {
	return &pb.IdReply{}, nil
}
func (s *IdgenService) CurrentTime(ctx context.Context, req *pb.CurrentTimeRequest) (*pb.CurrentTimeReply, error) {
	return &pb.CurrentTimeReply{}, nil
}
