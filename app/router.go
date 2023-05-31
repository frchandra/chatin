package app

import (
	"github.com/frchandra/chatin/app/controller"
	"github.com/frchandra/chatin/app/messenger"
	"github.com/frchandra/chatin/app/middleware"
	"github.com/gin-gonic/gin"
)

type Server struct {
	Web *gin.Engine
	Hub *messenger.Hub
}

func NewRouter(
	configController *controller.ConfigController,
	userMiddleware *middleware.UserMiddleware,
	userController *controller.UserController,
	roomController *controller.RoomController,
) *Server {
	router := gin.Default()

	//Health Check (scope: public)
	router.GET("/api/v1/health", configController.HealthCheck)

	//Public User Standard Auth Routes
	public := router
	public.POST("/api/v1/user/register", userController.Register)
	public.POST("/api/v1/user/sign_in", userController.SignIn)
	public.POST("/api/v1/user/login", userController.Login)
	public.POST("/api/v1/user/refresh", userController.RefreshToken)
	public.DELETE("/api/v1/migrate_fresh", configController.RefreshDb)
	public.DELETE("/api/v1/clear_memory", roomController.DeleteAllRoom)

	//Logged-In User Routes
	user := router.Use(userMiddleware.HandleUserAccess)
	user.POST("/api/v1/user/logout", userController.Logout)
	user.GET("/api/v1/user", userController.CurrentUser)
	user.POST("/api/v1/room", roomController.CreateRoom)
	user.GET("/ws/v1/room/join/:roomId", roomController.JoinRoom)

	//Logged-In Admin Routes
	admin := router.Use(userMiddleware.HandleAdminAccess)
	admin.GET("/api/v1/room", roomController.GetRooms)
	admin.GET("/api/v1/clients/:roomId", roomController.GetClients)
	admin.GET("api/v1/bot_less_room", roomController.GetBotLessRooms)

	server := &Server{
		Web: router,
		Hub: roomController.Hub,
	}

	return server
}
