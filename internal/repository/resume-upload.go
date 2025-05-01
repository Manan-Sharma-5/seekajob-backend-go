package repository

import (
	"context"
	"fmt"
	"os"
	"time"

	db "github.com/manan-sharma-5/seekajob-backend/internal/dbconfig"
	"github.com/manan-sharma-5/seekajob-backend/internal/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func UploadResume(resumeFile []byte, userID string) (string, error) {
	client, err := db.GetMongoClient()
	if err != nil {
		return "", err
	}
	coll := client.Collection("users")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	userPrimitive, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return "", fmt.Errorf("invalid user ID format")
	}

	user := coll.FindOne(ctx, bson.M{"_id": userPrimitive})
	if user.Err() == mongo.ErrNoDocuments {
		return "", fmt.Errorf("user not found")
	}
	if user.Err() != nil {
		return "", user.Err()
	}

	userDetails := model.User{}
	err = user.Decode(&userDetails)
	if err != nil {
		fmt.Println("Error decoding user:", err)
	}

	if !userDetails.IsCandidate {
		return "", fmt.Errorf("user is not a candidate")
	}

	resumePath := fmt.Sprintf("./resumes/%s_resume.pdf", userID)

	err = os.WriteFile(resumePath, resumeFile, 0644)
	if err != nil {
		return "", fmt.Errorf("failed to save resume: %v", err)
	}

	// Update the user's resume path in the database
	_, err = coll.UpdateOne(ctx, bson.M{"_id": userPrimitive}, bson.M{
		"$set": bson.M{
			"resume": resumePath,
			"updatedAt":  time.Now(),
		},
	})
	if err != nil {
		return "", fmt.Errorf("failed to update resume path in database: %v", err)
	}

	return resumePath, nil
}