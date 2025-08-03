package handler

import (
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

// In-memory storage for serverless - using session-based storage
var (
	users       []map[string]interface{}
	carts       map[int][]map[string]interface{}
	orders      []map[string]interface{}
	nextID      int = 1
	sessionCart []map[string]interface{}
	// Flag to prevent re-initialization
	initialized bool = false
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
			"status":  "running",
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
				"token":   "demo-jwt-token-12345",
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

		// Initialize session cart if it doesn't exist
		if sessionCart == nil {
			sessionCart = []map[string]interface{}{}
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

		// Check if item already in session cart
		found := false
		for i, cartItem := range sessionCart {
			if int(cartItem["item_id"].(int)) == req.ItemID {
				// Update quantity
				sessionCart[i]["quantity"] = sessionCart[i]["quantity"].(int) + 1
				log.Printf("[ADD TO CART] Updated quantity for item %d to %d", req.ItemID, sessionCart[i]["quantity"])
				c.JSON(http.StatusOK, gin.H{
					"message":      "Item quantity updated in cart",
					"item_name":    item["name"],
					"new_quantity": sessionCart[i]["quantity"],
				})
				found = true
				break
			}
		}

		if !found {
			// Add new item to session cart
			cartItem := map[string]interface{}{
				"id":         nextID,
				"item_id":    req.ItemID,
				"quantity":   1,
				"item":       item,
				"created_at": time.Now().Format(time.RFC3339),
			}
			nextID++
			sessionCart = append(sessionCart, cartItem)

			log.Printf("[ADD TO CART] Added new item %d (%s) to cart. Cart now has %d items", req.ItemID, item["name"], len(sessionCart))

			c.JSON(http.StatusCreated, gin.H{
				"message":      "Item added to cart successfully",
				"item_name":    item["name"],
				"cart_item_id": cartItem["id"],
			})
		}
	})

	router.GET("/carts", func(c *gin.Context) {
		// Return the actual session cart with all added items
		if sessionCart == nil {
			sessionCart = []map[string]interface{}{}
		}

		log.Printf("[GET /carts] Returning session cart with %d items", len(sessionCart))

		// Return cart in expected format
		c.JSON(http.StatusOK, gin.H{
			"id":         1,
			"user_id":    1,
			"cart_items": sessionCart,
		})
	})

	// Debug endpoint to clear cart (useful for testing)
	router.POST("/carts/clear", func(c *gin.Context) {
		clearSessionCart()
		c.JSON(http.StatusOK, gin.H{
			"message": "Cart cleared successfully",
		})
	})

	router.DELETE("/carts/:itemId", func(c *gin.Context) {
		itemIDStr := c.Param("itemId")
		itemID, err := strconv.Atoi(itemIDStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid item ID"})
			return
		}

		// Find and remove item from session cart
		if sessionCart == nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Cart is empty"})
			return
		}

		for i, cartItem := range sessionCart {
			if int(cartItem["item_id"].(int)) == itemID {
				// Remove item from session cart
				sessionCart = append(sessionCart[:i], sessionCart[i+1:]...)
				c.JSON(http.StatusOK, gin.H{
					"message": "Item removed from cart successfully",
					"item_id": itemID,
				})
				return
			}
		}

		c.JSON(http.StatusNotFound, gin.H{"error": "Item not found in cart"})
	})

	// Order endpoints
	router.POST("/orders", func(c *gin.Context) {
		// Check if there are items in session cart
		if sessionCart == nil || len(sessionCart) == 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Cart is empty"})
			return
		}

		// Calculate total from session cart
		total := 0
		for _, cartItem := range sessionCart {
			item := cartItem["item"].(map[string]interface{})
			price := int(item["price"].(int))
			quantity := int(cartItem["quantity"].(int))
			total += price * quantity
		}

		// Create order with session cart items in the format expected by frontend
		order := map[string]interface{}{
			"id":      nextID,
			"user_id": 1,
			"total":   total,
			"status":  "completed",
			"cart": map[string]interface{}{
				"cart_items": sessionCart,
			},
			"cart_id":    1,
			"created_at": time.Now().Format(time.RFC3339),
		}
		nextID++

		orders = append(orders, order)

		// Clear session cart after order
		sessionCart = []map[string]interface{}{}

		log.Printf("[CREATE ORDER] Order %d created with total %d. Orders array now has %d orders", order["id"], total, len(orders))

		c.JSON(http.StatusCreated, gin.H{
			"message":  "Order placed successfully",
			"order_id": order["id"],
			"total":    total,
		})
	})

	router.GET("/orders", func(c *gin.Context) {
		// Return actual orders from the orders array
		if orders == nil {
			orders = []map[string]interface{}{}
		}

		log.Printf("[GET /orders] Returning %d actual orders", len(orders))
		
		c.JSON(http.StatusOK, orders)
	})

	return router
}

func initData() {
	// Only initialize once to prevent data loss in serverless
	if initialized {
		log.Println("[INIT] Skipping re-initialization to preserve cart data")
		return
	}

	// Initialize in-memory data
	users = []map[string]interface{}{
		{"id": 1, "username": "testuser", "password": "password"},
	}
	carts = make(map[int][]map[string]interface{})
	orders = []map[string]interface{}{}
	sessionCart = []map[string]interface{}{}
	nextID = 2
	initialized = true
	log.Println("[INIT] Data initialized successfully - cart is empty")
}

// Helper function to clear session cart (useful for testing)
func clearSessionCart() {
	sessionCart = []map[string]interface{}{}
}
