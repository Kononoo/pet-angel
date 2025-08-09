//go:build wireinject
// +build wireinject

// The build tag makes sure the stub is not built in the final build.

package main

import (
	"github.com/go-kratos/kratos/v2"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/google/wire"

	"pet-angel/internal/biz"
	"pet-angel/internal/conf"
	"pet-angel/internal/data"
	"pet-angel/internal/server"
	"pet-angel/internal/service"
)

// wireApp init kratos application.
func wireApp(srv *conf.Server, dataConf *conf.Data, authConf *conf.Auth, logger log.Logger) (*kratos.App, func(), error) {
	panic(wire.Build(
		// data/infrastructure
		data.NewData,

		// repo providers
		data.NewGreeterRepo,
		data.NewAuthRepo,
		data.NewUserRepo,
		data.NewCommunityRepo,

		// interface bindings
		wire.Bind(new(biz.GreeterRepo), new(*data.GreeterRepo)),
		wire.Bind(new(biz.AuthRepo), new(*data.AuthRepo)),
		wire.Bind(new(biz.CommunityRepo), new(*data.CommunityRepoImpl)),

		// biz
		biz.NewGreeterUsecase,
		biz.NewAuthUsecase,
		biz.NewUserUsecase,
		biz.NewCommunityUsecase,

		// service
		service.NewGreeterService,
		service.NewAuthService,
		service.NewUserService,
		service.NewCommunityService,

		// server
		server.NewHTTPServer,
		server.NewGRPCServer,

		// app
		newApp,
	))
}
