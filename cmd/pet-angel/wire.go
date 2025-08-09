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
func wireApp(*conf.Bootstrap, log.Logger) (*kratos.App, func(), error) {
	panic(wire.Build(
		// infrastructure layer
		data.NewData,
		data.NewMysql,
		data.NewRedis,

		// data layer
		data.NewUserRepo,
		data.NewPetModelRepo,
		data.NewItemRepo,
		data.NewChatMessageRepo,
		data.NewCategoryRepo,
		data.NewPostRepo,
		data.NewCommentRepo,
		data.NewUserRelationRepo,
		data.NewUserInteractionRepo,
		data.NewNoteRepo,
		data.NewUserUnlockRecordRepo,

		// biz layer
		biz.NewAuthUsecase,
		biz.NewAvatarUsecase,
		biz.NewCommunityUsecase,
		biz.NewUserUsecase,
		biz.NewMessageUsecase,

		// service layer
		service.NewAuthService,
		service.NewAvatarService,
		service.NewCommunityService,
		service.NewUserService,
		service.NewMessageService,
		service.NewUploadService,

		// server layer
		server.NewHTTPServer,
		server.NewGRPCServer,

		// app layer
		newApp,
	))
}
