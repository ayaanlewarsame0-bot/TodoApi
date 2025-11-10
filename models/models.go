package models


import (
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type Todo struct {
	gorm.Model
	Title  string `json:"title"`
	Description   string `json:"description"`
	Status string `json:"status"`
	UserID uint   `json:"user_id"`
}
