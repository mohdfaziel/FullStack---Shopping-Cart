package handler

import (
	"net/http"
	"os"
	
	"github.com/gin-gonic/gin"
)

// Handler is the main entry point for Vercel
func Handler(w http.ResponseWriter, r *http.Request) {
	setupRouter().ServeHTTP(w, r)
}

func setupRouter() *gin.Engine {
	// Set Gin mode based on environment
	gin.SetMode(gin.ReleaseMode)

	// Create Gin router
	router := gin.Default()

	// CORS middleware with environment-based origin
	allowedOrigins := os.Getenv("ALLOWED_ORIGINS")
	if allowedOrigins == "" {
		allowedOrigins = "*"
	}

	router.Use(func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", allowedOrigins)
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Origin, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
		
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		
		c.Next()
	})

	// Test route
	router.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "Shopping Cart Backend API",
			"status": "running",
			"version": "1.0.0",
		})
	})

	// Simple items endpoint
	router.GET("/items", func(c *gin.Context) {
		items := []map[string]interface{}{
			{"id": 1, "name": "Laptop", "description": "High-performance laptop", "price": 59999, "stock": 10},
			{"id": 2, "name": "Smartphone", "description": "Latest smartphone", "price": 29999, "stock": 15},
			{"id": 3, "name": "Headphones", "description": "Wireless headphones", "price": 9999, "stock": 20},
			{"id": 4, "name": "Keyboard", "description": "Mechanical keyboard", "price": 7999, "stock": 25},
			{"id": 5, "name": "Mouse", "description": "Wireless mouse", "price": 1999, "stock": 30},
			{"id": 6, "name": "Monitor", "description": "4K monitor", "price": 24999, "stock": 8},
			{"id": 7, "name": "Tablet", "description": "10-inch tablet", "price": 34999, "stock": 12},
			{"id": 8, "name": "Webcam", "description": "HD webcam", "price": 4999, "stock": 18},
		}
		c.JSON(http.StatusOK, gin.H{"items": items})
	})

	return router
}
