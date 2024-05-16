package routes

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"twitterGo/bd"
	"twitterGo/models"
	"twitterGo/strategies"
	"twitterGo/utils"

	events "github.com/aws/aws-lambda-go/events"
)

const BODY = "body"

func Login(ctx context.Context) models.ResponseApi {
	var user models.User
	var response models.ResponseApi
	response.ReponseCode = http.StatusBadRequest

	bodyRequest := ctx.Value(models.Key(BODY)).(string)
	//convertimos string JSON to models.User
	err := json.Unmarshal([]byte(bodyRequest), &user)
	if err != nil {
		response.Message = "Usuario y/o Contraseña Inválidos " + err.Error()
		return response
	}
	msgValidate := ValidateRequestLogin(user)
	if len(msgValidate) > 0 {
		response.Message = strings.Join(msgValidate, ", ")
		return response
	}

	token, err := strategies.GenerateJWT(ctx, user)
	if err != nil {
		response.Message = "Error Generando el Token " + err.Error()
		return response
	}
	responseLogin := models.ResponseLogin{
		Token: token,
	}
	responseLoginJSON, err2 := json.Marshal(responseLogin)
	if err2 != nil {
		response.Message = "Error formater JSON  " + err2.Error()
		return response
	}
	//create cookie
	cookie := &http.Cookie{
		Name:    "token",
		Value:   token,
		Expires: time.Now().Add(time.Hour * 24),
	}
	cookiesStr := cookie.String()
	resposeMethod := &events.APIGatewayProxyResponse{
		StatusCode: http.StatusOK,
		Body:       string(responseLoginJSON),
		Headers: map[string]string{
			"Content-Type":                "application/json",
			"Access-Control-Allow-Origin": "*",
			"Set-Cookie":                  cookiesStr,
		},
	}
	response.ReponseCode = http.StatusOK
	response.Message = string(token)
	response.Data = resposeMethod
	return response
}

func ValidateRequestLogin(bodyRequest models.User) []string {
	var msgValidate []string
	if len(bodyRequest.Email) == 0 {
		msgValidate = append(msgValidate, "Email es requerido")
	}
	if len(bodyRequest.Password) < 6 {
		msgValidate = append(msgValidate, "La Contraseña debe tener por lo menos 6 caracteres")
	}
	if len(msgValidate) == 0 {
		userBD, existEmail, _ := bd.ExistUserEmail(bodyRequest.Email)
		if !existEmail {
			msgValidate = append(msgValidate, "No existe un Usuario con el correo "+bodyRequest.Email)
		} else {
			err := utils.ComparePassword(bodyRequest.Password, userBD.Password)
			if err != nil {
				msgValidate = append(msgValidate, "Contraseña Inválida")
			}
		}
	}
	return msgValidate
}
