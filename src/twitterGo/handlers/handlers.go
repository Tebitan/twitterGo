package handlers

import (
	"context"
	"fmt"
	"net/http"
	"slices"
	"twitterGo/awsgo"
	"twitterGo/models"
	"twitterGo/routes"
	"twitterGo/strategies"

	"github.com/aws/aws-lambda-go/events"
)

const (
	PathRegistro      = "registro"
	PathLogin         = "login"
	PathPerfil        = "perfil"
	PathObtenerAvatar = "obteneravatar"
	PathObtenerBanner = "obtenerbanner"
	Autorizacion      = "autorizacion"
)

/*
Realiza la funcion de un Enrutador
@param ctx  context.Context
@param request events.APIGatewayProxyRequest
@return models.ResponseApi
*/
func Handlers(ctx context.Context, request events.APIGatewayProxyRequest) models.ResponseApi {
	path := ctx.Value(models.Key(awsgo.Path)).(string)
	method := ctx.Value(models.Key(awsgo.Method)).(string)
	fmt.Printf("Handlers [[START]] METHOD: %s , PATH: %s , REQUEST_BODY:{ %+v } ", method, path, request.Body)
	var response models.ResponseApi
	response.ReponseCode = http.StatusBadRequest
	response.Message = "Method Invalid"
	isOk, statusCode, msg, claim := validateAuthorization(path, ctx, request)
	if !isOk {
		response.ReponseCode = statusCode
		response.Message = msg
		return response
	}

	switch method {
	case http.MethodGet:
		switch path {
		case PathPerfil:
			return routes.GetPerfil(request)
		}

	case http.MethodPost:
		switch path {
		case PathRegistro:
			return routes.RegistroUser(ctx)
		case PathLogin:
			return routes.Login(ctx)
		}

	case http.MethodPut:
		switch path {
		case PathPerfil:
			return routes.UpdatePerfil(ctx, claim)

		}
	case http.MethodDelete:
		switch path {

		}

	}
	fmt.Printf("Handlers [[END]] METHOD: %s , PATH: %s , RESPONSE:{ %+v } ", path, method, response)
	return response
}

/*
Realiza la validacion de la autorizacion
@param path URL
@param ctx  context.Context
@param request events.APIGatewayProxyRequest
@return (4)(bool,int,string,models.Claim)
*/
func validateAuthorization(path string, ctx context.Context, request events.APIGatewayProxyRequest) (bool, int, string, models.Claim) {
	//rutas con exepcion
	rotesException := []string{PathRegistro, PathLogin, PathObtenerAvatar, PathObtenerBanner}
	if slices.Contains(rotesException, path) {
		return true, http.StatusOK, "", models.Claim{}
	}
	token := request.Headers[Autorizacion]
	//existe el token
	if len(token) == 0 {
		return false, http.StatusUnauthorized, "Token Requerido", models.Claim{}
	}
	//validamos el token
	claim, isOK, msg, err := strategies.ProcessToken(token, ctx.Value(models.Key(awsgo.Jwtsing)).(string))
	if !isOK {
		errorString := msg
		if err != nil {
			errorString = err.Error()
		}
		fmt.Println("Error en el TOKEN " + errorString)
		return false, http.StatusUnauthorized, errorString, models.Claim{}
	}
	return true, http.StatusOK, msg, *claim
}
