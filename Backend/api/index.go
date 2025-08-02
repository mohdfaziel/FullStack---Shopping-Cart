package handler

import (
	"net/http"
	"os"
	"strconv"
	"time"
	
	"github.com/gin-gonic/gin"
)

// In-memory storage for serverless
var (
	users    []map[string]interface{}
	carts    map[int][]map[string]interface{}
	orders   []map[string]interface{}
	nextID   int = 1
)

// Handler is the main entry point for Vercel
func Handler(w http.ResponseWriter, r *http.Request) {
	setupRouter().ServeHTTP(w, r)
}

func setupRouter() *gin.Engine {
	// Initialize data
	initData()
	
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

	// Items endpoint - Return array directly (not wrapped in object)
	router.GET("/items", func(c *gin.Context) {
		items := []map[string]interface{}{
			{"id": 1, "name": "Laptop", "description": "High-performance laptop", "price": 59999, "stock": 10, "status": "available", "created_at": time.Now().AddDate(0, 0, -30).Format(time.RFC3339)},
			{"id": 2, "name": "Smartphone", "description": "Latest smartphone", "price": 29999, "stock": 15, "status": "available", "created_at": time.Now().AddDate(0, 0, -25).Format(time.RFC3339)},
			{"id": 3, "name": "Headphones", "description": "Wireless headphones", "price": 9999, "stock": 20, "status": "available", "created_at": time.Now().AddDate(0, 0, -20).Format(time.RFC3339)},
			{"id": 4, "name": "Keyboard", "description": "Mechanical keyboard", "price": 7999, "stock": 25, "status": "available", "created_at": time.Now().AddDate(0, 0, -15).Format(time.RFC3339)},
			{"id": 5, "name": "Mouse", "description": "Wireless mouse", "price": 1999, "stock": 30, "status": "available", "created_at": time.Now().AddDate(0, 0, -10).Format(time.RFC3339)},
			{"id": 6, "name": "Monitor", "description": "4K monitor", "price": 24999, "stock": 8, "status": "available", "created_at": time.Now().AddDate(0, 0, -5).Format(time.RFC3339)},
			{"id": 7, "name": "Tablet", "description": "10-inch tablet", "price": 34999, "stock": 12, "status": "available", "created_at": time.Now().AddDate(0, 0, -3).Format(time.RFC3339)},
			{"id": 8, "name": "Webcam", "description": "HD webcam", "price": 4999, "stock": 18, "status": "available", "created_at": time.Now().AddDate(0, 0, -1).Format(time.RFC3339)},
		}
		// Return items array directly (not wrapped in object) to fix frontend map error
		c.JSON(http.StatusOK, items)
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
		
		userID := 1 // Simplified
		
		// Initialize cart if it doesn't exist
		if carts[userID] == nil {
			carts[userID] = []map[string]interface{}{}
		}
		
		// Find item details
		items := []map[string]interface{}{
			{"id": 1, "name": "Laptop", "description": "High-performance laptop", "price": 59999, "stock": 10, "status": "available", "created_at": time.Now().AddDate(0, 0, -30).Format(time.RFC3339)},
			{"id": 2, "name": "Smartphone", "description": "Latest smartphone", "price": 29999, "stock": 15, "status": "available", "created_at": time.Now().AddDate(0, 0, -25).Format(time.RFC3339)},
			{"id": 3, "name": "Headphones", "description": "Wireless headphones", "price": 9999, "stock": 20, "status": "available", "created_at": time.Now().AddDate(0, 0, -20).Format(time.RFC3339)},
			{"id": 4, "name": "Keyboard", "description": "Mechanical keyboard", "price": 7999, "stock": 25, "status": "available", "created_at": time.Now().AddDate(0, 0, -15).Format(time.RFC3339)},
			{"id": 5, "name": "Mouse", "description": "Wireless mouse", "price": 1999, "stock": 30, "status": "available", "created_at": time.Now().AddDate(0, 0, -10).Format(time.RFC3339)},
			{"id": 6, "name": "Monitor", "description": "4K monitor", "price": 24999, "stock": 8, "status": "available", "created_at": time.Now().AddDate(0, 0, -5).Format(time.RFC3339)},
			{"id": 7, "name": "Tablet", "description": "10-inch tablet", "price": 34999, "stock": 12, "status": "available", "created_at": time.Now().AddDate(0, 0, -3).Format(time.RFC3339)},
			{"id": 8, "name": "Webcam", "description": "HD webcam", "price": 4999, "stock": 18, "status": "available", "created_at": time.Now().AddDate(0, 0, -1).Format(time.RFC3339)},
		}
		var item map[string]interface{}
		for _, itm := range items {
			if int(itm["id"].(int)) == req.ItemID {
				item = itm
				break
			}
		}
		
		if item == nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Item not found"})
			return
		}
		
		// Check if item already in cart
		for i, cartItem := range carts[userID] {
			if int(cartItem["item_id"].(int)) == req.ItemID {
				// Update quantity
				carts[userID][i]["quantity"] = carts[userID][i]["quantity"].(int) + 1
				c.JSON(http.StatusOK, gin.H{"message": "Item quantity updated in cart"})
				return
			}
		}
		
		// Add new item to cart
		cartItem := map[string]interface{}{
			"id":         nextID,
			"item_id":    req.ItemID,
			"quantity":   1,
			"item":       item,
			"created_at": time.Now().Format(time.RFC3339),
		}
		nextID++
		carts[userID] = append(carts[userID], cartItem)
		
		c.JSON(http.StatusCreated, gin.H{"message": "Item added to cart successfully"})
	})

	router.GET("/carts", func(c *gin.Context) {
		userID := 1 // Simplified
		
		userCart := carts[userID]
		if userCart == nil {
			userCart = []map[string]interface{}{}
		}
		
		// Return cart in expected format
		c.JSON(http.StatusOK, gin.H{
			"id":         userID,
			"user_id":    userID,
			"cart_items": userCart,
		})
	})

	router.DELETE("/carts/:itemId", func(c *gin.Context) {
		itemIDStr := c.Param("itemId")
		itemID, err := strconv.Atoi(itemIDStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid item ID"})
			return
		}
		
		userID := 1 // Simplified
		
		if carts[userID] == nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Cart not found"})
			return
		}
		
		// Remove item from cart
		for i, cartItem := range carts[userID] {
			if int(cartItem["item_id"].(int)) == itemID {
				carts[userID] = append(carts[userID][:i], carts[userID][i+1:]...)
				c.JSON(http.StatusOK, gin.H{"message": "Item removed from cart successfully"})
				return
			}
		}
		
		c.JSON(http.StatusNotFound, gin.H{"error": "Item not found in cart"})
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
		userID := 1 // Simplified
		
		var userOrders []map[string]interface{}
		for _, order := range orders {
			if int(order["user_id"].(int)) == userID {
				userOrders = append(userOrders, order)
			}
		}
		
		c.JSON(http.StatusOK, userOrders)
	})

	return router
}

func initData() {
	// Initialize in-memory data
	users = []map[string]interface{}{
		{"id": 1, "username": "testuser", "password": "password"},
	}
	carts = make(map[int][]map[string]interface{})
	orders = []map[string]interface{}{}
	nextID = 2
}
