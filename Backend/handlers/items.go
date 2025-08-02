package handlers

import (
	"net/http"
	"shopping-cart-backend/database"
	"shopping-cart-backend/models"
	
	"github.com/gin-gonic/gin"
)

type CreateItemRequest struct {
	Name   string `json:"name" binding:"required"`
	Status string `json:"status"`
}

// CreateItem handles POST /items
func CreateItem(c *gin.Context) {
	var req CreateItemRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if req.Status == "" {
		req.Status = "available"
	}

	item := models.Item{
		Name:   req.Name,
		Status: req.Status,
	}

	if err := database.DB.Create(&item).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create item"})
		return
	}

	c.JSON(http.StatusCreated, item)
}

// GetItems handles GET /items
func GetItems(c *gin.Context) {
	var items []models.Item
	if err := database.DB.Find(&items).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch items"})
		return
	}

	c.JSON(http.StatusOK, items)
}
