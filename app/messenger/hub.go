package messenger

import "time"

type Room struct { ///
	Id        string            `json:"id"`
	Name      string            `json:"name"`
	Clients   map[string]Client `json:"-"`
	CreatedAt time.Time         `json:"created_at"`
	UpdatedAt time.Time         `json:"updated_at"`
	DeletedAt time.Time         `json:"deleted_at"`
}

type Hub struct {
	Rooms      map[string]*Room
	Register   chan Client
	Unregister chan Client
	Broadcast  chan *Message
}

func NewHub() *Hub {
	return &Hub{
		Rooms:      make(map[string]*Room),
		Register:   make(chan Client),
		Unregister: make(chan Client),
		Broadcast:  make(chan *Message, 5),
	}
}

func (h *Hub) Run() {
	for {
		select {
		case cl := <-h.Register: //new client has been added
			if _, ok := h.Rooms[cl.GetRoomId()]; ok { //check: is the client's roomId is actually exist?
				r := h.Rooms[cl.GetRoomId()]
				if _, ok = r.Clients[cl.GetId()]; !ok { //check: is the client hasn't been added to the room?
					r.Clients[cl.GetId()] = cl //add the client to this room
				}
			}
		case cl := <-h.Unregister: //client left the room
			if _, ok := h.Rooms[cl.GetRoomId()]; ok { //check: is the client's roomId is actually exist?
				if _, ok = h.Rooms[cl.GetRoomId()].Clients[cl.GetId()]; ok { //check: is the client is actually exist in this room?
					delete(h.Rooms[cl.GetRoomId()].Clients, cl.GetId()) //delete the client
					close(cl.GetMessage())                              //close the client's message channel
				}
			}
		case m := <-h.Broadcast: //the client send a new message
			if _, ok := h.Rooms[m.RoomId]; ok { //check: is the client's roomId is actually exist?
				for _, cl := range h.Rooms[m.RoomId].Clients { //broadcast the new messages to every client in the room
					cl.GetMessage() <- m
				}
			}
		}
	}
}
