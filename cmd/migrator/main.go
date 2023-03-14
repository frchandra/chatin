package main

import (
	"fmt"
	"github.com/frchandra/chatin/injector"
)

func main() {
	fmt.Println("running migrations")
	migrator := injector.InitializeMigrator()
	migrator.RunMigration()
}
