package model

import "go.mongodb.org/mongo-driver/bson/primitive"

type Room struct { ///
	Id       primitive.ObjectID `bson:"_id" json:"id"`
	Name     string             `bson:"name" json:"name"`
	Messages []Message          `bson:"messages" json:"messages"`
}
