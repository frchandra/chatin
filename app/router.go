package app

import (
	"github.com/frchandra/chatin/app/controller"
	"github.com/frchandra/chatin/app/middleware"
	"github.com/frchandra/chatin/app/websocket"
	"github.com/gin-gonic/gin"
)

type Server struct {
	Web *gin.Engine
	Hub *websocket.Hub
}

func NewRouter(userMiddleware *middleware.UserMiddleware, userController *controller.UserController, wsHandler *websocket.WsHandler) *Server {
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

	//Websocket Routes
	ws := router.Group("/ws/v1")
	ws.POST("/room", wsHandler.CreateRoom)
	ws.GET("/room/join/:roomId", wsHandler.JoinRoom)
	ws.GET("/room", wsHandler.GetRooms)
	ws.GET("/clients/:roomId", wsHandler.GetClients)

	server := &Server{
		Web: router,
		Hub: wsHandler.Hub,
	}

	return server
}
