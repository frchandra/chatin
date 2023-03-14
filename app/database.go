package app

import (
	"context"
	"github.com/frchandra/chatin/config"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

func NewDatabase(appConfig *config.AppConfig, log *logrus.Logger) *mongo.Database {
	clientOptions := options.Client().
		ApplyURI("mongodb://" + appConfig.DBHost + ":" + appConfig.DBPort).
		SetAuth(options.Credential{
			Username: appConfig.DBUser,
			Password: appConfig.DBPassword,
		})
	client, err := mongo.NewClient(clientOptions)
	if err != nil {
		log.Panic("failed on creating client. Error: " + err.Error())
	}
	if err = client.Connect(context.Background()); err != nil {
		log.Panic("failed on connecting to the database server. Error: " + err.Error())
	}
	if err = client.Ping(context.Background(), readpref.Primary()); err != nil {
		log.Panic("failed on contacting to the database server. Error: " + err.Error())
	}
	log.Info("application is successfully connected to the database server")
	return client.Database(appConfig.DBName)
}
