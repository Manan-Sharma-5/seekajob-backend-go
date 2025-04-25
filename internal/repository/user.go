package repository

import (
	"context"
	"time"

	db "github.com/manan-sharma-5/seekajob-backend/internal/dbconfig"
	"github.com/manan-sharma-5/seekajob-backend/internal/model"
	"golang.org/x/crypto/bcrypt"
)

func IsEmailTaken (email string) (bool, error) {
	dbClient, err := db.GetMongoClient()
	if err != nil {
		panic(err)
	}
	collection := dbClient.Collection("users")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	var user model.User
	err = collection.FindOne(ctx, model.User{Email: email}).Decode(&user)
	if err != nil {
		if err.Error() == "mongo: no documents in result" {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

func CreateUser(user *model.User) (*model.User, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)

	if err != nil {
		panic(err)
	}

	user.Password = string(hashedPassword)

	dbClient, err := db.GetMongoClient()
	if err != nil {
		panic(err)
	}
	collection := dbClient.Collection("users")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	userMongo, err := collection.InsertOne(ctx, user)
	if err != nil {
		panic(err)
	}
	user.ID = userMongo.InsertedID.(string)
	return user, nil
}

func VerifyUser(email string, password string, isCandidate bool) (*model.User, error) {
	dbClient, err := db.GetMongoClient()
	if err != nil {
		panic(err)
	}
	collection := dbClient.Collection("users")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	var user model.User
	err = collection.FindOne(ctx, model.User{Email: email}).Decode(&user)
	if err != nil {
		panic(err)
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		panic(err)
	}
	if user.IsCandidate != isCandidate {
		return nil, nil
	}
	return &user, nil
}