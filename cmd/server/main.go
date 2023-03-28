package main

import (
	"github.com/frchandra/chatin/config"
	"github.com/frchandra/chatin/injector"
)

func main() {
	appConfig := config.NewAppConfig()
	server := injector.InitializeServer()
	go server.Hub.Run()                     //this must be run
	server.Web.Run(":" + appConfig.AppPort) //after this, if not the hub's channel won't get the data sent from the client's goroutine. ? maybe because this line of code is blocking ? => confirmed by simple experiment
}
