package repository

import (
	"context"
	"time"

	db "github.com/manan-sharma-5/seekajob-backend/internal/dbconfig"
	"github.com/manan-sharma-5/seekajob-backend/internal/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func CreateJob(job *model.Job, userID string) error {
	client, err := db.GetMongoClient()
	coll := client.Collection("jobs")
    compColl := client.Collection("companies")
	if err != nil {
		return err
	}

    companyName := job.Company.Name

    // Check if the company exists
    var company model.Company
    err = compColl.FindOne(context.Background(), bson.M{"name": companyName}).Decode(&company)
    if err != nil {
        if err == mongo.ErrNoDocuments {
            // If the company does not exist, create it
            company = model.Company{
                Name:      companyName,
                CreatedAt: time.Now(),
                UpdatedAt: time.Now(),
            }
            _, err = compColl.InsertOne(context.Background(), company)
            if err != nil {
                return err
            }
        } else {
            return err
        }
    }

    // Set the company ID in the job
    job.CompanyID = company.ID
	job.RecruiterID = userID
	job.CreatedAt = time.Now()
	job.UpdatedAt = time.Now()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	_, err = coll.InsertOne(ctx, job)
	if err != nil {
		return err
	}

	return nil
}

func GetAllJobs() ([]model.Job, error) {
    client, err := db.GetMongoClient()
    if err != nil {
        return nil, err
    }

    coll := client.Collection("jobs")
    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()

    cursor, err := coll.Find(ctx, bson.D{})
    if err != nil {
        return nil, err
    }
    defer cursor.Close(ctx)

    var jobs []model.Job
    if err := cursor.All(ctx, &jobs); err != nil {
        return nil, err
    }

    return jobs, nil
}

func GetJobApplicantWithRelations(applicantID string) (*model.JobApplicant, error) {
    coll, err := db.GetMongoClient()
    if err != nil {
        return nil, err
    }


    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()

    pipeline := mongo.Pipeline{
        // Match the job applicant
        {{Key: "$match", Value: bson.D{{Key: "_id", Value: applicantID}}}},
        
        // Lookup the Job
        {
            {Key: "$lookup", Value: bson.D{
                {Key: "from", Value: "jobs"},
                {Key: "localField", Value: "jobID"},
                {Key: "foreignField", Value: "_id"},
                {Key: "as", Value: "job"},
            }},
        },
        // Unwind the job array
        {{Key: "$unwind", Value: bson.D{{Key: "path", Value: "$job"}, {Key: "preserveNullAndEmptyArrays", Value: true}}}},

        // Lookup the Candidate (user)
        {
            {Key: "$lookup", Value: bson.D{
                {Key: "from", Value: "users"},
                {Key: "localField", Value: "userID"},
                {Key: "foreignField", Value: "_id"},
                {Key: "as", Value: "user"},
            }},
        },
        // Unwind the user array
        {{Key: "$unwind", Value: bson.D{{Key: "path", Value: "$user"}, {Key: "preserveNullAndEmptyArrays", Value: true}}}},
    }

    cursor, err := coll.Aggregate(ctx, pipeline)
    if err != nil {
        return nil, err
    }

    var results []model.JobApplicant
    if err := cursor.All(ctx, &results); err != nil {
        return nil, err
    }

    if len(results) == 0 {
        return nil, mongo.ErrNoDocuments
    }

    return &results[0], nil
}
