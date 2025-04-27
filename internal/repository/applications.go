package repository

import (
	"context"
	"fmt"
	"time"

	db "github.com/manan-sharma-5/seekajob-backend/internal/dbconfig"
	"github.com/manan-sharma-5/seekajob-backend/internal/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func ApplyJob(jobID string, userID string) error {
	client, err := db.GetMongoClient()
	if err != nil {
		return err
	}
	coll := client.Collection("job_applications")
	collUser := client.Collection("users")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	userPrimitive, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return fmt.Errorf("invalid user ID format")
	}
	jobPrimitive, err := primitive.ObjectIDFromHex(jobID)
	if err != nil {
		return fmt.Errorf("invalid job ID format")
	}
	user := collUser.FindOne(ctx, bson.M{"_id": userPrimitive})
	if user.Err() == mongo.ErrNoDocuments {
		return fmt.Errorf("user not found")
	}
	if user.Err() != nil {
		return user.Err()
	}

	existingApplication := coll.FindOne(context.Background(), bson.M{
		"job_id": jobPrimitive,
		"user_id": userPrimitive,
	})
	if existingApplication.Err() == nil {
		return fmt.Errorf("user has already applied for this job")
	}
	if existingApplication.Err() != mongo.ErrNoDocuments {
		return existingApplication.Err()
	}

	userDetails := model.User{}
	err = user.Decode(&userDetails)
	if err != nil {
		fmt.Println("Error decoding user:", err)
	}
	if !userDetails.IsCandidate {
		return fmt.Errorf("user is not a candidate")
	}

	jobIDObjectID, err := primitive.ObjectIDFromHex(jobID)
	if err != nil {
		return fmt.Errorf("invalid job ID format")
	}
	userIDObjectID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return fmt.Errorf("invalid user ID format")
	}

	application := &model.JobApplicant{
		JobID:    jobIDObjectID,
		UserID:   userIDObjectID,
		Status:   "applied",
		Resume: userDetails.Resume,
		Name:  userDetails.Name,
		Email: userDetails.Email,
	}

	_, err = coll.InsertOne(ctx, application)
	if err != nil {
		return err
	}

	return nil
}