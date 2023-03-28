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
		msg := c.InsertMessage("user " + c.Username + " telah meninggalkan ruangan")
		hub.Broadcast <- msg
		hub.Unregister <- c
	}()

	msg := c.InsertMessage("halo saya adalah bot, ada yang bisa saya bantu?")

	for { //listening to new incoming message
		message, ok := <-c.Message

		if message.Content == "0" && ok {
			msg = c.InsertMessage("selamat tinggal! -bot")
			hub.Broadcast <- msg //send the answer to the hub
			return
		}

		if message.Role != "bot" && ok { //if there is a new message that belongs to the user
			answer, err := c.DfUtil.DetectIntent(message.Content, message.RoomId) //send it to the dialogflow api //take the answer form the dialogflow
			if err != nil {
				fmt.Println(err.Error())
			}
			msg = c.InsertMessage(answer)
			hub.Broadcast <- msg //send the answer to the hub
		}
	}
}

func (c *DialogflowClient) InsertMessage(content string) *Message {
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
		Username:  c.Username, //bot
		Role:      c.Role,     //bot
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		DeletedAt: time.Time{},
	}
	return msg
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

func (c DialogflowClient) GetRole() string {
	return c.Role
}
