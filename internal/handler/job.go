package handler

import (
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/manan-sharma-5/seekajob-backend/internal/model"
	"github.com/manan-sharma-5/seekajob-backend/internal/repository"
)

func GetApplicantDetails(c *gin.Context) {
    id := c.Param("id")
    
    applicant, err := repository.GetJobApplicantWithRelations(id)
    if err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "Applicant not found"})
        return
    }

    c.JSON(http.StatusOK, applicant)
}

func CreateJob(c *gin.Context){
	var job model.CreateJobRequest
	if err := c.ShouldBindJSON(&job); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// fetch user id from cookie

	userID, err := c.Cookie("user_id")

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create job"})
		return
	}

	err = repository.CreateJob(&job, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create job"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Job created successfully"})
}

func GetJobByID(c *gin.Context) {
	jobID := c.Param("id")
	userID, err := c.Cookie("user_id")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get user ID"})
		return
	}
	jobs, hasApplied, recomjob, err := repository.GetJobByID(jobID, userID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK,
		gin.H{
			"jobs":         jobs,
			"hasApplied":  hasApplied,
			"recomjob": recomjob,
		},
	)
}

func GetAllJobs(c *gin.Context) {
	jobs, err := repository.GetAllJobs()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch jobs"})
		return
	}
	c.JSON(http.StatusOK, jobs)
}

func GetJobByRecruiterID(c *gin.Context) {
	recruiterID , err := c.Cookie("user_id")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get user ID"})
		return
	}
	
	jobs, err := repository.GetJobByRecruiterID(recruiterID)
	if err != nil {
		log.Println("Error fetching jobs by recruiter ID:", err)
		c.JSON(http.StatusNotFound, gin.H{"error": "Job not found"})
		return
	}
	c.JSON(http.StatusOK, jobs)
}

func DeleteJob(c *gin.Context) {
	jobID := c.Param("id")
	err := repository.DeleteJob(jobID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete job"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Job deleted successfully"})
}

func SearchJob(c *gin.Context) {
	query := c.Query("query")
	salary := c.Query("salary")
	location := c.Query("location")
	tags := c.Query("tags")
	if query == "" && salary == "" && location == "" && tags == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "At least one search parameter is required"})
		return
	}

	floatSalary := 0.0
	if salary != "" {
		var err error
		floatSalary, err = strconv.ParseFloat(salary, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid salary format"})
			return
		}
	}

	jobs, err := repository.SearchJob(query, floatSalary, location, tags)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to search jobs"})
		return
	}
	c.JSON(http.StatusOK, jobs)
}