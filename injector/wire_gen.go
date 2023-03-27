// Code generated by Wire. DO NOT EDIT.

//go:generate go run github.com/google/wire/cmd/wire
//go:build !wireinject
// +build !wireinject

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

// Injectors from injector.go:

func InitializeServer() *app.Server {
	appConfig := config.NewAppConfig()
	client := app.NewCache(appConfig)
	tokenUtil := util.NewTokenUtil(client, appConfig)
	logger := app.NewLogger(appConfig)
	database := app.NewDatabase(appConfig, logger)
	logUtil := util.NewLogUtil(logger)
	userRepository := repository.NewUserRepository(database, logUtil)
	userService := service.NewUserService(userRepository, tokenUtil)
	userMiddleware := middleware.NewUserMiddleware(tokenUtil, logger, userService)
	userController := controller.NewUserController(userService, tokenUtil, appConfig, logUtil)
	roomRepository := repository.NewRoomRepository(database, logUtil)
	roomService := service.NewRoomService(roomRepository)
	sessionsClient := app.NewChatBot(appConfig, logger)
	dialogflowUtil := util.NewDialogflowUtil(sessionsClient, appConfig)
	hub := messenger.NewHub()
	roomController := controller.NewRoomController(roomService, userService, dialogflowUtil, hub)
	server := app.NewRouter(userMiddleware, userController, roomController, dialogflowUtil)
	return server
}

func InitializeMigrator() *database.Migrator {
	appConfig := config.NewAppConfig()
	logger := app.NewLogger(appConfig)
	mongoDatabase := app.NewDatabase(appConfig, logger)
	migrator := database.NewMigrator(mongoDatabase, logger)
	return migrator
}

// injector.go:

var MiddlewareSet = wire.NewSet(middleware.NewUserMiddleware)

var UserSet = wire.NewSet(repository.NewUserRepository, service.NewUserService, controller.NewUserController)

var RoomSet = wire.NewSet(repository.NewRoomRepository, service.NewRoomService, controller.NewRoomController)

var UtilSet = wire.NewSet(util.NewTokenUtil, util.NewLogUtil, util.NewDialogflowUtil)

var MessengerSet = wire.NewSet(messenger.NewHub)
