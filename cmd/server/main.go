package main

import (
	"github.com/frchandra/chatin/config"
	"github.com/frchandra/chatin/injector"
)

func main() {
	appConfig := config.NewAppConfig()
	server := injector.InitializeServer()
	server.Run(":" + appConfig.AppPort)
}
