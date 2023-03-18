package websocket

import (
	"github.com/frchandra/chatin/app/validation"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"net/http"
)

type WsHandler struct {
	Hub *Hub
}

func NewHandler(hub *Hub) *WsHandler {
	return &WsHandler{Hub: hub}
}

func (h *WsHandler) CreateRoom(c *gin.Context) {
	var request validation.CreateRoomRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "fail", "error": err.Error()})
		return
	}

	h.Hub.Rooms[request.Id] = &Room{ //TODO: persist to DB
		Id:      request.Id,
		Name:    request.Name,
		Clients: make(map[string]*Client),
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

func (h *WsHandler) JoinRoom(c *gin.Context) {
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "fail", "error": err.Error()})
		return
	}

	roomId := c.Param("roomId")
	clientId := c.Query("userId")
	username := c.Query("username")

	client := &Client{
		Conn:     conn,
		Message:  make(chan *Message, 10),
		Id:       clientId,
		RoomId:   roomId,
		Username: username,
	}

	message := &Message{
		Content:  "A new user has joined the room",
		RoomId:   roomId,
		Username: username,
	}

	//Register a new client through the register channel
	h.Hub.Register <- client
	//Broadcast that message
	h.Hub.Broadcast <- message

	//writeMessage
	go client.writeMessage()
	//readMessage (blocking)
	client.readMessage(h.Hub)
}

func (h *WsHandler) GetRooms(c *gin.Context) {
	rooms := make([]validation.GetRoomResponse, 0)
	for _, r := range h.Hub.Rooms {
		rooms = append(rooms, validation.GetRoomResponse{
			Id:   r.Id,
			Name: r.Name,
		})
	}
	c.JSON(http.StatusOK, rooms)
	return
}

func (h *WsHandler) GetClients(c *gin.Context) {
	var clients []validation.GetClientResponse
	roomId := c.Param("roomId")
	if _, ok := h.Hub.Rooms[roomId]; !ok {
		clients = make([]validation.GetClientResponse, 0)
		c.JSON(http.StatusBadRequest, clients)
		return
	}
	for _, c := range h.Hub.Rooms[roomId].Clients {
		clients = append(clients, validation.GetClientResponse{
			Id:       c.Id,
			Username: c.Username,
		})
	}
	c.JSON(http.StatusOK, clients)
	return
}
