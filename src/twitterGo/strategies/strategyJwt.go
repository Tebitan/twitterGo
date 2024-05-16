package strategies

import (
	"context"
	"errors"
	"strings"
	"time"
	"twitterGo/bd"

	"twitterGo/awsgo"
	"twitterGo/models"

	"github.com/golang-jwt/jwt/v5"
)

const Bearer = "Bearer"

/*
Realiza la validacion del TOKEN JWT
@param token
@param JWTSing Semilla de encripcion
@return models.Claim,bool,stringerror
*/
func ProcessToken(token string, JWTSing string) (*models.Claim, bool, string, error) {
	key := []byte(JWTSing)
	var claims models.Claim
	splitToken := strings.Split(token, Bearer)
	if len(splitToken) != 2 {
		return &claims, false, string(""), errors.New("formato de token invalido")
	}
	token = strings.TrimSpace(splitToken[1])
	newToken, err := jwt.ParseWithClaims(token, &claims, func(tk *jwt.Token) (interface{}, error) {
		return key, nil
	})
	if err != nil {
		return &claims, false, string(""), errors.New("error generando el token")
	} else if !newToken.Valid {
		return &claims, false, string(""), errors.New("token invalido")
	}
	user, existUser, _ := bd.ExistUserEmail(claims.Email)
	claims.ID = user.ID
	return &claims, existUser, claims.ID.Hex(), nil
}

/*
Realiza la Creacion del TOKEN JWT
@param ctx context.Context
@param user models.User
@return string,error
*/
func GenerateJWT(ctx context.Context, user models.User) (string, error) {
	jwtSign := ctx.Value(models.Key(awsgo.Jwtsing)).(string)
	jwtSignBytes := []byte(jwtSign)
	payload := jwt.MapClaims{
		"email": user.Email,
		"exp":   time.Now().Add(time.Hour * 24).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, payload)
	tokenStr, err := token.SignedString(jwtSignBytes)
	return tokenStr, err
}
