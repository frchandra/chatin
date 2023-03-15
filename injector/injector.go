//go:build wireinject
// +build wireinject

package injector

import (
	"github.com/frchandra/chatin/app"
	"github.com/frchandra/chatin/app/controller"
	"github.com/frchandra/chatin/app/middleware"
	"github.com/frchandra/chatin/app/repository"
	"github.com/frchandra/chatin/app/service"
	"github.com/frchandra/chatin/app/util"
	"github.com/frchandra/chatin/config"
	"github.com/frchandra/chatin/database"
	"github.com/gin-gonic/gin"
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

var UtilSet = wire.NewSet(
	util.NewTokenUtil,
	util.NewLogUtil,
)

func InitializeServer() *gin.Engine {
	wire.Build(
		config.NewAppConfig,
		app.NewDatabase,
		app.NewCache,
		app.NewLogger,
		UtilSet,
		MiddlewareSet,
		UserSet,
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
