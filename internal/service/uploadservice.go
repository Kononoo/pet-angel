package service

import (
	"context"

	pb "pet-angel/api/upload/v1"
)

type UploadServiceService struct {
	pb.UnimplementedUploadServiceServer
}

func NewUploadServiceService() *UploadServiceService {
	return &UploadServiceService{}
}

func (s *UploadServiceService) GetPresign(ctx context.Context, req *pb.GetPresignRequest) (*pb.GetPresignReply, error) {
	return &pb.GetPresignReply{}, nil
}
func (s *UploadServiceService) UploadDone(ctx context.Context, req *pb.UploadDoneRequest) (*pb.UploadDoneReply, error) {
	return &pb.UploadDoneReply{}, nil
}
