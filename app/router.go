package app

import (
	"github.com/frchandra/chatin/app/controller"
	"github.com/gin-gonic/gin"
	"net/http"
)

func NewRouter(
	userController *controller.UserController,
) *gin.Engine {
	router := gin.Default()

	public := router.Group("/api/v1")
	public.GET("/home", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"OK": "OK",
		})
	})
	public.POST("/user/register", userController.Register)

	return router
}
