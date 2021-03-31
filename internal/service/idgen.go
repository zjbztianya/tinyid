package service

import (
	"context"
	"github.com/go-kratos/kratos/v2/log"
	"time"
	"tinyid/internal/biz"

	durpb "google.golang.org/protobuf/types/known/durationpb"
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

func (s *IdgenService) SegmentID(ctx context.Context, req *pb.SegmentRequest) (*pb.IdReply, error) {
	return &pb.IdReply{Id: 0}, nil
}

func (s *IdgenService) SnowflakeID(ctx context.Context, req *pb.SnowflakeRequest) (*pb.IdReply, error) {
	id, err := s.uc.Get(ctx, "")
	s.log.Debugf("generator snowflake id %d", id)
	if err != nil {
		return nil, err
	}
	return &pb.IdReply{Id: id}, nil
}

func (s *IdgenService) CurrentTime(ctx context.Context, req *pb.CurrentTimeRequest) (*pb.CurrentTimeReply, error) {
	return &pb.CurrentTimeReply{Time: durpb.New(time.Duration(time.Now().UnixNano()))}, nil
}
