package bd

import (
	"context"
	"twitterGo/models"
	"twitterGo/utils"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

const tableUsers = "usuarios"
const colId = "_id"
const colName = "name"
const colLastName = "lastName"
const colBirthdate = "birthdate"
const colEmail = "email"
const colAvatar = "avatar"
const colBanner = "banner"
const colBiography = "biography"
const colLocation = "location"
const colWebsite = "website"

/*
Valida si ya existe un Usuario con ese EMAIL
@return models.User, bool, error
*/
func ExistUserEmail(email string) (models.User, bool, string) {
	ctx := context.TODO()
	db := MongoCN.Database(DataBaseName)
	tableUser := db.Collection(tableUsers)
	condition := bson.M{colEmail: email}
	var user models.User
	err := tableUser.FindOne(ctx, condition).Decode(&user)
	if err != nil {
		return user, false, ""
	}
	return user, true, user.ID.Hex()
}

/*
Guardar Usuario BBDD
@param user models.User
@return (3)string, bool, error
*/
func SaveUser(user models.User) (string, bool, error) {
	pasEncrypt, err := utils.Encrypt(user.Password)
	if err != nil {
		return "", false, err
	}
	user.Password = pasEncrypt
	ctx := context.TODO()
	db := MongoCN.Database(DataBaseName)
	tableUser := db.Collection(tableUsers)
	newUser, err := tableUser.InsertOne(ctx, user)
	if err != nil {
		return "", false, err
	}
	userId, _ := newUser.InsertedID.(primitive.ObjectID)
	return userId.String(), true, nil
}

/*
Buscar el Usuario por ID
@param id
@return models.User, error
*/
func FindUserById(id string) (models.User, error) {
	ctx := context.TODO()
	db := MongoCN.Database(DataBaseName)
	tableUser := db.Collection(tableUsers)
	decodeID, _ := primitive.ObjectIDFromHex(id)
	condition := bson.M{colId: decodeID}
	var user models.User
	err := tableUser.FindOne(ctx, condition).Decode(&user)
	return user, err
}

/*
Modifica el Perfil del  Usuario BBDD
@param id string 'Id del user'
@param user models.User
@return error
*/
func UpdateUser(id string, user models.User) error {
	ctx := context.TODO()
	db := MongoCN.Database(DataBaseName)
	tableUser := db.Collection(tableUsers)
	decodeID, _ := primitive.ObjectIDFromHex(id)
	updateString := bson.M{"$set": dataUpdate(user)}
	condition := bson.M{colId: bson.M{"$eq": decodeID}}
	_, err := tableUser.UpdateOne(ctx, condition, updateString)
	return err
}

func dataUpdate(user models.User) map[string]interface{} {
	registro := make(map[string]interface{})
	if len(user.Name) > 0 {
		registro[colName] = user.Name
	}
	if len(user.LastName) > 0 {
		registro[colLastName] = user.LastName
	}
	if user.Birthdate.IsZero() {
		registro[colBirthdate] = user.Birthdate
	}
	if len(user.Avatar) > 0 {
		registro[colAvatar] = user.Avatar
	}
	if len(user.Banner) > 0 {
		registro[colBanner] = user.Banner
	}
	if len(user.Biography) > 0 {
		registro[colBiography] = user.Biography
	}
	if len(user.Location) > 0 {
		registro[colLocation] = user.Location
	}
	if len(user.Website) > 0 {
		registro[colWebsite] = user.Website
	}
	return registro
}
