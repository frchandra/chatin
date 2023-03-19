//go:build wireinject
// +build wireinject

package injector

import (
	"github.com/frchandra/chatin/app"
	"github.com/frchandra/chatin/app/controller"
	"github.com/frchandra/chatin/app/messenger"
	"github.com/frchandra/chatin/app/middleware"
	"github.com/frchandra/chatin/app/repository"
	"github.com/frchandra/chatin/app/service"
	"github.com/frchandra/chatin/app/util"
	"github.com/frchandra/chatin/config"
	"github.com/frchandra/chatin/database"
	"github.com/google/wire"
)

var MiddlewareSet = wire.NewSet(
	middleware.NewUserMiddleware,
)

var UserSet = wire.NewSet(
	repository.NewUserRepository,
	service.NewUserService,
	controller.NewUserController,
)

var RoomSet = wire.NewSet(
	controller.NewRoomController,
)

var UtilSet = wire.NewSet(
	util.NewTokenUtil,
	util.NewLogUtil,
)

var MessengerSet = wire.NewSet(
	messenger.NewHub,
)

func InitializeServer() *app.Server {
	wire.Build(
		config.NewAppConfig,
		app.NewDatabase,
		app.NewCache,
		app.NewLogger,
		MessengerSet,
		UtilSet,
		MiddlewareSet,
		UserSet,
		RoomSet,
		app.NewRouter,
	)
	return nil
}

func InitializeMigrator() *database.Migrator {
	wire.Build(
		config.NewAppConfig,
		app.NewLogger,
		app.NewDatabase,
		database.NewMigrator,
	)
	return nil
}
