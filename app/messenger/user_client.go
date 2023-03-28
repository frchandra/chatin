package messenger

import (
	"github.com/frchandra/chatin/app/model"
	"github.com/frchandra/chatin/app/service"
	"github.com/gorilla/websocket"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"log"
	"time"
)

type UserClient struct {
	Conn        *websocket.Conn
	Message     chan *Message
	Id          string `json:"id"`
	RoomId      string `json:"room_id"`
	Username    string `json:"username"`
	Role        string `json:"role"`
	RoomService *service.RoomService
}

func (c *UserClient) Subscriber() {
	defer func() {
		_ = c.Conn.Close()
	}()

	room, _ := c.RoomService.GetOneById(c.RoomId) //get previous conversation/messages from database
	messages := room.Messages

	for _, message := range messages { //sent previous conversation/messages to this client only
		_ = c.Conn.WriteJSON(message)
	}

	for {
		message, ok := <-c.Message
		if !ok {
			return
		}
		_ = c.Conn.WriteJSON(message)
	}
}

func (c *UserClient) Publisher(hub *Hub) {
	defer func() {
		msg := c.InsertMessage("user " + c.Username + " telah meninggalkan ruangan")
		hub.Broadcast <- msg
		hub.Unregister <- c
		_ = c.Conn.Close()
	}()

	msg := c.InsertMessage("user " + c.Username + " telah bergabung kedalam ruangan")
	//hub.Broadcast <- msg

	for {
		_, m, err := c.Conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error: %v", err)
			}
			break
		}
		msg = c.InsertMessage(string(m))
		hub.Broadcast <- msg
	}
}

func (c *UserClient) InsertMessage(content string) *Message {
	result, _ := c.RoomService.InsertMessage(c.RoomId, &model.Message{ //insert payload to database
		Id:       primitive.NewObjectID(),
		Content:  content,
		Username: c.Username,
		Role:     c.Role,
		RoomId:   c.RoomId,
	})
	msg := &Message{ //create messenger payload
		Id:        result.Id.Hex(),
		Content:   content,
		RoomId:    c.RoomId,
		Username:  c.Username,
		Role:      c.Role,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		DeletedAt: time.Time{},
	}
	return msg
}

func (c *UserClient) GetId() string {
	return c.Id
}

func (c *UserClient) GetRoomId() string {
	return c.RoomId
}

func (c *UserClient) GetUsername() string {
	return c.Username
}

func (c *UserClient) GetMessage() chan *Message {
	return c.Message
}
