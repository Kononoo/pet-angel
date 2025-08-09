package server

import (
	authv1 "pet-angel/api/auth/v1"
	avatv1 "pet-angel/api/avatar/v1"
	communityv1 "pet-angel/api/community/v1"
	v1 "pet-angel/api/helloworld/v1"
	msgv1 "pet-angel/api/message/v1"
	userv1 "pet-angel/api/user/v1"
	"pet-angel/internal/conf"
	"pet-angel/internal/service"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware/recovery"
	"github.com/go-kratos/kratos/v2/transport/grpc"
)

// NewGRPCServer new a gRPC server.
func NewGRPCServer(c *conf.Server, greeter *service.GreeterService, auth *service.AuthService, user *service.UserService, community *service.CommunityService, avatar *service.AvatarService, message *service.MessageService, logger log.Logger) *grpc.Server {
	var opts = []grpc.ServerOption{
		grpc.Middleware(
			recovery.Recovery(),
		),
	}
	if c.Grpc.Network != "" {
		opts = append(opts, grpc.Network(c.Grpc.Network))
	}
	if c.Grpc.Addr != "" {
		opts = append(opts, grpc.Address(c.Grpc.Addr))
	}
	if c.Grpc.Timeout != nil {
		opts = append(opts, grpc.Timeout(c.Grpc.Timeout.AsDuration()))
	}

	srv := grpc.NewServer(opts...)
	v1.RegisterGreeterServer(srv, greeter)
	authv1.RegisterAuthServiceServer(srv, auth)
	userv1.RegisterUserServiceServer(srv, user)
	communityv1.RegisterCommunityServiceServer(srv, community)
	avatv1.RegisterAvatarServiceServer(srv, avatar)
	msgv1.RegisterMessageServiceServer(srv, message)
	return srv
}
