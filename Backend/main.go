package main

import (
	"os"
	"shopping-cart-backend/database"
	"shopping-cart-backend/handlers"
	"shopping-cart-backend/middleware"
	
	"github.com/gin-gonic/gin"
)

func main() {
	// Connect to database
	database.Connect()

	// Set Gin mode based on environment
	if os.Getenv("GIN_MODE") == "" {
		gin.SetMode(gin.ReleaseMode)
	}

	// Create Gin router
	r := gin.Default()

	// CORS middleware with environment-based origin
	allowedOrigins := os.Getenv("ALLOWED_ORIGINS")
	if allowedOrigins == "" {
		allowedOrigins = "*"
	}

	r.Use(func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", allowedOrigins)
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Origin, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
		
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		
		c.Next()
	})

	// Public routes
	r.POST("/users", handlers.CreateUser)
	r.GET("/users", handlers.GetUsers)
	r.POST("/users/login", handlers.Login)
	
	r.POST("/items", handlers.CreateItem)
	r.GET("/items", handlers.GetItems)

	// Protected routes (require authentication)
	protected := r.Group("/")
	protected.Use(middleware.AuthMiddleware())
	{
		protected.POST("/carts", handlers.CreateCart)
		protected.GET("/carts", handlers.GetCart)
		protected.DELETE("/carts/:itemId", handlers.RemoveFromCart)
		
		protected.POST("/orders", handlers.CreateOrder)
		protected.GET("/orders", handlers.GetOrders)
	}

	// Start server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8081"
	}
	r.Run(":" + port)
}
