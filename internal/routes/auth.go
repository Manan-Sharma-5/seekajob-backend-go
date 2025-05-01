package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/manan-sharma-5/seekajob-backend/internal/handler"
	"github.com/manan-sharma-5/seekajob-backend/internal/middleware"
)

func RegisterAuthRoutes(rg *gin.RouterGroup) {
    auth := rg.Group("/auth")
    {
        auth.POST("/signup", handler.Signup)
        auth.POST("/login", handler.Login)
        auth.POST("/logout", middleware.AuthMiddleware(), handler.Logout)
        auth.GET("/user", middleware.AuthMiddleware(), handler.GetUser)
        auth.POST("/upload-resume", middleware.AuthMiddleware(), handler.UploadResume)
    }
}
