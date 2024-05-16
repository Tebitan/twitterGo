package routes

import (
	"context"
	"encoding/json"
	"net/http"

	"twitterGo/bd"
	"twitterGo/models"

	events "github.com/aws/aws-lambda-go/events"
)

func GetPerfil(request events.APIGatewayProxyRequest) models.ResponseApi {
	var response models.ResponseApi
	response.ReponseCode = http.StatusBadRequest

	userID := request.QueryStringParameters["id"]
	if len(userID) < 1 {
		response.Message = "El parámetro ID es obligatorio"
		return response
	}

	user, err := bd.FindUserById(userID)
	if err != nil {
		response.ReponseCode = http.StatusInternalServerError
		response.Message = "Ocurrió un error al intentar buscar el registro " + err.Error()
		return response
	}
	//lo colocamos vacio , para que no lo envie en la respuesta
	user.Password = ""
	userJSON, err := json.Marshal(user)
	if err != nil {
		response.ReponseCode = http.StatusInternalServerError
		response.Message = "Error formater JSON " + err.Error()
		return response
	}
	response.ReponseCode = http.StatusOK
	response.Message = string(userJSON)
	return response
}

func UpdatePerfil(ctx context.Context, claim models.Claim) models.ResponseApi {
	var response models.ResponseApi
	response.ReponseCode = http.StatusBadRequest
	bodyRequest := ctx.Value(models.Key(Body)).(string)
	//convertimos string JSON to models.User
	var user models.User
	err := json.Unmarshal([]byte(bodyRequest), &user)
	if err != nil {
		response.Message = err.Error()
		return response
	}
	msgUpdPerfil := updatePerfilBD(claim.ID.Hex(), user)
	if len(msgUpdPerfil) > 0 {
		response.ReponseCode = http.StatusInternalServerError
		response.Message = msgUpdPerfil
		return response
	}
	response.ReponseCode = http.StatusOK
	response.Message = "OK"
	return response
}

func updatePerfilBD(id string, user models.User) string {
	err := bd.UpdateUser(id, user)
	if err != nil {
		return "ERROR Update Perfil " + err.Error()
	}
	return ""
}
