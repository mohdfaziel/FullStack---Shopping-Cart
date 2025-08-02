package handlers

import (
	"net/http"
	"shopping-cart-backend/database"
	"shopping-cart-backend/models"
	
	"github.com/gin-gonic/gin"
)

type AddItemToCartRequest struct {
	ItemID uint `json:"item_id" binding:"required"`
}

// CreateCart handles POST /carts (add item to cart)
func CreateCart(c *gin.Context) {
	userID := c.GetUint("user_id")
	
	var req AddItemToCartRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Check if item exists and is available
	var item models.Item
	if err := database.DB.First(&item, req.ItemID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Item not found"})
		return
	}

	if item.Status != "available" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Item not available"})
		return
	}

	// Get or create user's active cart
	var cart models.Cart
	err := database.DB.Where("user_id = ? AND status = ?", userID, "active").First(&cart).Error
	
	if err != nil {
		// Create new cart
		cart = models.Cart{
			UserID: userID,
			Name:   "Shopping Cart",
			Status: "active",
		}
		if err := database.DB.Create(&cart).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create cart"})
			return
		}

		// Update user's cart_id
		database.DB.Model(&models.User{}).Where("id = ?", userID).Update("cart_id", cart.ID)
	}

	// Check if item already in cart
	var existingCartItem models.CartItem
	if err := database.DB.Where("cart_id = ? AND item_id = ?", cart.ID, req.ItemID).First(&existingCartItem).Error; err == nil {
		// Item exists, increase quantity
		existingCartItem.Quantity++
		if err := database.DB.Save(&existingCartItem).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update item quantity"})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"message": "Item quantity updated in cart",
			"cart_id": cart.ID,
			"item_id": req.ItemID,
			"quantity": existingCartItem.Quantity,
		})
		return
	}

	// Add new item to cart
	cartItem := models.CartItem{
		CartID:   cart.ID,
		ItemID:   req.ItemID,
		Quantity: 1,
	}

	if err := database.DB.Create(&cartItem).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add item to cart"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Item added to cart successfully",
		"cart_id": cart.ID,
		"item_id": req.ItemID,
	})
}

// GetCart handles GET /carts
func GetCart(c *gin.Context) {
	userID := c.GetUint("user_id")

	var cart models.Cart
	if err := database.DB.Where("user_id = ? AND status = ?", userID, "active").
		Preload("Items").
		Preload("CartItems.Item").First(&cart).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "No active cart found"})
		return
	}

	c.JSON(http.StatusOK, cart)
}

// RemoveFromCart handles DELETE /carts/:itemId (remove item from cart)
func RemoveFromCart(c *gin.Context) {
	userID := c.GetUint("user_id")
	itemID := c.Param("itemId")

	// Get user's active cart
	var cart models.Cart
	if err := database.DB.Where("user_id = ? AND status = ?", userID, "active").First(&cart).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "No active cart found"})
		return
	}

	// Find and remove the cart item
	var cartItem models.CartItem
	if err := database.DB.Where("cart_id = ? AND item_id = ?", cart.ID, itemID).First(&cartItem).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Item not found in cart"})
		return
	}

	// Delete the cart item
	if err := database.DB.Delete(&cartItem).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to remove item from cart"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Item removed from cart successfully",
		"cart_id": cart.ID,
		"item_id": itemID,
	})
}
