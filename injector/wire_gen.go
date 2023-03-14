// Code generated by Wire. DO NOT EDIT.

//go:generate go run github.com/google/wire/cmd/wire
//go:build !wireinject
// +build !wireinject

package injector

import (
	"github.com/frchandra/chatin/app"
	"github.com/frchandra/chatin/app/controller"
	"github.com/frchandra/chatin/app/repository"
	"github.com/frchandra/chatin/app/service"
	"github.com/frchandra/chatin/app/util"
	"github.com/frchandra/chatin/config"
	"github.com/frchandra/chatin/database"
	"github.com/gin-gonic/gin"
	"github.com/google/wire"
)

// Injectors from injector.go:

func InitializeServer() *gin.Engine {
	appConfig := config.NewAppConfig()
	logger := app.NewLogger(appConfig)
	database := app.NewDatabase(appConfig, logger)
	logUtil := util.NewLogUtil(logger)
	userRepository := repository.NewUserRepository(database, logUtil)
	userService := service.NewUserService(userRepository)
	userController := controller.NewUserController(userService, logUtil)
	engine := app.NewRouter(userController)
	return engine
}

func InitializeMigrator() *database.Migrator {
	appConfig := config.NewAppConfig()
	logger := app.NewLogger(appConfig)
	mongoDatabase := app.NewDatabase(appConfig, logger)
	migrator := database.NewMigrator(mongoDatabase, logger)
	return migrator
}

// injector.go:

var UserSet = wire.NewSet(repository.NewUserRepository, service.NewUserService, controller.NewUserController)

var UtilSet = wire.NewSet(util.NewLogUtil)
