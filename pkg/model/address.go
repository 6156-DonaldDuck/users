package model

import (
	"gorm.io/gorm"
)

type Address struct {
	gorm.Model
	StreetName1 string `json:"street_name_1"`
	StreetName2 string `json:"street_name_2"`
	City        string `json:"city"`
	Region      string `json:"region"`
	CountryCode string `json:"country_code"`
	PostalCode  string `json:"postal_code"`
	UserId      uint   `json:"user_id"`
}