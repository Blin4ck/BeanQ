package models

import "gorm.io/gorm"

type Menu struct {
	gorm.Model
	ID       uint       `gorm:"primaryKey"`
	Title    string     `json:"title"`
	Position []Position `gorm:"many2many:menu_positions;" json:"positions"`
}
