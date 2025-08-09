package service

import (
	"context"

	pb "pet-angel/api/auth/v1"
)

type AuthServiceService struct {
	pb.UnimplementedAuthServiceServer
}

func NewAuthServiceService() *AuthServiceService {
	return &AuthServiceService{}
}

func (s *AuthServiceService) Login(ctx context.Context, req *pb.LoginRequest) (*pb.LoginReply, error) {
	return &pb.LoginReply{}, nil
}
func (s *AuthServiceService) Relogin(ctx context.Context, req *pb.ReloginRequest) (*pb.ReloginReply, error) {
	return &pb.ReloginReply{}, nil
}
func (s *AuthServiceService) GetUserInfo(ctx context.Context, req *pb.GetUserInfoRequest) (*pb.GetUserInfoReply, error) {
	return &pb.GetUserInfoReply{}, nil
}
func (s *AuthServiceService) UpdateUserInfo(ctx context.Context, req *pb.UpdateUserInfoRequest) (*pb.UpdateUserInfoReply, error) {
	return &pb.UpdateUserInfoReply{}, nil
}
