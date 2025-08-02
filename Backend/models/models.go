package models

import (
	"time"
)

// User model
type User struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	Username  string    `json:"username" gorm:"unique;not null"`
	Password  string    `json:"-" gorm:"not null"` // "-" prevents password from being serialized
	Token     string    `json:"token"`
	CartID    *uint     `json:"cart_id"`
	CreatedAt time.Time `json:"created_at"`
	
	// Relationships
	Cart   *Cart   `json:"cart,omitempty" gorm:"foreignKey:CartID"`
	Orders []Order `json:"orders,omitempty" gorm:"foreignKey:UserID"`
}

// Item model
type Item struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	Name      string    `json:"name" gorm:"not null"`
	Status    string    `json:"status" gorm:"default:'available'"` // available, unavailable
	CreatedAt time.Time `json:"created_at"`
	
	// Relationships
	CartItems []CartItem `json:"cart_items,omitempty" gorm:"foreignKey:ItemID"`
}

// Cart model
type Cart struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	UserID    uint      `json:"user_id" gorm:"not null"`
	Name      string    `json:"name"`
	Status    string    `json:"status" gorm:"default:'active'"` // active, ordered
	CreatedAt time.Time `json:"created_at"`
	
	// Relationships
	User      User       `json:"user,omitempty" gorm:"foreignKey:UserID"`
	Items     []Item     `json:"items,omitempty" gorm:"many2many:cart_items;"`
	CartItems []CartItem `json:"cart_items,omitempty" gorm:"foreignKey:CartID"`
}

// CartItem model (join table)
type CartItem struct {
	CartID   uint `json:"cart_id" gorm:"primaryKey"`
	ItemID   uint `json:"item_id" gorm:"primaryKey"`
	Quantity int  `json:"quantity" gorm:"default:1"`
	
	// Relationships
	Cart Cart `json:"cart,omitempty" gorm:"foreignKey:CartID"`
	Item Item `json:"item,omitempty" gorm:"foreignKey:ItemID"`
}

// Order model
type Order struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	CartID    uint      `json:"cart_id" gorm:"not null"`
	UserID    uint      `json:"user_id" gorm:"not null"`
	CreatedAt time.Time `json:"created_at"`
	
	// Relationships
	Cart Cart `json:"cart,omitempty" gorm:"foreignKey:CartID"`
	User User `json:"user,omitempty" gorm:"foreignKey:UserID"`
}
