package db

import (
    "time"
)

// User model represents the users of your application.
type User struct {
    ID        int       `json:"id" gorm:"primaryKey"`
    Name      string    `json:"name"`
    Email     string    `json:"email" gorm:"unique;not null"`
    Password  string    `json:"password" gorm:"not null"`
    Country   string    `json:"country"`   // Field to store country or region
	State     string    `json:"state"`
    CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime"`
    UpdatedAt time.Time `json:"updated_at" gorm:"autoUpdateTime"`
    // Relationship with Transaction model
    Transactions []Transaction `json:"transactions" gorm:"foreignKey:UserID"`
    // Relationship with Session model
    Sessions []Session `json:"sessions" gorm:"foreignKey:UserID"`
}

// Transaction model represents a user's transactions.
type Transaction struct {
    ID        int     `json:"id" gorm:"primaryKey"`
    UserID    int     `json:"user_id" gorm:"not null"`  // Foreign key referencing User
    Amount    float64 `json:"amount" gorm:"not null"`
    Date      time.Time `json:"date" gorm:"autoCreateTime"`
    // Foreign key relation back to the User model
    User      User    `json:"-" gorm:"foreignKey:UserID"`
}


// Session model represents sessions for a user (for tracking user interactions without saving CSV/DB data).
type Session struct {
    ID        int       `json:"id" gorm:"primaryKey"`
    UserID    int       `json:"user_id" gorm:"not null"`  // Foreign key referencing User
    Token     string    `json:"token" gorm:"not null"`
    ExpiresAt time.Time `json:"expires_at"`
    // Foreign key relation back to the User model
    User      User      `json:"-" gorm:"foreignKey:UserID"`
}

