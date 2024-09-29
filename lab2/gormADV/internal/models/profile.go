package models

import "gorm.io/gorm"

type Profile struct {
	gorm.Model
	UserID            uint   `json:"user_id"`
	Bio               string `json:"bio"`
	ProfilePictureURL string `json:"profile_picture_url"`
}
