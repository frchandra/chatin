package controller

import (
	"github.com/frchandra/chatin/app/messenger"
	"github.com/frchandra/chatin/app/model"
	"github.com/frchandra/chatin/app/service"
	"github.com/frchandra/chatin/app/util"
	"github.com/frchandra/chatin/app/validation"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"net/http"
	"time"
)

type RoomController struct {
	roomService *service.RoomService
	userService *service.UserService
	dfUtil      *util.DialogflowUtil
	Hub         *messenger.Hub
}

func NewRoomController(roomService *service.RoomService, userService *service.UserService, dfUtil *util.DialogflowUtil, hub *messenger.Hub) *RoomController {
	return &RoomController{roomService: roomService, userService: userService, dfUtil: dfUtil, Hub: hub}
}

func (r *RoomController) CreateRoom(c *gin.Context) {
	contextData, _ := c.Get("accessDetails")              //from the context passed by the user middleware, get the details about the current user that make request from the context passed by user middleware
	accessDetails, _ := contextData.(*util.AccessDetails) //type assertion
	user, _ := r.userService.GetOneById(accessDetails.UserId)

	messages := []model.Message{ //create message payload
		model.Message{
			Id:       primitive.NewObjectID(),
			Content:  "room " + user.Username + "_room is created",
			Username: user.Username,
			Role:     user.Role,
		},
	}

	roomResult, err := r.roomService.InsertOne(&model.Room{Name: user.Username + "_room", Messages: messages}) //persist to db
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "fail", "error": err.Error()})
	}

	r.Hub.Rooms[roomResult.Id.Hex()] = &messenger.Room{ //create messenger room
		Id:        roomResult.Id.Hex(),
		Name:      user.Username + "_room",
		Clients:   make(map[string]messenger.Client),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		DeletedAt: time.Time{},
	}
	c.JSON(http.StatusOK, r.Hub.Rooms[roomResult.Id.Hex()])
	return
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func (r *RoomController) JoinRoom(c *gin.Context) {
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil) //switch http to ws protocol
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "fail", "error": err.Error()})
		return
	}

	contextData, _ := c.Get("accessDetails")              //from the context passed by the user middleware, get the details about the current user that make request from the context passed by user middleware
	accessDetails, _ := contextData.(*util.AccessDetails) //type assertion
	user, _ := r.userService.GetOneById(accessDetails.UserId)

	userClient := &messenger.UserClient{ //create messenger userClient
		Conn:        conn,
		Message:     make(chan *messenger.Message, 10),
		Id:          user.Id.Hex(),
		RoomId:      c.Param("roomId"),
		Username:    user.Username,
		Role:        user.Role,
		RoomService: r.roomService,
	}

	if user.Role != "admin" {
		botClient := &messenger.DialogflowClient{ //create messenger bot client
			Message:     make(chan *messenger.Message, 10),
			Id:          "bot_" + user.Id.Hex(),
			RoomId:      c.Param("roomId"),
			Username:    "bot",
			Role:        "bot",
			DfUtil:      r.dfUtil,
			RoomService: r.roomService,
		}
		r.Hub.Register <- botClient
		go botClient.Publisher(r.Hub)
	}

	r.Hub.Register <- userClient //Register a new userClient through the register channel
	go userClient.Subscriber()   //writeMessage (non-blocking)
	userClient.Publisher(r.Hub)  //readMessage (blocking)
}

func (r *RoomController) GetRooms(c *gin.Context) {
	rooms := make([]messenger.Room, 0)
	for _, room := range r.Hub.Rooms {
		rooms = append(rooms, *room)
	}
	c.JSON(http.StatusOK, rooms)
	return
}

func (r *RoomController) GetBotLessRooms(c *gin.Context) {
	rooms := make([]messenger.Room, 0)
	for _, room := range r.Hub.Rooms {
		isBotExist := false
		for _, client := range r.Hub.Rooms[room.Id].Clients {
			if client.GetRole() == "bot" {
				isBotExist = true
			}
		}
		if isBotExist == false {
			rooms = append(rooms, *room)
		}
	}
	c.JSON(http.StatusOK, rooms)
	return
}

func (r *RoomController) GetClients(c *gin.Context) {
	var clients []validation.GetClientResponse
	roomId := c.Param("roomId")
	if _, ok := r.Hub.Rooms[roomId]; !ok {
		clients = make([]validation.GetClientResponse, 0)
		c.JSON(http.StatusBadRequest, clients)
		return
	}
	for _, client := range r.Hub.Rooms[roomId].Clients {
		clients = append(clients, validation.GetClientResponse{
			Id:       client.GetId(),
			Username: client.GetUsername(),
			Role:     client.GetRole(),
		})
	}
	c.JSON(http.StatusOK, clients)
	return
}
