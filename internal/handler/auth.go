package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/manan-sharma-5/seekajob-backend/internal/model"
	service "github.com/manan-sharma-5/seekajob-backend/internal/servivce"
)

func Signup(c *gin.Context) {
    var req model.SignupRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    userID, err := service.SignupUser(req)
    if err != nil {
        if err == service.ErrEmailAlreadyInUse {
            c.JSON(http.StatusConflict, gin.H{"error": "Email already in use"})
        } else {
            c.JSON(http.StatusInternalServerError, gin.H{"error": "Signup failed"})
        }
        return
    }

    c.SetCookie("user_id", userID, 3600, "/", "", false, true)
    c.SetCookie("is_candidate", "true", 3600, "/", "", false, true)
    c.JSON(http.StatusOK, gin.H{"message": "Signup successful"})
}

func Login(c *gin.Context) {
    var req model.LoginRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    userID, err := service.LoginUser(req)
    if err != nil {
        if err == service.ErrInvalidCredentials {
            c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
        } else {
            c.JSON(http.StatusInternalServerError, gin.H{"error": "Login failed"})
        }
        return
    }
    c.SetCookie("user_id", userID, 3600, "/", "", false, true)
    c.JSON(http.StatusOK, gin.H{"message": "Login successful"})
}