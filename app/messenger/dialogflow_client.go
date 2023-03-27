package messenger

import (
	"fmt"
	"github.com/frchandra/chatin/app/util"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type DialogflowClient struct {
	Message  chan *Message
	Id       string `json:"id"`
	RoomId   string `json:"room_id"`
	Username string `json:"username"`
	Role     string `json:"role"`
	DfUtil   *util.DialogflowUtil
}

func (c *DialogflowClient) Subscriber() {
	defer func() {

	}()

}

func (c *DialogflowClient) Publisher(hub *Hub) { //this should be non-blocking process
	defer func() {
		hub.Unregister <- c
	}()
	for { //listening to new incoming message
		message, ok := <-c.Message
		fmt.Println("DF MESSAGE <<<<<<<<<")
		fmt.Println(message.Content)
		fmt.Println("DF MESSAGE >>>>>>>>>")
		fmt.Println("")

		if message.Role != "bot" && ok { //if there is a new message that belongs to the user
			answer, err := c.DfUtil.DetectIntent(message.Content) //send it to the dialogflow api //take the answer form the dialogflow
			if err != nil {
				fmt.Println("DF ERROR <<<<<<<<<")
				fmt.Println(err.Error())
				fmt.Println("DF ERROR >>>>>>>>>")
			}
			msg := &Message{ //create messenger payload
				Id:        primitive.NewObjectID().Hex(),
				Content:   answer,
				RoomId:    c.RoomId,
				Username:  c.Username, //bot
				Role:      c.Role,     //bot
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
				DeletedAt: time.Time{},
			}
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
