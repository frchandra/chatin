package model

import "go.mongodb.org/mongo-driver/bson/primitive"

type Message struct { ///
	Id       primitive.ObjectID `bson:"_id" json:"id"`
	Content  string             `bson:"content" json:"content"`
	Username string             `bson:"username" json:"username"`
	Role     string             `bson:"role" json:"role"`
}
