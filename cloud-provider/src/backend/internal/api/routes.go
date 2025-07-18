package api

import (
	"github.com/gin-gonic/gin"
    "cloud-provider/backend/src/internal/api/handlers"
    "cloud-provider/backend/src/internal/config"
)

func SetupRoutes(router *gin.Engine, cfg *config.Config) {
	// Create handlers
	projectHandler := handlers.NewProjectHandler(cfg) 

	// API v1 routes
	v1 := router.Group("/api/v1")
	{
		projects := v1.Group("/projects")
        {
            projects.GET("", projectHandler.ListProjects)
            projects.GET("/:id", projectHandler.GetProject)
        }

	}
}
