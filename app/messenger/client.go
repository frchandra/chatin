package messenger

import (
	"fmt"
	"github.com/frchandra/chatin/app/model"
	"github.com/frchandra/chatin/app/service"
	"github.com/gorilla/websocket"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"log"
	"time"
)

type Client struct {
	Conn        *websocket.Conn
	Message     chan *Message
	Id          string `json:"id"`
	RoomId      string `json:"room_id"`
	Username    string `json:"username"`
	RoomService *service.RoomService
}

type Message struct {
	Id        string    `json:"id"`
	Content   string    `json:"content"`
	Username  string    `json:"username"`
	Role      string    `json:"role"`
	RoomId    string    `json:"roomId"`
	CreatedAt time.Time `bson:"created_at" json:"created_at"`
	UpdatedAt time.Time `bson:"updated_at" json:"updated_at"`
	DeletedAt time.Time `bson:"deleted_at" json:"deleted_at"`
}

func (c *Client) WriteMessage() {
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

func (c *Client) ReadMessage(hub *Hub) {
	defer func() {
		hub.Unregister <- c
		_ = c.Conn.Close()
	}()

	for {
		_, m, err := c.Conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error: %v", err)
			}
			break
		}

		msgId := primitive.NewObjectID()

		msg := &Message{
			Id:        msgId.Hex(),
			Content:   string(m),
			RoomId:    c.RoomId,
			Username:  c.Username,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
			DeletedAt: time.Time{},
		}

		fmt.Println("######### RECEIVE NEW MSG : ")
		fmt.Println(msg)

		result, err := c.RoomService.InsertMessage(msg.RoomId, &model.Message{
			Id:       msgId,
			Content:  msg.Content,
			Username: msg.Username,
			Role:     "user", //TODO: make this dynamic
		})
		if err != nil {
			fmt.Println("######### DB OPS ERROR")
			fmt.Println(err.Error())
		} else {
			fmt.Println("######### DB OPS SUCCESS")
			fmt.Println(result)
		}

		hub.Broadcast <- msg
	}
}
