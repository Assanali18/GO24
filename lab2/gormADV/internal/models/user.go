package models

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Name    string   `json:"name" validate:"required"`
	Age     int      `json:"age" validate:"required,gte=0"`
	Profile *Profile `gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"profile"`
}
