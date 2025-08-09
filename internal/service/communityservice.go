package service

import (
	"context"

	pb "pet-angel/api/community/v1"
)

type CommunityServiceService struct {
	pb.UnimplementedCommunityServiceServer
}

func NewCommunityServiceService() *CommunityServiceService {
	return &CommunityServiceService{}
}

func (s *CommunityServiceService) GetCategories(ctx context.Context, req *pb.GetCategoriesRequest) (*pb.GetCategoriesReply, error) {
	return &pb.GetCategoriesReply{}, nil
}
func (s *CommunityServiceService) GetPostList(ctx context.Context, req *pb.GetPostListRequest) (*pb.GetPostListReply, error) {
	return &pb.GetPostListReply{}, nil
}
func (s *CommunityServiceService) GetPostDetail(ctx context.Context, req *pb.GetPostDetailRequest) (*pb.GetPostDetailReply, error) {
	return &pb.GetPostDetailReply{}, nil
}
func (s *CommunityServiceService) CreatePost(ctx context.Context, req *pb.CreatePostRequest) (*pb.CreatePostReply, error) {
	return &pb.CreatePostReply{}, nil
}
func (s *CommunityServiceService) LikePost(ctx context.Context, req *pb.LikePostRequest) (*pb.LikePostReply, error) {
	return &pb.LikePostReply{}, nil
}
func (s *CommunityServiceService) UnlikePost(ctx context.Context, req *pb.UnlikePostRequest) (*pb.UnlikePostReply, error) {
	return &pb.UnlikePostReply{}, nil
}
func (s *CommunityServiceService) GetCommentList(ctx context.Context, req *pb.GetCommentListRequest) (*pb.GetCommentListReply, error) {
	return &pb.GetCommentListReply{}, nil
}
func (s *CommunityServiceService) CreateComment(ctx context.Context, req *pb.CreateCommentRequest) (*pb.CreateCommentReply, error) {
	return &pb.CreateCommentReply{}, nil
}
func (s *CommunityServiceService) LikeComment(ctx context.Context, req *pb.LikeCommentRequest) (*pb.LikeCommentReply, error) {
	return &pb.LikeCommentReply{}, nil
}
func (s *CommunityServiceService) UnlikeComment(ctx context.Context, req *pb.UnlikeCommentRequest) (*pb.UnlikeCommentReply, error) {
	return &pb.UnlikeCommentReply{}, nil
}
