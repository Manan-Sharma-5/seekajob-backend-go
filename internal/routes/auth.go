package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/manan-sharma-5/seekajob-backend/internal/handler"
)

func RegisterAuthRoutes(r *gin.Engine) {
    auth := r.Group("/auth")
    {
        auth.POST("/signup", handler.Signup)
    }
}
