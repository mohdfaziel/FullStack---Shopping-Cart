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
		allowedOrigins = "https://full-stack-shopping-cart.vercel.app"
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

	// Authentication endpoints
	router.POST("/users", func(c *gin.Context) {
		var req struct {
			Username string `json:"username"`
			Password string `json:"password"`
		}
		
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		
		// Simple validation
		if req.Username == "" || req.Password == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Username and password are required"})
			return
		}
		
		// In a real app, you'd hash the password and save to database
		// For demo purposes, just return success
		c.JSON(http.StatusCreated, gin.H{
			"message": "User created successfully",
			"user_id": 1,
		})
	})

	router.POST("/users/login", func(c *gin.Context) {
		var req struct {
			Username string `json:"username"`
			Password string `json:"password"`
		}
		
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		
		// Simple demo authentication - accept any non-empty username/password
		if req.Username != "" && req.Password != "" {
			c.JSON(http.StatusOK, gin.H{
				"token": "demo-jwt-token-12345",
				"user_id": 1,
				"message": "Login successful",
			})
		} else {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid username or password"})
		}
	})

	router.GET("/users", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"users": []map[string]interface{}{
				{"id": 1, "username": "demo-user"},
			},
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

	// Cart endpoints
	router.POST("/carts", func(c *gin.Context) {
		var req struct {
			ItemID   int `json:"item_id"`
			Quantity int `json:"quantity"`
		}
		
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		
		c.JSON(http.StatusCreated, gin.H{
			"message": "Item added to cart successfully",
			"cart_item_id": req.ItemID,
		})
	})

	router.GET("/carts", func(c *gin.Context) {
		// Return sample cart data
		cartItems := []map[string]interface{}{
			{
				"id": 1,
				"user_id": 1,
				"item_id": 1,
				"quantity": 2,
				"item": map[string]interface{}{
					"id": 1,
					"name": "Laptop",
					"description": "High-performance laptop",
					"price": 59999,
					"stock": 10,
				},
			},
		}
		
		c.JSON(http.StatusOK, gin.H{
			"id": 1,
			"user_id": 1,
			"cart_items": cartItems,
		})
	})

	router.DELETE("/carts/:itemId", func(c *gin.Context) {
		itemId := c.Param("itemId")
		c.JSON(http.StatusOK, gin.H{
			"message": "Item removed from cart successfully",
			"item_id": itemId,
		})
	})

	// Order endpoints
	router.POST("/orders", func(c *gin.Context) {
		c.JSON(http.StatusCreated, gin.H{
			"message": "Order placed successfully",
			"order_id": 1,
			"total": 119998,
		})
	})

	router.GET("/orders", func(c *gin.Context) {
		orders := []map[string]interface{}{
			{
				"id": 1,
				"user_id": 1,
				"total": 119998,
				"created_at": "2025-08-03T00:00:00Z",
				"items": []map[string]interface{}{
					{
						"id": 1,
						"item_id": 1,
						"quantity": 2,
						"item": map[string]interface{}{
							"id": 1,
							"name": "Laptop",
							"price": 59999,
						},
					},
				},
			},
		}
		
		c.JSON(http.StatusOK, gin.H{"orders": orders})
	})

	return router
}
