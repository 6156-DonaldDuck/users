package model

import (
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	FirstName  string `json:"first_name"`
	LastName  string `json:"last_name"`
	PhoneNumber  string `json:"phone_number"`
	Email  string `json:"email"`
	AddressID uint `json:"address_id"`
}