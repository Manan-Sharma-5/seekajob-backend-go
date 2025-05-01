package repository

import (
	"context"
	"log"
	"time"

	db "github.com/manan-sharma-5/seekajob-backend/internal/dbconfig"
	"github.com/manan-sharma-5/seekajob-backend/internal/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
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

func GetUserByID(id string) (*model.UserWithApplicationsAndJobs, error) {
    dbClient, err := db.GetMongoClient()
    if err != nil {
        return nil, err
    }
    collection := dbClient.Collection("users")

    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()

    oid, err := primitive.ObjectIDFromHex(id)
    if err != nil {
        return nil, err
    }

	pipeline := mongo.Pipeline{
		{
			primitive.E{Key: "$match", Value: bson.D{
				{Key: "_id", Value: oid},
			}},
		},
		{
			primitive.E{Key: "$lookup", Value: bson.D{
				{Key: "from", Value: "job_applications"},
				{Key: "localField", Value: "_id"},
				{Key: "foreignField", Value: "userID"},
				{Key: "as", Value: "jobApplications"},
			}},
		},
		{
			primitive.E{Key: "$unwind", Value: bson.D{
				{Key: "path", Value: "$jobApplications"},
				{Key: "preserveNullAndEmptyArrays", Value: true},
			}},
		},
		{
			primitive.E{Key: "$lookup", Value: bson.D{
				{Key: "from", Value: "jobs"},
				{Key: "localField", Value: "jobApplications.jobID"},
				{Key: "foreignField", Value: "_id"},
				{Key: "as", Value: "jobApplications.job"},
			}},
		},
		{
			primitive.E{Key: "$unwind", Value: bson.D{
				{Key: "path", Value: "$jobApplications.job"},
				{Key: "preserveNullAndEmptyArrays", Value: true},
			}},
		},
		{
			primitive.E{Key: "$group", Value: bson.D{
				{Key: "_id", Value: "$_id"},
				{Key: "name", Value: bson.D{{Key: "$first", Value: "$name"}}},
				{Key: "resume", Value: bson.D{{Key: "$first", Value: "$resume"}}},
				{Key: "email", Value: bson.D{{Key: "$first", Value: "$email"}}},
				{Key: "isCandidate", Value: bson.D{{Key: "$first", Value: "$isCandidate"}}},
				{Key: "createdAt", Value: bson.D{{Key: "$first", Value: "$createdAt"}}},
				{Key: "updatedAt", Value: bson.D{{Key: "$first", Value: "$updatedAt"}}},
				{Key: "jobApplications", Value: bson.D{{Key: "$push", Value: "$jobApplications"}}},
				{Key: "appliedJobs", Value: bson.D{{Key: "$push", Value: "$jobApplications.job"}}},
			}},
		},
	}
	
	

    cursor, err := collection.Aggregate(ctx, pipeline)
    if err != nil {
        return nil, err
    }

    var users []model.UserWithApplicationsAndJobs
    if err := cursor.All(ctx, &users); err != nil {
        return nil, err
    }

    if len(users) == 0 {
        return nil, mongo.ErrNoDocuments
    }

    return &users[0], nil
}
