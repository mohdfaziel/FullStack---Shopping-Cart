package handlers

import (
	"net/http"
	"shopping-cart-backend/database"
	"shopping-cart-backend/models"
	
	"github.com/gin-gonic/gin"
)

// CreateOrder handles POST /orders (checkout)
func CreateOrder(c *gin.Context) {
	userID := c.GetUint("user_id")

	// Get user's active cart
	var cart models.Cart
	if err := database.DB.Where("user_id = ? AND status = ?", userID, "active").First(&cart).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "No active cart found"})
		return
	}

	// Check if cart has items
	var cartItemCount int64
	database.DB.Model(&models.CartItem{}).Where("cart_id = ?", cart.ID).Count(&cartItemCount)
	if cartItemCount == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Cart is empty"})
		return
	}

	// Create order
	order := models.Order{
		CartID: cart.ID,
		UserID: userID,
	}

	if err := database.DB.Create(&order).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create order"})
		return
	}

	// Update cart status to ordered
	cart.Status = "ordered"
	if err := database.DB.Save(&cart).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update cart status"})
		return
	}

	// Clear user's cart_id (so they can create a new cart)
	database.DB.Model(&models.User{}).Where("id = ?", userID).Update("cart_id", nil)

	c.JSON(http.StatusCreated, gin.H{
		"message":  "Order created successfully",
		"order_id": order.ID,
		"cart_id":  cart.ID,
	})
}

// GetOrders handles GET /orders
func GetOrders(c *gin.Context) {
	userID := c.GetUint("user_id")

	var orders []models.Order
	if err := database.DB.Where("user_id = ?", userID).
		Preload("Cart").
		Preload("Cart.CartItems.Item").Find(&orders).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch orders"})
		return
	}

	c.JSON(http.StatusOK, orders)
}
