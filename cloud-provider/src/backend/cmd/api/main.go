package main

import (
    "log"
    "net/http"

    "github.com/gin-gonic/gin"
    "cloud-provider/src/backend/internal/config"
    "cloud-provider/src/backend/internal/api/handlers"
)

func main() {
    // Load configuration
    log.Println("Loading configuration...")
    cfg, err := config.Load()
    log.Println("Configuration loaded successfully!")
    if err != nil {
        log.Fatal("Failed to load config:", err)
    }

    // Initialize Gin router
    router := gin.Default()

    // Add CORS middleware if needed
    router.Use(func(c *gin.Context) {
        c.Header("Access-Control-Allow-Origin", "*")
        c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
        c.Header("Access-Control-Allow-Headers", "Content-Type, Authorization")
        
        if c.Request.Method == "OPTIONS" {
            c.AbortWithStatus(204)
            return
        }
        
        c.Next()
    })

    // Health check endpoint
    router.GET("/health", func(c *gin.Context) {
        c.JSON(http.StatusOK, gin.H{
            "status": "healthy",
            "service": "cloud-provider-api",
        })
    })

    // Hello endpoint
    router.GET("/hello", func(c *gin.Context) {
        c.JSON(http.StatusOK, gin.H{
            "message": "Hello from Cloud Provider API",
        })
    })

    // Initialize handlers
    log.Println("About to create ProjectHandler...")

    projectHandler := handlers.NewProjectHandler(cfg)
    // API v1 routes
    v1 := router.Group("/api/v1")
    {
        // Project routes
        v1.GET("/projects", projectHandler.ListProjects)
        v1.GET("/projects/:id", projectHandler.GetProject)
		v1.POST("/projects", projectHandler.CreateProject)
        v1.POST("/networks", projectHandler.CreateNetwork)
        v1.POST("/subnets", projectHandler.CreateSubnet)
        v1.POST("/routers", projectHandler.CreateRouter)
        v1.POST("/router-interfaces", projectHandler.CreateRouterInterface)
        v1.POST("/flavors", projectHandler.CreateFlavor)
    }

    // Start server
    log.Println("Starting Cloud Provider API server on port 8080")
    if err := router.Run(":8080"); err  != nil {
        log.Fatal("Failed to start server:", err)
    }
}