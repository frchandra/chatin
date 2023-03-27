package util

import (
	dialogflow "cloud.google.com/go/dialogflow/apiv2"
	"cloud.google.com/go/dialogflow/apiv2/dialogflowpb"
	"context"
	"fmt"
	"github.com/frchandra/chatin/config"
)

type DialogflowUtil struct {
	projectId     string
	lang          string
	timeZone      string
	sessionClient *dialogflow.SessionsClient
	config        *config.AppConfig
}

func NewDialogflowUtil(sessionClient *dialogflow.SessionsClient, config *config.AppConfig) *DialogflowUtil {
	return &DialogflowUtil{sessionClient: sessionClient, config: config}
}

func (u DialogflowUtil) DetectIntent(text string) (string, error) {
	sessionPath := fmt.Sprintf("projects/%s/agent/sessions/%s", u.config.DialogflowProjectId, u.config.DialogflowSessionId)
	textInput := dialogflowpb.TextInput{Text: text, LanguageCode: u.config.DialogflowLanguage}
	queryTextInput := dialogflowpb.QueryInput_Text{Text: &textInput}
	queryInput := dialogflowpb.QueryInput{Input: &queryTextInput}
	request := dialogflowpb.DetectIntentRequest{Session: sessionPath, QueryInput: &queryInput}

	response, err := u.sessionClient.DetectIntent(context.Background(), &request)
	if err != nil {
		return "", err
	}

	queryResult := response.GetQueryResult()
	fulfillmentText := queryResult.GetFulfillmentText()
	return fulfillmentText, nil

}
