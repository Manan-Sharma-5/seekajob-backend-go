package repository

import (
	"context"
	"fmt"
	"time"

	db "github.com/manan-sharma-5/seekajob-backend/internal/dbconfig"
	"github.com/manan-sharma-5/seekajob-backend/internal/model"
	"go.mongodb.org/mongo-driver/mongo"
)

func ApplyJob(jobID string, userID string) error {
	client, err := db.GetMongoClient()
	if err != nil {
		return err
	}
	coll := client.Collection("job_applications")
	collUser := client.Collection("users")
	user := collUser.FindOne(context.Background(), model.User{ID: userID})
	if user.Err() != nil {
		return user.Err()
	}
	if user.Err() == mongo.ErrNoDocuments {
		return fmt.Errorf("user not found")
	}

	existingApplication := coll.FindOne(context.Background(), model.JobApplicant{JobID: jobID, UserID: userID})
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

	application := &model.JobApplicant{
		JobID:    jobID,
		UserID:   userID,
		Status:   "applied",
		Resume: userDetails.Resume,
		Name:  userDetails.Name,
		Email: userDetails.Email,
	}


	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	_, err = coll.InsertOne(ctx, application)
	if err != nil {
		return err
	}

	return nil
}