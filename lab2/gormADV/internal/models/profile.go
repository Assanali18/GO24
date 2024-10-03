package models

import "gorm.io/gorm"

type Profile struct {
	gorm.Model
	UserID            uint   `json:"user_id"`
	Bio               string `json:"bio" validate:"omitempty"`
	ProfilePictureURL string `json:"profile_picture_url" validate:"omitempty,url"`
}
