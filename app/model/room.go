package model

type Room struct { ///
	Id       string    `bson:"id" json:"id"`
	Name     string    `bson:"name" json:"name"`
	Messages []Message `bson:"messages" json:"messages"`
}
