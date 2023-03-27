package messenger

import (
	"fmt"
	"github.com/frchandra/chatin/app/model"
	"github.com/frchandra/chatin/app/service"
	"github.com/frchandra/chatin/app/util"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type DialogflowClient struct {
	Message     chan *Message
	Id          string `json:"id"`
	RoomId      string `json:"room_id"`
	Username    string `json:"username"`
	Role        string `json:"role"`
	DfUtil      *util.DialogflowUtil
	RoomService *service.RoomService
}

func (c *DialogflowClient) Subscriber() {
	defer func() {

	}()

}

func (c *DialogflowClient) Publisher(hub *Hub) { //this should be non-blocking process
	defer func() {
		hub.Unregister <- c
	}()

	msgId := primitive.NewObjectID() //send hello message
	msg := &Message{                 //create messenger payload
		Id:        msgId.Hex(),
		Content:   "hello i am a bot, how can i help you",
		RoomId:    c.RoomId,
		Username:  c.Username, //bot
		Role:      c.Role,     //bot
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		DeletedAt: time.Time{},
	}
	_, _ = c.RoomService.InsertMessage(msg.RoomId, &model.Message{ //insert payload to database
		Id:       msgId,
		Content:  msg.Content,
		Username: msg.Username,
		Role:     msg.Role,
		RoomId:   msg.RoomId,
	})
	hub.Broadcast <- msg //send the answer to the hub

	for { //listening to new incoming message
		message, ok := <-c.Message
		if message.Content == "0" && ok {
			msgId = primitive.NewObjectID()
			msg = &Message{ //create messenger payload
				Id:        msgId.Hex(),
				Content:   "good bye! -bot",
				RoomId:    c.RoomId,
				Username:  c.Username, //bot
				Role:      c.Role,     //bot
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
				DeletedAt: time.Time{},
			}
			_, _ = c.RoomService.InsertMessage(msg.RoomId, &model.Message{ //insert payload to database
				Id:       msgId,
				Content:  msg.Content,
				Username: msg.Username,
				Role:     msg.Role,
				RoomId:   msg.RoomId,
			})
			hub.Broadcast <- msg //send the answer to the hub
			hub.Unregister <- c
			return
		}

		if message.Role != "bot" && ok { //if there is a new message that belongs to the user
			answer, err := c.DfUtil.DetectIntent(message.Content) //send it to the dialogflow api //take the answer form the dialogflow
			if err != nil {
				fmt.Println(err.Error())
			}
			msgId = primitive.NewObjectID()
			msg = &Message{ //create messenger payload
				Id:        msgId.Hex(),
				Content:   answer,
				RoomId:    c.RoomId,
				Username:  c.Username, //bot
				Role:      c.Role,     //bot
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
				DeletedAt: time.Time{},
			}
			_, _ = c.RoomService.InsertMessage(msg.RoomId, &model.Message{ //insert payload to database
				Id:       msgId,
				Content:  msg.Content,
				Username: msg.Username,
				Role:     msg.Role,
				RoomId:   msg.RoomId,
			})
			hub.Broadcast <- msg //send the answer to the hub
		}

	}
}

func (c *DialogflowClient) GetId() string {
	return c.Id
}

func (c *DialogflowClient) GetRoomId() string {
	return c.RoomId
}

func (c *DialogflowClient) GetUsername() string {
	return c.Username
}

func (c *DialogflowClient) GetMessage() chan *Message {
	return c.Message
}
