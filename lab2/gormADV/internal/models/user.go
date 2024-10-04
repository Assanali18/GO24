package models

import "gorm.io/gorm"

// User represents a user in the system.
// swagger:model
type User struct {
	gorm.Model
	// The user's name.
	// example: John Doe
	Name string `json:"name" validate:"required"`
	// The user's age.
	// example: 30
	Age int `json:"age" validate:"required,gte=0"`
	// The user's profile.
	Profile *Profile `json:"profile"`
}

// UserListResponse represents a paginated list of users.
// swagger:model
type UserListResponse struct {
	// The list of users.
	Users []User `json:"users"`
	// The total number of users.
	// example: 100
	TotalItems int `json:"total_items"`
	// The current page number.
	// example: 1
	Page int `json:"page"`
	// The size of each page.
	// example: 10
	PageSize int `json:"page_size"`
	// The total number of pages.
	// example: 10
	TotalPages int `json:"total_pages"`
}
