package controller

import (
	"fmt"
	"github.com/frchandra/chatin/config"
	"github.com/frchandra/chatin/database"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

type ConfigController struct {
	migrator *database.Migrator
	config   *config.AppConfig
}

func NewConfigController(migrator *database.Migrator, config *config.AppConfig) *ConfigController {
	return &ConfigController{migrator: migrator, config: config}
}

func (r *ConfigController) HealthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message":   "success",
		"app_name":  r.config.AppName,
		"app_url":   r.config.AppUrl,
		"time_unix": time.Now().Unix(),
		"time":      time.Now(),
	})
	return
}

func (r *ConfigController) RefreshDb(c *gin.Context) {
	fmt.Println("running migrations")
	r.migrator.RunMigration()
	c.JSON(http.StatusOK, gin.H{"message": "success"})
	return
}
