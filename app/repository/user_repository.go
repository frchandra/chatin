package repository

import (
	"context"
	"github.com/frchandra/chatin/app/model"
	"github.com/frchandra/chatin/app/util"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type UserRepository struct {
	db  *mongo.Database
	log *util.LogUtil
}

func NewUserRepository(db *mongo.Database, log *util.LogUtil) *UserRepository {
	return &UserRepository{db: db, log: log}
}

func (r *UserRepository) GetOrInsertOne(user *model.User) model.User {
	result := r.FindOne(user)
	if result.Err() != nil {
		_, err := r.InsertOne(user)
		if err != nil {
			r.log.BasicErrorLog(err, "UserRepository@GetOrInsertOne")
		}
		return *user
	} else {
		var userResult model.User
		_ = result.Decode(&userResult)
		return userResult
	}
}

func (r *UserRepository) FindOne(user *model.User) *mongo.SingleResult {
	filter := bson.M{
		"name":  user.Name,
		"email": user.Email,
	}
	result := r.db.Collection("users").FindOne(context.Background(), filter)
	return result
}

func (r *UserRepository) InsertOne(user *model.User) (*mongo.InsertOneResult, error) {
	userBson, _ := bson.Marshal(user)
	result, err := r.db.Collection("users").InsertOne(context.Background(), userBson)
	return result, err
}
