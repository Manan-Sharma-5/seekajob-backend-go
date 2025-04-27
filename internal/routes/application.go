package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/manan-sharma-5/seekajob-backend/internal/handler"
	"github.com/manan-sharma-5/seekajob-backend/internal/middleware"
)

func ApplyRoutes(rg *gin.RouterGroup) {
    auth := rg.Group("/apply")
    {
        auth.POST("/job", middleware.AuthMiddleware(), handler.ApplyJob)
    }
}
