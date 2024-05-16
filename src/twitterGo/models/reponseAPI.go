package models

import "github.com/aws/aws-lambda-go/events"

type ResponseApi struct {
	ReponseCode int                             `json:"reponseCode"`
	Message     string                          `json:"message"`
	Data        *events.APIGatewayProxyResponse `json:"data"`
}
