package main

import (
	"github.com/gin-gonic/gin"
	"github.com/manan-sharma-5/seekajob-backend/internal/routes"
)

func main() {
    r := gin.Default()
    // Start the server
    r.Use(func(c *gin.Context) {
        origin := c.Request.Header.Get("Origin")
        if origin == "http://localhost:3000" { // Update with allowed origins
            c.Writer.Header().Set("Access-Control-Allow-Origin", origin)
        }
        c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
        c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
        c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
        if c.Request.Method == "OPTIONS" {
            c.AbortWithStatus(204)
            return
        }
        c.Next()
    })

    api := r.Group("/api")

    routes.RegisterAuthRoutes(api)
    routes.ApplyRoutes(api)
    routes.JobRoytes(api)

    r.Use(gin.Logger())
    r.Use(gin.Recovery())
    r.Use(gin.ErrorLogger())
    r.Use(gin.CustomRecovery(func(c *gin.Context, err interface{}) {
        c.JSON(500, gin.H{"error": "Internal Server Error"})
        c.Abort()
    }))
    r.Run()
}
