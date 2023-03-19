package repository

import (
	"github.com/frchandra/chatin/app/util"
	"go.mongodb.org/mongo-driver/mongo"
)

type RoomRepository struct {
	db  *mongo.Database
	log *util.LogUtil
}

func NewRoomRepository(db *mongo.Database, log *util.LogUtil) *RoomRepository {
	return &RoomRepository{db: db, log: log}
}
