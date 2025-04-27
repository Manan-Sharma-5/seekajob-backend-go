package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/manan-sharma-5/seekajob-backend/internal/handler"
)

func RegisterAuthRoutes(rg *gin.RouterGroup) {
    auth := rg.Group("/auth")
    {
        auth.POST("/signup", handler.Signup)
        auth.POST("/login", handler.Login)
    }
}
