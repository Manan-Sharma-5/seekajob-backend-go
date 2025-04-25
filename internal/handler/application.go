package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/manan-sharma-5/seekajob-backend/internal/model"
	service "github.com/manan-sharma-5/seekajob-backend/internal/servivce"
)

func ApplyJob(c *gin.Context) {
    var req model.ApplyJobRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

	// Call the service to apply for the job

	err := service.ApplyJobService(&req)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Job application successful"})
}
