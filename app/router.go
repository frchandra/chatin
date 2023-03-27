package app

import (
	"github.com/frchandra/chatin/app/controller"
	"github.com/frchandra/chatin/app/messenger"
	"github.com/frchandra/chatin/app/middleware"
	"github.com/frchandra/chatin/app/util"
	"github.com/gin-gonic/gin"
	"net/http"
)

type Server struct {
	Web *gin.Engine
	Hub *messenger.Hub
}

func NewRouter(
	userMiddleware *middleware.UserMiddleware,
	userController *controller.UserController,
	roomController *controller.RoomController,

	dfUtil *util.DialogflowUtil,
) *Server {
	router := gin.Default()

	//Public User Standard Auth Routes
	public := router
	public.POST("/api/v1/user/register", userController.Register)
	public.POST("/api/v1/user/sign_in", userController.SignIn)
	public.POST("/api/v1/user/login", userController.Login)
	public.POST("/api/v1/user/refresh", userController.RefreshToken)

	public.POST("/invoke", func(c *gin.Context) {
		response, err := dfUtil.DetectIntent("hello")
		if err != nil {
			c.JSON(http.StatusOK, gin.H{"status": "error", "message": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"status": "OK", "message": response})
		return
	})

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

	server := &Server{
		Web: router,
		Hub: roomController.Hub,
	}

	return server
}
