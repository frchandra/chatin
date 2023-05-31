package factory

import (
	"context"
	"github.com/frchandra/chatin/app/model"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

type UserFactory struct {
	db  *mongo.Database
	log *logrus.Logger
}

func NewUserFactory(db *mongo.Database, log *logrus.Logger) *UserFactory {
	return &UserFactory{db: db, log: log}
}

func (f *UserFactory) RunFactory() error {
	defaultPass, _ := bcrypt.GenerateFromPassword([]byte("password"), bcrypt.DefaultCost)

	users := []any{
		&model.User{Id: primitive.NewObjectID(), Username: "eko", Email: "eko@mail.com", Password: string(defaultPass), Role: "admin"},
		&model.User{Id: primitive.NewObjectID(), Username: "dagus", Email: "dagus@mail.com", Password: string(defaultPass), Role: "user"},
		&model.User{Id: primitive.NewObjectID(), Username: "bekti", Email: "bekti@mail.com", Password: string(defaultPass), Role: "user"},
		&model.User{Id: primitive.NewObjectID(), Username: "juni", Email: "juni@mail.com", Password: string(defaultPass), Role: "user"},
	}
	if _, err := f.db.Collection("users").InsertMany(context.Background(), users); err != nil {
		f.log.Error("cannot seeding database. Error: " + err.Error())
		return err
	}
	return nil
}
