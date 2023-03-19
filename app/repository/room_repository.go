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

func (r *RoomRepository) GetOneByName(room *model.Room) (model.Room, error) {
	filter := bson.M{
		"name": room.Name,
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
