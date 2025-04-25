package main

import (
	"github.com/gin-gonic/gin"
	"github.com/manan-sharma-5/seekajob-backend/internal/routes"
)

func main() {
    r := gin.Default()
    routes.RegisterAuthRoutes(r)
    r.Run()
}
