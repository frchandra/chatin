package messenger

import (
	"github.com/gorilla/websocket"
	"log"
)

type Client struct {
	Conn     *websocket.Conn
	Message  chan *Message
	Id       string `json:"id"`
	RoomId   string `json:"room_id"`
	Username string `json:"username"`
}

type Message struct { ///
	Content  string `json:"content"`
	RoomId   string `json:"roomId"`
	Username string `json:"username"`
}

func (c *Client) WriteMessage() {
	defer func() {
		_ = c.Conn.Close()
	}()

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

		msg := &Message{
			Content:  string(m),
			RoomId:   c.RoomId,
			Username: c.Username,
		}

		hub.Broadcast <- msg
	}
}
