//go:build wireinject
// +build wireinject

// The build tag makes sure the stub is not built in the final build.

package main

import (
	"pet-angel/internal/biz"
	"pet-angel/internal/conf"
	"pet-angel/internal/data"
	"pet-angel/internal/server"
	"pet-angel/internal/service"

	"github.com/go-kratos/kratos/v2"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/google/wire"
)

// wireApp init kratos application.
func wireApp(*conf.Server, *conf.Data, log.Logger) (*kratos.App, func(), error) {
	panic(wire.Build(
		// server
		server.NewHTTPServer,
		server.NewGRPCServer,
		newApp,

		// data layer
		data.NewData,
		data.NewGreeterRepo,

		// biz layer
		biz.NewGreeterUsecase,

		// service layer
		service.NewGreeterService,
	))
}
