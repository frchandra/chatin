package model

import (
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type Room struct { ///
	Id        primitive.ObjectID `bson:"_id" json:"id"`
	Name      string             `bson:"name" json:"name"`
	Messages  []Message          `bson:"messages" json:"messages"`
	CreatedAt time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt time.Time          `bson:"updated_at" json:"updated_at"`
	DeletedAt time.Time          `bson:"deleted_at" json:"deleted_at"`
}

func (r *Room) MarshalBSON() ([]byte, error) {
	if r.CreatedAt.IsZero() {
		r.CreatedAt = time.Now()
	}
	r.UpdatedAt = time.Now()
	type my Room
	return bson.Marshal((*my)(r))
}
