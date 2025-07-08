package postgres

import (
	"time"
)

// User represents a customer or admin user.
type User struct {
	ID        uint   `gorm:"primaryKey"`
	Email     string `gorm:"unique;not null"`
	Password  string `gorm:"not null"`
	Name      string
	CreatedAt time.Time
	UpdatedAt time.Time
	Addresses []Address
	Orders    []Order
	Wishlist  []Product `gorm:"many2many:user_wishlist"`
	Ratings   []Rating
}

// Address represents a user's shipping or billing address.
type Address struct {
	ID        uint   `gorm:"primaryKey"`
	UserID    uint   `gorm:"not null;index"`
	Line1     string `gorm:"not null"`
	Line2     string
	City      string `gorm:"not null"`
	State     string `gorm:"not null"`
	ZipCode   string `gorm:"not null"`
	Country   string `gorm:"not null"`
	IsDefault bool   `gorm:"default:false"`
}

// Category represents a product category.
type Category struct {
	ID       uint   `gorm:"primaryKey"`
	Name     string `gorm:"unique;not null"`
	Products []Product
}

// Product represents an item for sale.
type Product struct {
	ID          uint   `gorm:"primaryKey"`
	Name        string `gorm:"not null"`
	Description string
	Price       float64 `gorm:"not null"`
	CategoryID  uint    `gorm:"index"`
	Category    Category
	Images      []ProductImage
	Ratings     []Rating
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type Inventory struct {
	ItemID uint   `gorm:"primaryKey"`
	Stock  uint64 `gorm:"not null"` // Number of items in stock
}

// ProductImage represents an image for a product.
type ProductImage struct {
	ID        uint   `gorm:"primaryKey"`
	ProductID uint   `gorm:"index;not null"`
	URL       string `gorm:"not null"`
	AltText   string
}

// Rating represents a user's rating and review for a product.
type Rating struct {
	ID        uint `gorm:"primaryKey"`
	UserID    uint `gorm:"index;not null"`
	ProductID uint `gorm:"index;not null"`
	Score     int  `gorm:"not null"` // 1-5
	Comment   string
	CreatedAt time.Time
}

// Order represents a customer order.
type Order struct {
	ID            uint    `gorm:"primaryKey"`
	UserID        uint    `gorm:"index;not null"`
	AddressID     uint    `gorm:"not null"`
	Status        string  `gorm:"not null"` // e.g., "processing", "shipped", "delivered"
	TotalAmount   float64 `gorm:"not null"`
	PaymentMethod string  `gorm:"not null"` // "stripe", "paypal"
	OrderItems    []OrderItem
	CreatedAt     time.Time
	UpdatedAt     time.Time
}

// OrderItem represents a product within an order.
type OrderItem struct {
	ID        uint    `gorm:"primaryKey"`
	OrderID   uint    `gorm:"index;not null"`
	ProductID uint    `gorm:"not null"`
	Quantity  int     `gorm:"not null"`
	Price     float64 `gorm:"not null"` // Price at time
}
