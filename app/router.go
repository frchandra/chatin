package app

import (
	"github.com/frchandra/chatin/app/controller"
	"github.com/frchandra/chatin/app/middleware"
	"github.com/gin-gonic/gin"
)

func NewRouter(
	userMiddleware *middleware.UserMiddleware,

	userController *controller.UserController,
) *gin.Engine {
	router := gin.Default()

	//Public User Standard Auth Routes
	public := router.Group("/api/v1")
	public.POST("/user/register", userController.Register)
	public.POST("/user/sign_in", userController.SignIn)
	public.POST("/user/login", userController.Login)
	public.POST("/user/refresh", userController.RefreshToken)

	//Logged-In User Routes
	user := router.Group("/api/v1").Use(userMiddleware.HandleUserAccess)
	user.POST("/user/logout", userController.Logout)
	user.GET("/user", userController.CurrentUser)

	return router
}
