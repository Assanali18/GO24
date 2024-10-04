package models

import "gorm.io/gorm"

// Profile represents a user's profile.
// swagger:model
type Profile struct {
	gorm.Model
	// The ID of the user this profile belongs to.
	// example: 1
	UserID uint `json:"user_id"`
	// The user's bio.
	// example: Software Developer
	Bio string `json:"bio" validate:"omitempty"`
	// The URL to the user's profile picture.
	// example: http://example.com/profile.jpg
	ProfilePictureURL string `json:"profile_picture_url" validate:"omitempty,url"`
}
