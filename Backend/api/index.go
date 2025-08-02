package handler

import (
	"net/http"
	"os"
	"shopping-cart-backend/database"
	"shopping-cart-backend/handlers"
	"shopping-cart-backend/middleware"
	
	"github.com/gin-gonic/gin"
)

// Handler is the main entry point for Vercel
func Handler(w http.ResponseWriter, r *http.Request) {
	setupRouter().ServeHTTP(w, r)
}

func setupRouter() *gin.Engine {
	// Connect to database
	database.Connect()

	// Set Gin mode based on environment
	if os.Getenv("GIN_MODE") == "" {
		gin.SetMode(gin.ReleaseMode)
	}

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

	// Public routes
	router.POST("/users", handlers.CreateUser)
	router.GET("/users", handlers.GetUsers)
	router.POST("/users/login", handlers.Login)
	
	router.POST("/items", handlers.CreateItem)
	router.GET("/items", handlers.GetItems)

	// Protected routes (require authentication)
	protected := router.Group("/")
	protected.Use(middleware.AuthMiddleware())
	{
		protected.POST("/carts", handlers.CreateCart)
		protected.GET("/carts", handlers.GetCart)
		protected.DELETE("/carts/:itemId", handlers.RemoveFromCart)
		
		protected.POST("/orders", handlers.CreateOrder)
		protected.GET("/orders", handlers.GetOrders)
	}

	return router
}
