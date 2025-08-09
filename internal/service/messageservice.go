package service

import (
	"context"

	pb "pet-angel/api/message/v1"
)

type MessageServiceService struct {
	pb.UnimplementedMessageServiceServer
}

func NewMessageServiceService() *MessageServiceService {
	return &MessageServiceService{}
}

func (s *MessageServiceService) GetMessageList(ctx context.Context, req *pb.GetMessageListRequest) (*pb.GetMessageListReply, error) {
	return &pb.GetMessageListReply{}, nil
}
func (s *MessageServiceService) UnlockMessage(ctx context.Context, req *pb.UnlockMessageRequest) (*pb.UnlockMessageReply, error) {
	return &pb.UnlockMessageReply{}, nil
}
