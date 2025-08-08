package service

import (
	"context"
	"errors"
	v1 "pet-angel/api/helloworld/v1"
	"pet-angel/internal/biz"
)

// GreeterService is a greeter service.
type GreeterService struct {
	v1.UnimplementedGreeterServer

	uc *biz.GreeterUsecase
}

// NewGreeterService new a greeter service.
func NewGreeterService(uc *biz.GreeterUsecase) *GreeterService {
	return &GreeterService{uc: uc}
}

// SayHello implements helloworld.GreeterServer.
func (s *GreeterService) SayHello(ctx context.Context, in *v1.HelloRequest) (*v1.HelloReply, error) {
	g, err := s.uc.CreateGreeter(ctx, &biz.Greeter{Hello: in.Name})
	if err != nil {
		return nil, err
	}
	// 返回业务数据，会被统一响应编码器包装到 HTTP 响应的 data 字段中
	rsp := &v1.HelloReply{
		Message: "Hello " + g.Hello,
		Pet:     "小猫咪～",
		Hello:   "主人你好呀！",
	}
	err = errors.New("user not found")
	return rsp, err
}
