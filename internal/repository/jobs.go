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
	"go.mongodb.org/mongo-driver/mongo/options"
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
        {{Key: "$match", Value: bson.D{{Key: "_id", Value: applicantIDObj}}}},
        
        {
            {Key: "$lookup", Value: bson.D{
                {Key: "from", Value: "jobs"},
                {Key: "localField", Value: "jobID"},
                {Key: "foreignField", Value: "_id"},
                {Key: "as", Value: "job"},
            }},
        },
        {{Key: "$unwind", Value: bson.D{{Key: "path", Value: "$job"}, {Key: "preserveNullAndEmptyArrays", Value: true}}}},

        {
            {Key: "$lookup", Value: bson.D{
                {Key: "from", Value: "users"},
                {Key: "localField", Value: "userID"},
                {Key: "foreignField", Value: "_id"},
                {Key: "as", Value: "user"},
            }},
        },
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

func GetJobByID(jobID string, userID string) (*model.Job, bool, []model.Job, error) {
    client, err := db.GetMongoClient()
    if err != nil {
        return nil, false, nil, err
    }

    coll := client.Collection("jobs")
    applicantColl := client.Collection("job_applications")
    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()

    var job model.Job
    objobID, err := primitive.ObjectIDFromHex(jobID)
    if err != nil {
        return nil, false, nil, err
    }
    err = coll.FindOne(ctx, bson.M{"_id": objobID}).Decode(&job)
    if err != nil {
        return nil, false, nil, err
    }

    obuserID, err := primitive.ObjectIDFromHex(userID)
    if err != nil {
        return nil, false, nil, err
    }

    // Check if the user has applied for this job
    jobApplicant := true
    jobApplier := model.JobApplicant{}
    err = applicantColl.FindOne(ctx, bson.M{"jobID": objobID, "userID": obuserID}).Decode(&jobApplier)
    if err != nil {
        if err == mongo.ErrNoDocuments {
            // User has not applied for this job
            jobApplicant = false
        } else {
            // Some other error occurred
            return nil, false, nil, err
        }
    }

    filter := bson.M{
        "_id": bson.M{"$ne": objobID}, // Exclude current job
        "$or": []bson.M{
            {"companyID": job.CompanyID},
            {"tags": bson.M{"$in": job.Tags}},
        },
    }

    options := options.Find()
    options.SetLimit(3)
    options.SetSort(bson.M{"createdAt": -1})

    cursor, err := coll.Find(ctx, filter, options)
    if err != nil {
        return &job, jobApplicant, nil, err
    }

    var recommendedJobs []model.Job
    if err := cursor.All(ctx, &recommendedJobs); err != nil {
        return &job, jobApplicant, nil, err
    }

    return &job, jobApplicant, recommendedJobs, nil
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

func DeleteJob(jobID string) error {
    dbClient, err := db.GetMongoClient()
    if err != nil {
        return err
    }

    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()

    jobsCollection := dbClient.Collection("jobs")
    jobApplicationsCollection := dbClient.Collection("job_applications")

    oid, err := primitive.ObjectIDFromHex(jobID)
    if err != nil {
        return err
    }

    _, err = jobsCollection.DeleteOne(ctx, bson.M{"_id": oid})
    if err != nil {
        return err
    }

    _, err = jobApplicationsCollection.DeleteMany(ctx, bson.M{"jobID": oid})
    if err != nil {
        return err
    }

    return nil
}

func SearchJob(jobName string, salary float64, location string, tags string) ([]model.Job, error) {
    client, err := db.GetMongoClient()
    if err != nil {
        return nil, err
    }

    coll := client.Collection("jobs")
    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()

    filter := bson.M{}
    if jobName != "" {
        filter["title"] = bson.M{"$regex": jobName, "$options": "i"}
    }
    if salary != 0 {
        filter["salary"] = bson.M{"$gte": salary}
    }
    if location != "" {
        filter["location"] = bson.M{"$regex": location, "$options": "i"}
    }
    if tags != "" {
        filter["tags"] = bson.M{"$regex": tags, "$options": "i"}
    }
    options := options.Find()
    options.SetLimit(10)
    options.SetSort(bson.M{"createdAt": -1})
    options.SetProjection(bson.M{"_id": 1, "title": 1, "location": 1, "salary": 1, "tags": 1})
    options.SetSkip(0)


    cursor, err := coll.Find(ctx, filter)
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