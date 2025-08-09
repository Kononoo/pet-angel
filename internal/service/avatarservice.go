package service

import (
	"context"

	pb "pet-angel/api/avatar/v1"
)

type AvatarServiceService struct {
	pb.UnimplementedAvatarServiceServer
}

func NewAvatarServiceService() *AvatarServiceService {
	return &AvatarServiceService{}
}

func (s *AvatarServiceService) GetModels(ctx context.Context, req *pb.GetModelsRequest) (*pb.GetModelsReply, error) {
	return &pb.GetModelsReply{}, nil
}
func (s *AvatarServiceService) SetPetModel(ctx context.Context, req *pb.SetPetModelRequest) (*pb.SetPetModelReply, error) {
	return &pb.SetPetModelReply{}, nil
}
func (s *AvatarServiceService) GetItems(ctx context.Context, req *pb.GetItemsRequest) (*pb.GetItemsReply, error) {
	return &pb.GetItemsReply{}, nil
}
func (s *AvatarServiceService) UseItem(ctx context.Context, req *pb.UseItemRequest) (*pb.UseItemReply, error) {
	return &pb.UseItemReply{}, nil
}
func (s *AvatarServiceService) Chat(ctx context.Context, req *pb.ChatRequest) (*pb.ChatReply, error) {
	return &pb.ChatReply{}, nil
}
