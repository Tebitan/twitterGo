package utils

import (
	"golang.org/x/crypto/bcrypt"
)

/*
Realiza la Encriptacion
@param data dato a encriptar
@return (2)(string, error)
*/
func Encrypt(data string) (string, error) {
	costo := 8
	bytes, err := bcrypt.GenerateFromPassword([]byte(data), costo)
	if err != nil {
		return err.Error(), err
	}
	return string(bytes), nil
}

/*
Realiza la comparacion de contraseña
@param passwordDecode clave encriptar
@param passwordEncode clave desencriptada
@return error
*/
func ComparePassword(passwordDecode string, passwordEncode string) error {
	passwordDecodeBytes := []byte(passwordDecode)
	passwordEncodeBytes := []byte(passwordEncode)
	return bcrypt.CompareHashAndPassword(passwordEncodeBytes, passwordDecodeBytes)
}
