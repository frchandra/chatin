package factory

import (
	"context"
	"github.com/frchandra/chatin/app/model"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/mongo"
)

type UserFactory struct {
	db  *mongo.Database
	log *logrus.Logger
}

func NewUserFactory(db *mongo.Database, log *logrus.Logger) *UserFactory {
	return &UserFactory{db: db, log: log}
}

func (f *UserFactory) RunFactory() error {
	users := []any{
		model.User{Name: "eko", Email: "eko@mail.com", Password: "password"},
		model.User{Name: "dagus", Email: "dagus@mail.com", Password: "password"},
		model.User{Name: "bekti", Email: "bekti@mail.com", Password: "password"},
		model.User{Name: "juni", Email: "juni@mail.com", Password: "password"},
	}
	if _, err := f.db.Collection("users").InsertMany(context.Background(), users); err != nil {
		f.log.Error("cannot seeding database. Error: " + err.Error())
		return err
	}
	return nil
}
