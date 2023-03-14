package repository

import (
	"context"
	"errors"
	"github.com/frchandra/chatin/app/model"
	"github.com/frchandra/chatin/app/util"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type UserRepository struct {
	db  *mongo.Database
	log *util.LogUtil
}

func NewUserRepository(db *mongo.Database, log *util.LogUtil) *UserRepository {
	return &UserRepository{db: db, log: log}
}

func (r *UserRepository) GetOne(user *model.User) *mongo.SingleResult {
	filter := bson.M{
		"name":  user.Name,
		"email": user.Email,
	}
	result := r.db.Collection("users").FindOne(context.Background(), filter)
	return result
}

func (r *UserRepository) GetOneWithPassword(user *model.User) *mongo.SingleResult {
	filter := bson.M{
		"name":     user.Name,
		"email":    user.Email,
		"password": user.Password,
	}
	result := r.db.Collection("users").FindOne(context.Background(), filter)
	return result
}

func (r *UserRepository) GetOneById(id string) *mongo.SingleResult {
	filter := bson.M{
		"_id": id,
	}
	result := r.db.Collection("users").FindOne(context.Background(), filter)
	return result
}

func (r *UserRepository) GetOrInsertOne(user *model.User) (model.User, error) {
	result := r.GetOne(user)
	if result.Err() != nil { //if user is not exist
		if _, err := r.InsertOne(user); err != nil { //create user
			r.log.BasicErrorLog(err, "UserRepository@GetOrInsertOne")
			return *user, err
		}
		return *user, nil
	} else { //if user exist
		var userResult model.User
		_ = result.Decode(&userResult)
		return userResult, nil //return the found user
	}
}

func (r *UserRepository) InsertOne(user *model.User) (*mongo.InsertOneResult, error) {
	filter := bson.D{
		{
			"$or",
			bson.A{
				bson.M{"name": user.Name},
				bson.M{"email": user.Email},
			},
		},
	}
	if result := r.db.Collection("users").FindOne(context.Background(), filter); result.Err() == nil { //Able to find user with this username/email. This means user with this username/email is already exist
		return nil, errors.New("user with this username/email is already exist")
	}
	user.Id = primitive.NewObjectID()
	result, err := r.db.Collection("users").InsertOne(context.Background(), user)
	return result, err
}
