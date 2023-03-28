package messenger

type Client interface {
	Subscriber()
	Publisher(hub *Hub)
	GetId() string
	GetRoomId() string
	GetUsername() string
	GetMessage() chan *Message
	GetRole() string
}
