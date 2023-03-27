package messenger

import "time"

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
