package awsgo

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
)

var Ctx context.Context
var Cfg aws.Config
var err error

// Varaibles del context Ctx
const (
	Method     = "method"
	User       = "user"
	Password   = "password"
	Host       = "host"
	Database   = "database"
	Jwtsing    = "jwtSign"
	Body       = "body"
	BucketName = "bucketName"
	Path       = "path"
)

/*
Realiza la Inicializacion de la configuracion de AWS
@param region resgion EJ: "us-east-1"
@return bool
*/
func InitAWS(region string) {
	//Crea un context vacio
	Ctx = context.TODO()
	Cfg, err = config.LoadDefaultConfig(Ctx, config.WithDefaultRegion(region))
	if err != nil {
		panic("Error al cargar la configuracion .aws/config " + err.Error())
	}
}
