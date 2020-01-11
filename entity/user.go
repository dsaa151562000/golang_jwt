package entity

import (
	"github.com/jinzhu/gorm"
)

// User is user models property
type User struct {
	gorm.Model
	FirstName string `json:"firstname"`
	LastName  string `json:"lastname"`
	Email     string `json:"email"`
	Password  string `json:"password"`
}

// Email 1:N  User -< Email
type Email struct {
	gorm.Model
	UserID int    `json:"user_id" gorm:"index"`
	Email  string `json:"email"`
}

// News 1:N  User -<
type News struct {
	gorm.Model
	UserID int    `json:"user_id" gorm:"index"`
	Title  string `json:"title"`
	Image  string `json:"image" gorm:"size:255"`
	Body   string `json:"body" sql:"type:text;"`
}
