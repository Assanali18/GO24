package models

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Name    string  `json:"name"`
	Age     int     `json:"age"`
	Profile Profile `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"profile"`
}
