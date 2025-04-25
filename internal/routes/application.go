package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/manan-sharma-5/seekajob-backend/internal/handler"
)

func ApplyRoutes(r *gin.Engine) {
    auth := r.Group("/apply")
    {
        auth.POST("/job", handler.ApplyJob)
    }
}
