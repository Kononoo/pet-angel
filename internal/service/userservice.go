package service

import (
	"context"

	pb "pet-angel/api/user/v1"
)

type UserServiceService struct {
	pb.UnimplementedUserServiceServer
}

func NewUserServiceService() *UserServiceService {
	return &UserServiceService{}
}

func (s *UserServiceService) FollowUser(ctx context.Context, req *pb.FollowUserRequest) (*pb.FollowUserReply, error) {
	return &pb.FollowUserReply{}, nil
}
func (s *UserServiceService) UnfollowUser(ctx context.Context, req *pb.UnfollowUserRequest) (*pb.UnfollowUserReply, error) {
	return &pb.UnfollowUserReply{}, nil
}
func (s *UserServiceService) GetUserProfile(ctx context.Context, req *pb.GetUserProfileRequest) (*pb.GetUserProfileReply, error) {
	return &pb.GetUserProfileReply{}, nil
}
func (s *UserServiceService) GetFollowList(ctx context.Context, req *pb.GetFollowListRequest) (*pb.GetFollowListReply, error) {
	return &pb.GetFollowListReply{}, nil
}
func (s *UserServiceService) GetLikeList(ctx context.Context, req *pb.GetLikeListRequest) (*pb.GetLikeListReply, error) {
	return &pb.GetLikeListReply{}, nil
}
