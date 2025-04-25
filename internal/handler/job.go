package handler

import (
	"net/http"

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
	var job model.Job
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