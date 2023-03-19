package controller

import (
	"github.com/frchandra/chatin/app/messenger"
	"github.com/frchandra/chatin/app/model"
	"github.com/frchandra/chatin/app/service"
	"github.com/frchandra/chatin/app/validation"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"net/http"
	"time"
)

type RoomController struct {
	roomSvc *service.RoomService
	Hub     *messenger.Hub
}

func NewRoomController(roomSvc *service.RoomService, hub *messenger.Hub) *RoomController {
	return &RoomController{roomSvc: roomSvc, Hub: hub}
}

func (r *RoomController) CreateRoom(c *gin.Context) {
	var request validation.CreateRoomRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "fail", "error": err.Error()})
		return
	}

	messages := []model.Message{
		model.Message{
			Id:      primitive.NewObjectID(),
			Content: "room " + request.Name + " is created",
		},
	}

	roomResult, err := r.roomSvc.InsertOne(&model.Room{Name: request.Name, Messages: messages})
	if err != nil { //persist to db
		c.JSON(http.StatusBadRequest, gin.H{"status": "fail", "error": err.Error()})
	}

	r.Hub.Rooms[roomResult.Id.Hex()] = &messenger.Room{
		Id:        roomResult.Id.Hex(),
		Name:      request.Name,
		Clients:   make(map[string]*messenger.Client),
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
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "fail", "error": err.Error()})
		return
	}

	roomId := c.Param("roomId")
	clientId := c.Query("userId")
	username := c.Query("username")

	client := &messenger.Client{
		Conn:        conn,
		Message:     make(chan *messenger.Message, 10),
		Id:          clientId,
		RoomId:      roomId,
		Username:    username,
		RoomService: r.roomSvc,
	}

	message := &messenger.Message{
		Content:  "A new user has joined the room",
		RoomId:   roomId,
		Username: username,
	}

	r.Hub.Register <- client   //Register a new client through the register channel
	r.Hub.Broadcast <- message //Broadcast that message

	go client.WriteMessage()  //writeMessage (non-blocking)
	client.ReadMessage(r.Hub) //readMessage (blocking)
}

func (r *RoomController) GetRooms(c *gin.Context) {
	rooms := make([]messenger.Room, 0)
	for _, room := range r.Hub.Rooms {
		rooms = append(rooms, *room)
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
			Id:       client.Id,
			Username: client.Username,
		})
	}
	c.JSON(http.StatusOK, clients)
	return
}
