package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Job struct {
    ID                  string       `bson:"_id,omitempty" json:"id"`
    Title               string       `bson:"title" json:"title"`
    Description         string       `bson:"description" json:"description"`
    DetailedDescription string       `bson:"detailedDescription" json:"detailedDescription"`
    CompanyID           string       `bson:"companyID" json:"companyID"`
    Company             *Company     `bson:"company,omitempty" json:"company,omitempty"`
    Experience          string       `bson:"experience" json:"experience"`
    Location            string       `bson:"location" json:"location"`
    Salary              float64      `bson:"salary" json:"salary"`
    RecruiterID         string       `bson:"recruiterId" json:"recruiterId"`
    Recruiter           *User   `bson:"recruiter,omitempty" json:"recruiter,omitempty"`
    Category            string       `bson:"category" json:"category"`
    Candidates          []string     `bson:"candidates" json:"candidates"`
    CandidatesApplied   []*User `bson:"candidatesApplied,omitempty" json:"candidatesApplied,omitempty"`
    Tags                []string     `bson:"tags" json:"tags"`
    CreatedAt           time.Time    `bson:"createdAt" json:"createdAt"`
    UpdatedAt           time.Time    `bson:"updatedAt" json:"updatedAt"`
    JobApplicant        []JobApplicant `bson:"jobApplicant,omitempty" json:"jobApplicant,omitempty"`
}

type Company struct {
    ID        string    `bson:"_id,omitempty" json:"id"`
    Name      string    `bson:"name" json:"name"`
    CreatedAt time.Time `bson:"createdAt" json:"createdAt"`
    UpdatedAt time.Time `bson:"updatedAt" json:"updatedAt"`
    Jobs      []Job     `bson:"jobs,omitempty" json:"jobs,omitempty"`
}

type JobApplicant struct {
    ID        primitive.ObjectID `bson:"_id,omitempty" json:"id"`
    JobID     primitive.ObjectID `bson:"jobID" json:"jobID"`
    UserID    primitive.ObjectID `bson:"userID" json:"userID"`

    Job       *Job               `bson:"job,omitempty" json:"job,omitempty"`
    User      *User         `bson:"user,omitempty" json:"user,omitempty"`

    Name      string             `bson:"name" json:"name"`
    Email     string             `bson:"email" json:"email"`
    Resume    string             `bson:"resume" json:"resume"`
    Status    string             `bson:"status" json:"status"`
    CreatedAt time.Time          `bson:"createdAt" json:"createdAt"`
    UpdatedAt time.Time          `bson:"updatedAt" json:"updatedAt"`
}


type CreateJobRequest struct {
	Title               string   `json:"title" binding:"required"`
	Description         string   `json:"description" binding:"required"`
	DetailedDescription string   `json:"detailedDescription" binding:"required"`
	Experience          string   `json:"experience" binding:"required"`
	Location            string   `json:"location" binding:"required"`
	Salary              float64  `json:"salary" binding:"required"`
	Category            string   `json:"category" binding:"required"`
	Tags                []string `json:"tags" binding:"required"`
	Company             string   `json:"company" binding:"required"`
}

type ApplyJobRequest struct {
	JobID     string `json:"job_id" binding:"required"`
	UserID    string `json:"user_id"`
}

type UserWithApplicationsAndJobs struct {
    ID             string          `bson:"_id,omitempty" json:"id"`
    Name           string          `bson:"name" json:"name"`
    Email          string          `bson:"email" json:"email"`
    IsCandidate    bool            `bson:"isCandidate" json:"isCandidate"`
    CreatedAt      time.Time       `bson:"createdAt" json:"createdAt"`
    UpdatedAt      time.Time       `bson:"updatedAt" json:"updatedAt"`
    AppliedJobs    []Job     `bson:"appliedJobs" json:"appliedJobs"`
    JobApplications []JobApplicant `bson:"jobApplications" json:"jobApplications"`
}
