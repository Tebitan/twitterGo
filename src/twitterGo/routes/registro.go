package routes

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"

	"twitterGo/bd"
	"twitterGo/models"
)

const Body = "body"

func RegistroUser(ctx context.Context) models.ResponseApi {
	var user models.User
	var response models.ResponseApi
	response.ReponseCode = http.StatusBadRequest

	bodyRequest := ctx.Value(models.Key(Body)).(string)
	//convertimos string JSON to models.User
	err := json.Unmarshal([]byte(bodyRequest), &user)
	if err != nil {
		response.Message = err.Error()
		return response
	}
	msgValidate := ValidateRequestUser(user)
	if len(msgValidate) > 0 {
		response.Message = strings.Join(msgValidate, ", ")
		return response
	}
	msgSaveUser := SaveNewUser(user)
	if len(msgSaveUser) > 0 {
		response.ReponseCode = http.StatusInternalServerError
		response.Message = msgSaveUser
		return response
	}
	response.ReponseCode = http.StatusOK
	response.Message = "OK"
	return response
}

/*
Realiza la validacion de la request
@param bodyRequest models.User
@return []string Mensajes de validacion
*/
func ValidateRequestUser(bodyRequest models.User) []string {
	var msgValidate []string
	if len(bodyRequest.Email) == 0 {
		msgValidate = append(msgValidate, "Email es requerido")
	} else {
		_, existEmail, _ := bd.ExistUserEmail(bodyRequest.Email)
		if existEmail {
			msgValidate = append(msgValidate, "Ya existe un Usuario con el correo "+bodyRequest.Email)
		}
	}
	if len(bodyRequest.Password) < 6 {
		msgValidate = append(msgValidate, "La Contraseña debe tener por lo menos 6 caracteres")

	}
	return msgValidate
}

/*
Realiza el nuevo registro de Ususrio
@param bodyRequest models.User
@return string Mensajes de Error
*/
func SaveNewUser(bodyRequest models.User) string {
	_, saveUserBD, err := bd.SaveUser(bodyRequest)
	if err != nil {
		return "ERROR Save User " + err.Error()
	}
	if !saveUserBD {
		return "ERROR Save User No se ha logrado insertar el nuevo Usuario"
	}
	return ""
}
