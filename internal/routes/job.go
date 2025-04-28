package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/manan-sharma-5/seekajob-backend/internal/handler"
	"github.com/manan-sharma-5/seekajob-backend/internal/middleware"
)

func JobRoytes(rg *gin.RouterGroup) {
    job := rg.Group("/job")
    {
		job.GET("/applicant/:id", middleware.AuthMiddleware(), handler.GetApplicantDetails)
		job.POST("", middleware.AuthMiddleware(),handler.CreateJob)
		job.GET("/:id", middleware.AuthMiddleware(), handler.GetJobByID)
		job.GET("/", middleware.AuthMiddleware(), handler.GetAllJobs)
		job.GET("/applicants/:id", middleware.AuthMiddleware(), handler.GetJobApplicants)
		job.GET("/recruiter-jobs", middleware.AuthMiddleware(), handler.GetJobByRecruiterID)
		job.DELETE("/:id", middleware.AuthMiddleware(), handler.DeleteJob)
		job.GET("/search", middleware.AuthMiddleware(), handler.SearchJob)
    }
}
