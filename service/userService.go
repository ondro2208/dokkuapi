package service

import (
	"context"
	"errors"
	log "github.com/ondro2208/dokkuapi/logger"
	"github.com/ondro2208/dokkuapi/model"
	str "github.com/ondro2208/dokkuapi/store"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"net/http"
	"strconv"
	"time"
)

type UsersService interface {
	CreateUser(githubUser *model.GithubUser) (*model.User, int, string)
	GetExistingUser(githubUser *model.GithubUser) (*model.User, int, string)
	GetExistingUserById(userIdHex string) (*model.User, int, string)
	DeleteExistingUser(userIdHex string) error
	GetUserApplications(userIdHex string) ([]model.Application, int, string)
	UpdateUserWithApplication(appName string, userId primitive.ObjectID) (*model.Application, int, string)
	SetUserApplicationServices(newApp model.Application, userId primitive.ObjectID) (*model.Application, int, string)
	SetUserApplicationName(userId primitive.ObjectID, appId primitive.ObjectID, appName string) (int, string)
	DeleteUserApplication(userId primitive.ObjectID, appId primitive.ObjectID) (int, string, bool)
}

func NewUsersService(serviceStore *str.Store) UsersService {
	return &UsersServiceContext{store: serviceStore}
}

type UsersServiceContext struct {
	store *str.Store
}

func getCollection(client *mongo.Client, dbName string, collectionName string) (*mongo.Collection, context.Context) {
	collection := client.Database(dbName).Collection(collectionName)
	ctx, _ := context.WithTimeout(context.Background(), 30*time.Second)
	return collection, ctx
}

func (us *UsersServiceContext) CreateUser(githubUser *model.GithubUser) (*model.User, int, string) {
	var user = new(model.User)
	users, ctx := getCollection(us.store.Client, us.store.DbName, "users")
	err := users.FindOne(ctx, model.User{GithubId: githubUser.Id}).Decode(&user)
	if err != nil {
		if err != mongo.ErrNoDocuments {
			return nil, http.StatusInternalServerError, err.Error()
		}
		result, _ := users.InsertOne(ctx, model.User{GithubId: githubUser.Id, Username: githubUser.Login})
		users.FindOne(ctx, model.User{Id: result.InsertedID.(primitive.ObjectID)}).Decode(&user)
		return user, http.StatusCreated, "User created"
	} else {
		return nil, http.StatusConflict, "User already registered"
	}
}

func (us *UsersServiceContext) GetExistingUser(githubUser *model.GithubUser) (*model.User, int, string) {
	var user = new(model.User)
	users, ctx := getCollection(us.store.Client, us.store.DbName, "users")
	err := users.FindOne(ctx, model.User{GithubId: githubUser.Id}).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, http.StatusNotFound, "User doesn't exist"
		}
		log.ErrorLogger.Println("Can't find existing user")
		return nil, http.StatusInternalServerError, err.Error()
	}
	return user, http.StatusCreated, "User is successfully logged in"

}

func (us *UsersServiceContext) GetExistingUserById(userIdHex string) (*model.User, int, string) {
	idPrimitive, err := primitive.ObjectIDFromHex(userIdHex)
	if err != nil {
		log.ErrorLogger.Println("Parsing ObjectId from hex error: ", err.Error())
		return nil, http.StatusInternalServerError, "Can't find user"
	}
	var user = new(model.User)
	users, ctx := getCollection(us.store.Client, us.store.DbName, "users")
	err = users.FindOne(ctx, model.User{Id: idPrimitive}).Decode(&user)
	if err != nil {
		return nil, http.StatusInternalServerError, err.Error()
	}
	return user, http.StatusOK, "User founded by id"

}

func (us *UsersServiceContext) DeleteExistingUser(userIdHex string) error {
	idPrimitive, err := primitive.ObjectIDFromHex(userIdHex)
	if err != nil {
		log.ErrorLogger.Println("Parsing ObjectId from hex error: ", err.Error())
		return err
	}
	users, ctx := getCollection(us.store.Client, us.store.DbName, "users")
	res, err := users.DeleteOne(ctx, bson.M{"_id": idPrimitive})
	if err != nil {
		log.ErrorLogger.Println("Delete one user error: ", err.Error())
		return err
	}
	if res.DeletedCount != 1 {
		message := "Delete " + strconv.FormatInt(res.DeletedCount, 10) + " instead of 1"
		log.ErrorLogger.Println(message)
		return errors.New(message)
	}
	return nil
}

