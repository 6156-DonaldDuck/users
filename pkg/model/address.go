package model

import (
	"gorm.io/gorm"
)

type Address struct {
	gorm.Model
	StreetNumber  string `json:"street_number"`
	StreetName1  string `json:"street_name_1"`
	StreetName2  string `json:"street_name_2"`
	City  string `json:"city"`
	Region string `json:"region"`
	CountryCode  string `json:"country_code"`
	PostalCode  string `json:"postal_code"`
}