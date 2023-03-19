package repository

import (
	"context"
	"github.com/frchandra/chatin/app/model"
	"github.com/frchandra/chatin/app/util"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type RoomRepository struct {
	db  *mongo.Database
	log *util.LogUtil
}

func NewRoomRepository(db *mongo.Database, log *util.LogUtil) *RoomRepository {
	return &RoomRepository{db: db, log: log}
}

func (r *RoomRepository) GetOneById(roomId string) (model.Room, error) {
	roomIdBson, _ := primitive.ObjectIDFromHex(roomId)
	filter := bson.M{
		"_id": roomIdBson,
	}
	result := r.db.Collection("rooms").FindOne(context.Background(), filter)
	var resultRoom model.Room
	_ = result.Decode(&resultRoom)
	return resultRoom, result.Err()
}

func (r *RoomRepository) InsertOne(room *model.Room) (model.Room, error) {
	room.Id = primitive.NewObjectID()
	_, err := r.db.Collection("rooms").InsertOne(context.Background(), room)
	return *room, err
}

func (r *RoomRepository) InsertMessage(roomId string, message *model.Message) (model.Message, error) {
	roomIdBson, _ := primitive.ObjectIDFromHex(roomId)
	filter := bson.M{"_id": roomIdBson}
	newMessage := model.Message{
		Id:       primitive.NewObjectID(),
		Content:  message.Content,
		Username: message.Username,
		Role:     message.Role,
	}
	update := bson.D{{"$push", bson.M{"messages": newMessage}}}
	_, err := r.db.Collection("rooms").UpdateOne(context.Background(), filter, update)
	if err != nil {
		return model.Message{}, err
	}
	return *message, nil
}
