package messenger

import "time"

type Room struct { ///
	Id        string             `json:"id"`
	Name      string             `json:"name"`
	Clients   map[string]*Client `json:"clients"`
	CreatedAt time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt time.Time          `bson:"updated_at" json:"updated_at"`
	DeletedAt time.Time          `bson:"deleted_at" json:"deleted_at"`
}

type Hub struct {
	Rooms      map[string]*Room
	Register   chan *Client
	Unregister chan *Client
	Broadcast  chan *Message
}

func NewHub() *Hub {
	return &Hub{
		Rooms:      make(map[string]*Room),
		Register:   make(chan *Client),
		Unregister: make(chan *Client),
		Broadcast:  make(chan *Message, 5),
	}
}

func (h *Hub) Run() {
	for {
		select {
		case cl := <-h.Register: //new client has been added
			if _, ok := h.Rooms[cl.RoomId]; ok { //check: is the client's roomId is actually exist?
				r := h.Rooms[cl.RoomId]
				if _, ok := r.Clients[cl.Id]; !ok { //check: is the client hasn't been added to the room?
					r.Clients[cl.Id] = cl //add the client to this room
				}
			}
		case cl := <-h.Unregister: //client left the room
			if _, ok := h.Rooms[cl.RoomId]; ok { //check: is the client's roomId is actually exist?
				if _, ok := h.Rooms[cl.RoomId].Clients[cl.Id]; ok { //check: is the client is actually exist in this room?
					if len(h.Rooms[cl.RoomId].Clients) != 0 { //check: is there still another client in the room?
						h.Broadcast <- &Message{ //broadcast the default "user left the room" messages to the remaining client
							Content:  "user left the chat",
							RoomId:   cl.RoomId,
							Username: cl.Username,
						}
					}
					delete(h.Rooms[cl.RoomId].Clients, cl.Id) //delete the client
					close(cl.Message)                         //close the client's message channel
				}
			}
		case m := <-h.Broadcast: //the client send a new message
			if _, ok := h.Rooms[m.RoomId]; ok { //check: is the client's roomId is actually exist?
				for _, cl := range h.Rooms[m.RoomId].Clients { //broadcast the new messages to every client in the room
					cl.Message <- m
				}
			}
		}
	}
}