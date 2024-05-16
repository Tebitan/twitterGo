package bd

import (
	"context"
	"fmt"
	"twitterGo/awsgo"
	"twitterGo/models"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const formatConexionString = "mongodb+srv://%s:%s@%s/?retryWrites=true&w=majority"

//variables publicas

// conexion
var MongoCN *mongo.Client
var DataBaseName string

/*
Realiza la conexion a la BBDD Mongo
@param ctx Context
@return error
*/
func ConectionBD(ctx context.Context) error {
	user := ctx.Value(models.Key(awsgo.User)).(string)
	pass := ctx.Value(models.Key(awsgo.Password)).(string)
	host := ctx.Value(models.Key(awsgo.Host)).(string)
	conectionString := fmt.Sprintf(formatConexionString, user, pass, host)
	var clientOptions = options.Client().ApplyURI(conectionString)
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		fmt.Println("Error en la conexion de BBDD \n" + err.Error())
		return err
	}
	err = client.Ping(ctx, nil)
	if err != nil {
		fmt.Println("Error en el PING de la conexion de BBDD \n" + err.Error())
		return err
	}

	fmt.Println("Conexion exitosa con la BBDD")
	MongoCN = client
	DataBaseName = ctx.Value(models.Key(awsgo.Database)).(string)
	return nil
}

/*
Valida si la conexion esta activa
@return bool
*/
func IsConneted() bool {
	err := MongoCN.Ping(context.TODO(), nil)
	return err == nil
}
