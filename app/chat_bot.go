package app

import (
	dialogflow "cloud.google.com/go/dialogflow/apiv2"
	"context"
	"github.com/frchandra/chatin/config"
	"github.com/sirupsen/logrus"
	"google.golang.org/api/option"
)

func NewChatBot(appConfig *config.AppConfig, log *logrus.Logger) *dialogflow.SessionsClient {
	dfClient, err := dialogflow.NewSessionsClient(context.Background(), option.WithCredentialsFile(appConfig.DialogflowCredential))
	if err != nil {
		log.Panic("Error when connecting to dialogflow. Error " + err.Error())
	}
	return dfClient
}
