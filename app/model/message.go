package model

import (
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type Message struct { ///
	Id        primitive.ObjectID `bson:"_id" json:"id"`
	Content   string             `bson:"content" json:"content"`
	RoomId    string             `bson:"room_id" json:"room_id"`
	Username  string             `bson:"username" json:"username"`
	Role      string             `bson:"role" json:"role"`
	CreatedAt time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt time.Time          `bson:"updated_at" json:"updated_at"`
	DeletedAt time.Time          `bson:"deleted_at" json:"deleted_at"`
}

func (m *Message) MarshalBSON() ([]byte, error) {
	if m.CreatedAt.IsZero() {
		m.CreatedAt = time.Now()
	}
	m.UpdatedAt = time.Now()

	type my Message
	return bson.Marshal((*my)(m))
}
