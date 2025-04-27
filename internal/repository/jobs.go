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
)

func CreateJob(jobRequest *model.CreateJobRequest, userID string) error {
	client, err := db.GetMongoClient()
	coll := client.Collection("jobs")
    compColl := client.Collection("companies")
	if err != nil {
		return err
	}

    companyName := jobRequest.Company

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
            companyCreated, err := compColl.InsertOne(context.Background(), company)
            company.ID = companyCreated.InsertedID.(primitive.ObjectID).Hex()
            if err != nil {
                return err
            }
        } else {
            return err
        }
    }

    var job model.Job

    // Set the company ID in the job
    job.Title = jobRequest.Title
    job.Description = jobRequest.Description
    job.DetailedDescription = jobRequest.DetailedDescription
    job.Experience = jobRequest.Experience
    job.Location = jobRequest.Location
    job.Salary = jobRequest.Salary
    job.Category = jobRequest.Category
    job.Tags = jobRequest.Tags
    job.Company = &company
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
    client, err := db.GetMongoClient()
    coll := client.Collection("job_applications")
    if err != nil {
        return nil, err
    }

    applicantIDObj, err := primitive.ObjectIDFromHex(applicantID)
    if err != nil {
        return nil, err
    }


    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()

    pipeline := mongo.Pipeline{
        // Match the job applicant
        {{Key: "$match", Value: bson.D{{Key: "_id", Value: applicantIDObj}}}},
        
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

    log.Println("Job Applicant with relations:", results)

    return &results[0], nil
}

func GetJobByID(jobID string) (*model.Job, error) {
    client, err := db.GetMongoClient()
    if err != nil {
        return nil, err
    }

    coll := client.Collection("jobs")
    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()

    var job model.Job
    objobID, err := primitive.ObjectIDFromHex(jobID)
    if err != nil {
        return nil, err
    }
    err = coll.FindOne(ctx, bson.M{"_id": objobID}).Decode(&job)
    if err != nil {
        return nil, err
    }

    return &job, nil
}

func GetApplicationsWithJobID(jobID string) ([]model.JobApplicant, error) {
    client, err := db.GetMongoClient()
    if err != nil {
        return nil, err
    }

    coll := client.Collection("job_applications")
    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()

    objobID, err := primitive.ObjectIDFromHex(jobID)
    if err != nil {
        return nil, err
    }

    cursor, err := coll.Find(ctx, bson.M{"jobID": objobID})
    if err != nil {
        return nil, err
    }
    defer cursor.Close(ctx)

    var applications []model.JobApplicant
    if err := cursor.All(ctx, &applications); err != nil {
        return nil, err
    }

    return applications, nil
}

func GetJobByRecruiterID(recruiterID string) ([]model.Job, error) {
    client, err := db.GetMongoClient()
    if err != nil {
        return nil, err
    }

    log.Println("Recruiter ID:", recruiterID)

    coll := client.Collection("jobs")
    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()
    if err != nil {
        return nil, err
    }

    cursor, err := coll.Find(ctx, bson.M{"recruiterId": recruiterID})
    if err != nil {
        return nil, err
    }
    defer cursor.Close(ctx)

    log.Println("Cursor:", cursor)

    var jobs []model.Job
    if err := cursor.All(ctx, &jobs); err != nil {
        return nil, err
    }

    return jobs, nil
}