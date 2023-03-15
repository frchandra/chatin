package repository

import (
	"context"
	"errors"
	"github.com/frchandra/chatin/app/model"
	"github.com/frchandra/chatin/app/util"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"reflect"
)

type UserRepository struct {
	db  *mongo.Database
	log *util.LogUtil
}

func NewUserRepository(db *mongo.Database, log *util.LogUtil) *UserRepository {
	return &UserRepository{db: db, log: log}
}

func (r *UserRepository) GetOne(user *model.User) (model.User, error) {
	filter := bson.M{
		"name":  user.Name,
		"email": user.Email,
	}
	result := r.db.Collection("users").FindOne(context.Background(), filter)
	var resultUser model.User
	_ = result.Decode(&resultUser)
	return resultUser, result.Err()
}

func (r *UserRepository) GetOneByNameOrEmail(user *model.User) (model.User, error) {
	filter := bson.D{
		{"$or",
			bson.A{
				bson.D{{"name", user.Name}},
				bson.D{{"email", user.Email}},
			}},
	}
	result := r.db.Collection("users").FindOne(context.Background(), filter)
	var resultUser model.User
	_ = result.Decode(&resultUser)
	return resultUser, result.Err()
}

func (r *UserRepository) GetOneById(id string) (model.User, error) {
	bsonId, _ := primitive.ObjectIDFromHex(id)
	filter := bson.M{
		"_id": bsonId,
	}
	result := r.db.Collection("users").FindOne(context.Background(), filter)
	var resultUser model.User
	_ = result.Decode(&resultUser)
	return resultUser, result.Err()
}

func (r *UserRepository) GetOrInsertOne(user *model.User) (model.User, error) {
	resultUser, _ := r.GetOne(user)
	if reflect.ValueOf(resultUser).IsZero() { //if user is not exist
		resultUser, err := r.InsertOne(user) //create user
		return resultUser, err
	}
	//if user exist
	return resultUser, nil //return the found user
}

func (r *UserRepository) InsertOne(user *model.User) (model.User, error) {
	filter := bson.D{
		{
			"$or",
			bson.A{
				bson.M{"name": user.Name},
				bson.M{"email": user.Email},
			},
		},
	}
	var resultUser model.User
	if result := r.db.Collection("users").FindOne(context.Background(), filter); result.Err() == nil { //Able to find user with this username/email. This means user with this username/email is already exist
		return resultUser, errors.New("user with this username/email is already exist")
	}
	user.Id = primitive.NewObjectID()
	_, err := r.db.Collection("users").InsertOne(context.Background(), user)
	return *user, err
}
