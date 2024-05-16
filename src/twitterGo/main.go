package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"strings"

	events "github.com/aws/aws-lambda-go/events"
	lambda "github.com/aws/aws-lambda-go/lambda"

	"twitterGo/awsgo"
	"twitterGo/bd"
	"twitterGo/handlers"
	"twitterGo/models"
	"twitterGo/secretmanager"
)

const (
	SecretName  = "SecretName"
	BucketName  = "BucketName"
	UrlPrefix   = "UrlPrefix"
	Region      = "Region"
	NameService = "NameService"
)

// variables de entorno
var Envs = [5]string{SecretName, BucketName, UrlPrefix, Region, NameService}

func main() {
	lambda.Start(ExeLambda)
}

/*
Realiza la ejecuion de la lambda AWS
@param ctx  context.Context
@param request events.APIGatewayProxyRequest
@return events.APIGatewayProxyResponse
*/
func ExeLambda(ctx context.Context, request events.APIGatewayProxyRequest) (*events.APIGatewayProxyResponse, error) {
	var res *events.APIGatewayProxyResponse
	if !ValidateEnvs() {
		res = handlerError("Error en las variables de entorno", http.StatusBadRequest)
		return res, nil
	}
	awsgo.InitAWS(os.Getenv(Region))
	Secret, err := secretmanager.GetSecret(os.Getenv(SecretName))
	if err != nil {
		res = handlerError("Error en la lectura de Secret "+err.Error(), http.StatusBadRequest)
		return res, nil
	}
	saveEnvContext(Secret, request)
	err = bd.ConectionBD(awsgo.Ctx)
	if err != nil {
		res = handlerError("Error BBDD "+err.Error(), http.StatusInternalServerError)
		return res, nil
	}
	responseAPI := handlers.Handlers(awsgo.Ctx, request)
	if responseAPI.Data == nil {
		res = handlerError(responseAPI.Message, responseAPI.ReponseCode)
		return res, nil
	} else {
		return responseAPI.Data, nil
	}
}

/*
Realiza la validacion de las variables de entorno
@return bool
*/
func ValidateEnvs() bool {
	var existEnv bool = true
	for _, env := range Envs {
		_, existEnv = os.LookupEnv(env)
		if !existEnv {
			fmt.Printf("!!!No se encontro la variable de entorno [[ %s ]]", env)
			break
		}
	}
	return existEnv
}

/*
Realiza el manejo de errores
@param msg Mensaje
@param statusCode codigo de respuesta (400,500)
@return *events.APIGatewayProxyResponse
*/
func handlerError(msg string, statusCode int) *events.APIGatewayProxyResponse {
	return &events.APIGatewayProxyResponse{
		StatusCode: statusCode,
		Body:       msg,
		Headers: map[string]string{
			"Content-Type": "application/json",
		},
	}
}

/*
Realiza el guardado de ENV en el context
@param secret Data obtenida de AWS
@param request Peticion
*/
func saveEnvContext(secret models.Secret, request events.APIGatewayProxyRequest) {
	awsgo.Ctx = context.WithValue(awsgo.Ctx, models.Key(awsgo.Method), request.HTTPMethod)
	awsgo.Ctx = context.WithValue(awsgo.Ctx, models.Key(awsgo.User), secret.UserName)
	awsgo.Ctx = context.WithValue(awsgo.Ctx, models.Key(awsgo.Password), secret.Password)
	awsgo.Ctx = context.WithValue(awsgo.Ctx, models.Key(awsgo.Host), secret.Host)
	awsgo.Ctx = context.WithValue(awsgo.Ctx, models.Key(awsgo.Database), secret.Database)
	awsgo.Ctx = context.WithValue(awsgo.Ctx, models.Key(awsgo.Jwtsing), secret.JWTSing)
	awsgo.Ctx = context.WithValue(awsgo.Ctx, models.Key(awsgo.Body), request.Body)
	awsgo.Ctx = context.WithValue(awsgo.Ctx, models.Key(awsgo.BucketName), os.Getenv(BucketName))
	path := strings.Replace(request.PathParameters[os.Getenv(NameService)], os.Getenv(UrlPrefix), "", -1)
	awsgo.Ctx = context.WithValue(awsgo.Ctx, models.Key(awsgo.Path), path)
}
