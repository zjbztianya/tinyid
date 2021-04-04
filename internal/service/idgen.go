package service

import (
	"context"
	"github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/log"
	errors2 "github.com/pkg/errors"
	"time"
	"tinyid/internal/biz"

	durpb "google.golang.org/protobuf/types/known/durationpb"
	apierr "tinyid/api/idgen/errors"
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
	id, err := s.uc.GetSegmentID(ctx, req.Tag)
	if err != nil {
		s.log.Errorf("%v", errors2.Cause(err))
		return nil, errors.ResourceExhausted(apierr.Errors_NotGenSegmentID, "tag:%s", req.Tag)
	}
	return &pb.IdReply{Id: id}, nil
}

func (s *IdgenService) SnowflakeID(ctx context.Context, req *pb.SnowflakeRequest) (*pb.IdReply, error) {
	id, err := s.uc.GetSnowflakeID(ctx)
	if err != nil {
		return nil, errors.Internal(apierr.Errors_ClockBackwards, "snowflake")
	}
	return &pb.IdReply{Id: id}, nil
}

func (s *IdgenService) CurrentTime(ctx context.Context, req *pb.CurrentTimeRequest) (*pb.CurrentTimeReply, error) {
	return &pb.CurrentTimeReply{Time: durpb.New(time.Duration(time.Now().UnixNano()))}, nil
}
