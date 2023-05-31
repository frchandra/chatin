package controller

import (
	"fmt"
	"github.com/frchandra/chatin/database"
	"github.com/gin-gonic/gin"
	"net/http"
)

type ConfigController struct {
	migrator *database.Migrator
}

func NewConfigController(migrator *database.Migrator) *ConfigController {
	return &ConfigController{migrator: migrator}
}

func (r ConfigController) RefreshDb(c *gin.Context) {
	fmt.Println("running migrations")
	r.migrator.RunMigration()
	c.JSON(http.StatusOK, gin.H{"message": "success"})
	return
}
