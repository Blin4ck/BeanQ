package models

import "gorm.io/gorm"

type Client struct {
	gorm.Model
	ID        uint    `gorm:"primaryKey"`
	FirstName string  `json:"first_name"`
	LastName  string  `json:"last_name"`
	Email     string  `json:"email"`
	Pasword   string  `json:"password"`
	Orders    []Order `gorm:"foreignKey:ClientID" json:"orders"`
}
