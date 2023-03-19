package controller

import (
	"github.com/frchandra/chatin/app/messenger"
	"github.com/frchandra/chatin/app/validation"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"net/http"
)

type RoomController struct {
	Hub *messenger.Hub
}

func NewRoomController(hub *messenger.Hub) *RoomController {
	return &RoomController{Hub: hub}
}

func (r *RoomController) CreateRoom(c *gin.Context) {
	var request validation.CreateRoomRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "fail", "error": err.Error()})
		return
	}

	r.Hub.Rooms[request.Id] = &messenger.Room{ //TODO: persist to DB
		Id:      request.Id,
		Name:    request.Name,
		Clients: make(map[string]*messenger.Client),
	}
	c.JSON(http.StatusOK, request)
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
		Conn:     conn,
		Message:  make(chan *messenger.Message, 10),
		Id:       clientId,
		RoomId:   roomId,
		Username: username,
	}

	message := &messenger.Message{
		Content:  "A new user has joined the room",
		RoomId:   roomId,
		Username: username,
	}

	//Register a new client through the register channel
	r.Hub.Register <- client
	//Broadcast that message
	r.Hub.Broadcast <- message

	//writeMessage
	go client.WriteMessage()
	//readMessage (blocking)
	client.ReadMessage(r.Hub)
}

func (r *RoomController) GetRooms(c *gin.Context) {
	rooms := make([]validation.GetRoomResponse, 0)
	for _, r := range r.Hub.Rooms {
		rooms = append(rooms, validation.GetRoomResponse{
			Id:   r.Id,
			Name: r.Name,
		})
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
	for _, c := range r.Hub.Rooms[roomId].Clients {
		clients = append(clients, validation.GetClientResponse{
			Id:       c.Id,
			Username: c.Username,
		})
	}
	c.JSON(http.StatusOK, clients)
	return
}