func (us *UsersServiceContext) GetUserApplications(userIdHex string) ([]model.Application, int, string) {
	idPrimitive, err := primitive.ObjectIDFromHex(userIdHex)
	if err != nil {
		log.ErrorLogger.Println("Parsing ObjectId from hex error: ", err.Error())
		return nil, http.StatusInternalServerError, "Can't find user"
	}
	var user = new(model.User)
	users, ctx := getCollection(us.store.Client, us.store.DbName, "users")
	err = users.FindOne(ctx, model.User{Id: idPrimitive}).Decode(&user)
	if err != nil {
		return nil, http.StatusInternalServerError, err.Error()
	}
	return user.Applications, http.StatusOK, "User's apps founded"
}

func (us *UsersServiceContext) UpdateUserWithApplication(appName string, userId primitive.ObjectID) (*model.Application, int, string) {
	newApp := model.Application{
		Name: appName,
		Id:   primitive.NewObjectID(),
	}
	users, ctx := getCollection(us.store.Client, us.store.DbName, "users")
	result, err := users.UpdateOne(
		ctx,
		bson.M{"_id": userId},
		bson.D{
			{"$push", bson.M{"applications": newApp}},
		},
	)
	log.GeneralLogger.Println("Result after updating database ", result)
	if err != nil {
		log.ErrorLogger.Println(err)
		return nil, http.StatusInternalServerError, "Unable to store application"
	}

	if result.MatchedCount == 0 {
		return nil, http.StatusInternalServerError, "No user updated"
	}

	return &newApp, http.StatusCreated, "Application successfully created"
}

func (us *UsersServiceContext) SetUserApplicationServices(newApp model.Application, userId primitive.ObjectID) (*model.Application, int, string) {
	users, ctx := getCollection(us.store.Client, us.store.DbName, "users")
	result, err := users.UpdateOne(
		ctx,
		bson.M{"_id": userId, "applications._id": newApp.Id},
		bson.D{
			{"$set", bson.M{"applications.$.services": newApp.Services}},
		},
	)
	log.GeneralLogger.Println("Result after updating database ", result)
	if err != nil {
		log.ErrorLogger.Println(err)
		return nil, http.StatusUnprocessableEntity, "Unable to store services"
	}

	if result.MatchedCount == 0 {
		return nil, http.StatusInternalServerError, "No app updated"
	}

	return &newApp, http.StatusCreated, "Application successfully updated"
}

func (us *UsersServiceContext) DeleteUserApplication(userId primitive.ObjectID, appId primitive.ObjectID) (int, string, bool) {
	users, ctx := getCollection(us.store.Client, us.store.DbName, "users")
	result, err := users.UpdateOne(
		ctx,
		bson.M{"_id": userId},
		bson.D{
			{"$pull", bson.M{"applications": bson.M{"_id": appId}}},
		},
	)
	log.GeneralLogger.Println("Result after updating database ", result.MatchedCount)
	if err != nil {
		log.ErrorLogger.Println(err)
		return http.StatusUnprocessableEntity, "Can't delete application", true
	}

	if result.MatchedCount == 0 {
		return http.StatusInternalServerError, "No application deleted", false
	}

	return http.StatusAccepted, "Application deleted successfully", false
}

func (us *UsersServiceContext) SetUserApplicationName(userId primitive.ObjectID, appId primitive.ObjectID, appName string) (int, string) {
	users, ctx := getCollection(us.store.Client, us.store.DbName, "users")
	result, err := users.UpdateOne(
		ctx,
		bson.M{"_id": userId, "applications._id": appId},
		bson.D{
			{"$set", bson.M{"applications.$.appName": appName}},
		},
	)
	log.GeneralLogger.Println("Result after updating database ", result.ModifiedCount)
	if err != nil {
		log.ErrorLogger.Println(err)
		return http.StatusUnprocessableEntity, "Unable to rename application"
	}

	if result.MatchedCount == 0 {
		return http.StatusInternalServerError, "No app updated"
	}

	return http.StatusCreated, "Application name successfully updated"
}
