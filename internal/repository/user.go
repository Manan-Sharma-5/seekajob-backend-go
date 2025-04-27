package repository

import (
	"context"
	"log"
	"time"

	db "github.com/manan-sharma-5/seekajob-backend/internal/dbconfig"
	"github.com/manan-sharma-5/seekajob-backend/internal/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
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
	err = collection.FindOne(ctx, bson.M{"email": email}).Decode(&user)
	if err != nil {
		log.Println("Error finding user:", err)
		if err.Error() == "mongo: no documents in result" {
			return false, nil
		}
		return true, err
	}
	if (user.ID == "") {
		return true, nil
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
	user.ID = userMongo.InsertedID.(primitive.ObjectID).Hex()

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
	err = collection.FindOne(ctx, bson.M{"email": email}).Decode(&user)
	if err != nil {
		panic(err)
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		log.Println("Password mismatch:", err)
		return nil, err
	}
	if user.IsCandidate != isCandidate {
		return nil, nil
	}
	return &user, nil
}

func GetUserByID(id string) (*model.User, error) {
	dbClient, err := db.GetMongoClient()
	if err != nil {
		panic(err)
	}
	collection := dbClient.Collection("users")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	var user model.User
	log.Println("ID from here:", id)
	// Convert string ID to ObjectID
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		log.Println("Error converting ID:", err)
		return nil, err
	}
	err = collection.FindOne(ctx, bson.M{"_id": oid}).Decode(&user)
	if err != nil {
		panic(err)
	}
	return &user, nil
}
	