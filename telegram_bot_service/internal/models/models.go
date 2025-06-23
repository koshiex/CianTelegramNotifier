package models

import (
	"time"

	"gorm.io/gorm"
)

// User represents a Telegram user
type User struct {
	ID        int64     `gorm:"primaryKey" json:"id"`
	Username  string    `json:"username"`
	FirstName string    `json:"first_name"`
	LastName  string    `json:"last_name"`
	Language  string    `json:"language"`
	IsActive  bool      `gorm:"default:true" json:"is_active"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// Favorite represents a user's favorite listing
type Favorite struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	UserID    int64     `json:"user_id"`
	ListingID string    `json:"listing_id"`
	Title     string    `json:"title"`
	Price     string    `json:"price"`
	URL       string    `json:"url"`
	Note      string    `json:"note"`
	CreatedAt time.Time `json:"created_at"`
	User      User      `gorm:"foreignKey:UserID" json:"user"`
}

// Subscription represents a user's notification subscription
type Subscription struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	UserID    int64          `json:"user_id"`
	IsActive  bool           `gorm:"default:true" json:"is_active"`
	Settings  string         `json:"settings"` // JSON string with search settings
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at"`
	User      User           `gorm:"foreignKey:UserID" json:"user"`
}

// Listing represents a property listing from CIAN
type Listing struct {
	ID          string   `json:"id"`
	Title       string   `json:"title"`
	Price       string   `json:"price"`
	PriceValue  int      `json:"price_value"`
	Address     string   `json:"address"`
	URL         string   `json:"url"`
	Description string   `json:"description"`
	Photos      []string `json:"photos"`
	Area        string   `json:"area"`
	Rooms       string   `json:"rooms"`
	Floor       string   `json:"floor"`
	Metro       string   `json:"metro"`
	PublishedAt string   `json:"published_at"`
}
