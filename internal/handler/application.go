package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/manan-sharma-5/seekajob-backend/internal/model"
	"github.com/manan-sharma-5/seekajob-backend/internal/repository"
	service "github.com/manan-sharma-5/seekajob-backend/internal/servivce"
)

func ApplyJob(c *gin.Context) {
    var req model.ApplyJobRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }
	userID, err := c.Cookie("user_id")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get user ID"})
		return
	}

	// Set the user ID in the request
	req.UserID = userID

	// Call the service to apply for the job

	err = service.ApplyJobService(&req)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Job application successful"})
}

func GetJobApplicants (c *gin.Context) {
	jobID := c.Param("id")
	applicants, err := repository.GetApplicationsWithJobID(jobID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Job not found"})
		return
	}
	c.JSON(http.StatusOK, applicants)
}